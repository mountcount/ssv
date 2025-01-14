package runner

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	spectypes "github.com/ssvlabs/ssv-spec/types"

	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/pkg/errors"
	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	genesisspecssv "github.com/ssvlabs/ssv-spec-pre-cc/ssv"
	genesisspectypes "github.com/ssvlabs/ssv-spec-pre-cc/types"
	"go.uber.org/zap"

	"github.com/ssvlabs/ssv/logging/fields"
	"github.com/ssvlabs/ssv/protocol/genesis/ssv/runner/metrics"
)

type ValidatorRegistrationRunner struct {
	BaseRunner *BaseRunner

	beacon   genesisspecssv.BeaconNode
	network  genesisspecssv.Network
	signer   genesisspectypes.KeyManager
	valCheck qbft.ProposedValueCheckF

	metrics metrics.ConsensusMetrics
}

func NewValidatorRegistrationRunner(
	domainType spectypes.DomainType,
	beaconNetwork genesisspectypes.BeaconNetwork,
	share *genesisspectypes.Share,
	beacon genesisspecssv.BeaconNode,
	network genesisspecssv.Network,
	signer genesisspectypes.KeyManager,
) Runner {
	return &ValidatorRegistrationRunner{
		BaseRunner: &BaseRunner{
			BeaconRoleType: genesisspectypes.BNRoleValidatorRegistration,
			DomainType:     domainType,
			BeaconNetwork:  beaconNetwork,
			Share:          share,
		},

		beacon:  beacon,
		network: network,
		signer:  signer,
		metrics: metrics.NewConsensusMetrics(genesisspectypes.BNRoleValidatorRegistration),
	}
}

func (r *ValidatorRegistrationRunner) StartNewDuty(logger *zap.Logger, duty *genesisspectypes.Duty) error {
	return r.BaseRunner.baseStartNewNonBeaconDuty(logger, r, duty)
}

// HasRunningDuty returns true if a duty is already running (StartNewDuty called and returned nil)
func (r *ValidatorRegistrationRunner) HasRunningDuty() bool {
	return r.BaseRunner.hasRunningDuty()
}

func (r *ValidatorRegistrationRunner) ProcessPreConsensus(logger *zap.Logger, signedMsg *genesisspectypes.SignedPartialSignatureMessage) error {
	quorum, roots, err := r.BaseRunner.basePreConsensusMsgProcessing(r, signedMsg)
	if err != nil {
		return errors.Wrap(err, "failed processing validator registration message")
	}

	// quorum returns true only once (first time quorum achieved)
	if !quorum {
		return nil
	}

	// only 1 root, verified in basePreConsensusMsgProcessing
	root := roots[0]
	fullSig, err := r.GetState().ReconstructBeaconSig(r.GetState().PreConsensusContainer, root, r.GetShare().ValidatorPubKey)
	if err != nil {
		// If the reconstructed signature verification failed, fall back to verifying each partial signature
		r.BaseRunner.FallBackAndVerifyEachSignature(r.GetState().PreConsensusContainer, root)
		return errors.Wrap(err, "got pre-consensus quorum but it has invalid signatures")
	}
	specSig := phase0.BLSSignature{}
	copy(specSig[:], fullSig)

	feeRecipient := r.BaseRunner.Share.FeeRecipientAddress
	if err := r.beacon.SubmitValidatorRegistration(r.GetShare().ValidatorPubKey, feeRecipient, specSig); err != nil {
		return errors.Wrap(err, "could not submit validator registration")
	}

	logger.Debug("validator registration submitted successfully",
		fields.FeeRecipient(feeRecipient[:]),
		zap.String("signature", hex.EncodeToString(specSig[:])))

	r.GetState().Finished = true
	return nil
}

func (r *ValidatorRegistrationRunner) ProcessConsensus(logger *zap.Logger, signedMsg *qbft.SignedMessage) error {
	return errors.New("no consensus phase for validator registration")
}

func (r *ValidatorRegistrationRunner) ProcessPostConsensus(logger *zap.Logger, signedMsg *genesisspectypes.SignedPartialSignatureMessage) error {
	return errors.New("no post consensus phase for validator registration")
}

func (r *ValidatorRegistrationRunner) expectedPreConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	vr, err := r.calculateValidatorRegistration()
	if err != nil {
		return nil, genesisspectypes.DomainError, errors.Wrap(err, "could not calculate validator registration")
	}
	return []ssz.HashRoot{vr}, genesisspectypes.DomainApplicationBuilder, nil
}

// expectedPostConsensusRootsAndDomain an INTERNAL function, returns the expected post-consensus roots to sign
func (r *ValidatorRegistrationRunner) expectedPostConsensusRootsAndDomain() ([]ssz.HashRoot, phase0.DomainType, error) {
	return nil, [4]byte{}, errors.New("no post consensus roots for validator registration")
}

func (r *ValidatorRegistrationRunner) executeDuty(logger *zap.Logger, duty *genesisspectypes.Duty) error {
	logger.Debug("executing validator registration duty",
		zap.String("state_fee_recipient", hex.EncodeToString(r.BaseRunner.Share.FeeRecipientAddress[:])))
	vr, err := r.calculateValidatorRegistration()
	if err != nil {
		return errors.Wrap(err, "could not calculate validator registration")
	}

	// sign partial randao
	msg, err := r.BaseRunner.signBeaconObject(r, vr, duty.Slot, genesisspectypes.DomainApplicationBuilder)
	if err != nil {
		return errors.Wrap(err, "could not sign validator registration")
	}
	msgs := genesisspectypes.PartialSignatureMessages{
		Type:     genesisspectypes.ValidatorRegistrationPartialSig,
		Slot:     duty.Slot,
		Messages: []*genesisspectypes.PartialSignatureMessage{msg},
	}

	// sign msg
	signature, err := r.GetSigner().SignRoot(msgs, genesisspectypes.PartialSignatureType, r.GetShare().SharePubKey)
	if err != nil {
		return errors.Wrap(err, "could not sign randao msg")
	}
	signedPartialMsg := &genesisspectypes.SignedPartialSignatureMessage{
		Message:   msgs,
		Signature: signature,
		Signer:    r.GetShare().OperatorID,
	}

	// broadcast
	data, err := signedPartialMsg.Encode()
	if err != nil {
		return errors.Wrap(err, "failed to encode randao pre-consensus signature msg")
	}
	msgToBroadcast := &genesisspectypes.SSVMessage{
		MsgType: genesisspectypes.SSVPartialSignatureMsgType,
		MsgID:   genesisspectypes.NewMsgID(r.GetShare().DomainType, r.GetShare().ValidatorPubKey, r.BaseRunner.BeaconRoleType),
		Data:    data,
	}
	if err := r.GetNetwork().Broadcast(msgToBroadcast); err != nil {
		return errors.Wrap(err, "can't broadcast partial randao sig")
	}
	return nil
}

func (r *ValidatorRegistrationRunner) calculateValidatorRegistration() (*v1.ValidatorRegistration, error) {
	pk := phase0.BLSPubKey{}
	copy(pk[:], r.GetShare().ValidatorPubKey)

	epoch := r.BaseRunner.BeaconNetwork.EstimatedEpochAtSlot(r.BaseRunner.State.StartingDuty.Slot)

	return &v1.ValidatorRegistration{
		FeeRecipient: r.BaseRunner.Share.FeeRecipientAddress,
		GasLimit:     genesisspectypes.DefaultGasLimit,
		Timestamp:    r.BaseRunner.BeaconNetwork.EpochStartTime(epoch),
		Pubkey:       pk,
	}, nil
}

func (r *ValidatorRegistrationRunner) GetBaseRunner() *BaseRunner {
	return r.BaseRunner
}

func (r *ValidatorRegistrationRunner) GetNetwork() genesisspecssv.Network {
	return r.network
}

func (r *ValidatorRegistrationRunner) GetBeaconNode() genesisspecssv.BeaconNode {
	return r.beacon
}

func (r *ValidatorRegistrationRunner) GetShare() *genesisspectypes.Share {
	return r.BaseRunner.Share
}

func (r *ValidatorRegistrationRunner) GetState() *State {
	return r.BaseRunner.State
}

func (r *ValidatorRegistrationRunner) GetValCheckF() qbft.ProposedValueCheckF {
	return r.valCheck
}

func (r *ValidatorRegistrationRunner) GetSigner() genesisspectypes.KeyManager {
	return r.signer
}

// Encode returns the encoded struct in bytes or error
func (r *ValidatorRegistrationRunner) Encode() ([]byte, error) {
	return json.Marshal(r)
}

// Decode returns error if decoding failed
func (r *ValidatorRegistrationRunner) Decode(data []byte) error {
	return json.Unmarshal(data, &r)
}

// GetRoot returns the root used for signing and verification
func (r *ValidatorRegistrationRunner) GetRoot() ([32]byte, error) {
	marshaledRoot, err := r.Encode()
	if err != nil {
		return [32]byte{}, errors.Wrap(err, "could not encode DutyRunnerState")
	}
	ret := sha256.Sum256(marshaledRoot)
	return ret, nil
}
