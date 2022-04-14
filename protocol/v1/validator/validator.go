package validator

import (
	"context"
	"github.com/bloxapp/ssv/protocol/v1/qbft"
	"go.uber.org/zap"
	"time"

	forksprotocol "github.com/bloxapp/ssv/protocol/forks"
	beaconprotocol "github.com/bloxapp/ssv/protocol/v1/blockchain/beacon"
	"github.com/bloxapp/ssv/protocol/v1/message"
	p2pprotocol "github.com/bloxapp/ssv/protocol/v1/p2p"
	"github.com/bloxapp/ssv/protocol/v1/qbft/controller"
	qbftstorage "github.com/bloxapp/ssv/protocol/v1/qbft/storage"
	"github.com/bloxapp/ssv/utils/format"
)

type IValidator interface {
	Start()
	ExecuteDuty(slot uint64, duty *beaconprotocol.Duty)
	ProcessMsg(msg *message.SSVMessage) //TODO need to be as separate interface?
	GetShare() *beaconprotocol.Share
}

type Options struct {
	Context                    context.Context
	Logger                     *zap.Logger
	IbftStorage                qbftstorage.QBFTStore
	Network                    beaconprotocol.Network
	P2pNetwork                 p2pprotocol.Network
	Beacon                     beaconprotocol.Beacon
	Share                      *beaconprotocol.Share
	ForkVersion                forksprotocol.ForkVersion
	Signer                     beaconprotocol.Signer
	SyncRateLimit              time.Duration
	SignatureCollectionTimeout time.Duration
	ReadMode                   bool
	SyncDecided                controller.SyncDecided
	SyncRound                  controller.SyncRound
}

type Validator struct {
	ctx        context.Context
	logger     *zap.Logger
	network    beaconprotocol.Network
	p2pNetwork p2pprotocol.Network
	beacon     beaconprotocol.Beacon
	share      *beaconprotocol.Share
	signer     beaconprotocol.Signer

	ibfts controller.Controllers

	// flags
	readMode bool
}

func NewValidator(opt *Options) IValidator {
	logger := opt.Logger.With(zap.String("pubKey", opt.Share.PublicKey.SerializeToHexStr())).
		With(zap.Uint64("node_id", uint64(opt.Share.NodeID)))

	ibfts := setupIbfts(opt, logger)

	logger.Debug("new validator instance was created", zap.Strings("operators ids", opt.Share.HashOperators()))
	return &Validator{
		ctx:        opt.Context,
		logger:     logger,
		network:    opt.Network,
		p2pNetwork: opt.P2pNetwork,
		beacon:     opt.Beacon,
		share:      opt.Share,
		signer:     opt.Signer,
		ibfts:      ibfts,
		readMode:   opt.ReadMode,
	}
}

func (v *Validator) Start() {
}

func (v *Validator) GetShare() *beaconprotocol.Share {
	// TODO need lock?
	return v.share
}

// ProcessMsg processes a new msg, returns true if Decided, non nil byte slice if Decided (Decided value) and error
// Decided returns just once per instance as true, following messages (for example additional commit msgs) will not return Decided true
func (v *Validator) ProcessMsg(msg *message.SSVMessage) /*(bool, []byte, error)*/ {
	ibftController := v.ibfts.ControllerForIdentifier(msg.GetIdentifier())
	// synchronize process
	err := ibftController.ProcessMsg(msg)
	if err != nil {
		return
	}
	//	  TODO need to handle decided and decidedValue
}

// setupRunners return duty runners map with all the supported duty types
func setupIbfts(opt *Options, logger *zap.Logger) map[message.RoleType]controller.IController {
	ibfts := make(map[message.RoleType]controller.IController)
	ibfts[message.RoleTypeAttester] = setupIbftController(message.RoleTypeAttester, logger, opt)
	return ibfts
}

func setupIbftController(role message.RoleType, logger *zap.Logger, opt *Options) controller.IController {
	identifier := []byte(format.IdentifierFormat(opt.Share.PublicKey.Serialize(), role.String()))

	opts := controller.Options{
		Role:           role,
		Identifier:     identifier,
		Logger:         logger,
		Storage:        opt.IbftStorage,
		Network:        opt.P2pNetwork,
		InstanceConfig: qbft.DefaultConsensusParams(),
		ValidatorShare: opt.Share,
		Version:        opt.ForkVersion,
		Beacon:         opt.Beacon,
		Signer:         opt.Signer,
		SyncRateLimit:  opt.SyncRateLimit,
		SigTimeout:     opt.SignatureCollectionTimeout,
		ReadMode:       opt.ReadMode,
		SyncRound:      opt.SyncRound,
		SyncDecided:    opt.SyncDecided,
	}
	return controller.New(opts)
}
