package duties

import (
	"context"
	"testing"

	eth2apiv1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	genesisspectypes "github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/ssvlabs/ssv/utils/hashmap"

	spectypes "github.com/ssvlabs/ssv-spec/types"

	"github.com/ssvlabs/ssv/beacon/goclient"
	"github.com/ssvlabs/ssv/operator/duties/dutystore"
	"github.com/ssvlabs/ssv/protocol/v2/types"
)

func setupProposerGenesisDutiesMock(s *Scheduler, dutiesMap *hashmap.Map[phase0.Epoch, []*eth2apiv1.ProposerDuty]) (chan struct{}, chan []*genesisspectypes.Duty) {
	fetchDutiesCall := make(chan struct{})
	executeDutiesCall := make(chan []*genesisspectypes.Duty)

	s.beaconNode.(*MockBeaconNode).EXPECT().ProposerDuties(gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, epoch phase0.Epoch, indices []phase0.ValidatorIndex) ([]*eth2apiv1.ProposerDuty, error) {
			fetchDutiesCall <- struct{}{}
			duties, _ := dutiesMap.Get(epoch)
			return duties, nil
		}).AnyTimes()

	getShares := func(epoch phase0.Epoch) []*types.SSVShare {
		uniqueIndices := make(map[phase0.ValidatorIndex]bool)

		duties, _ := dutiesMap.Get(epoch)
		for _, d := range duties {
			uniqueIndices[d.ValidatorIndex] = true
		}

		shares := make([]*types.SSVShare, 0, len(uniqueIndices))
		for index := range uniqueIndices {
			share := &types.SSVShare{
				Share: spectypes.Share{
					ValidatorIndex: index,
				},
			}
			shares = append(shares, share)
		}

		return shares
	}

	s.validatorProvider.(*MockValidatorProvider).EXPECT().SelfParticipatingValidators(gomock.Any()).DoAndReturn(getShares).AnyTimes()
	s.validatorProvider.(*MockValidatorProvider).EXPECT().ParticipatingValidators(gomock.Any()).DoAndReturn(getShares).AnyTimes()

	return fetchDutiesCall, executeDutiesCall
}

func expectedExecutedGenesisProposerDuties(handler *ProposerHandler, duties []*eth2apiv1.ProposerDuty) []*genesisspectypes.Duty {
	expectedDuties := make([]*genesisspectypes.Duty, 0)
	for _, d := range duties {
		expectedDuties = append(expectedDuties, handler.toGenesisSpecDuty(d, genesisspectypes.BNRoleProposer))
	}
	return expectedDuties
}

func TestScheduler_Proposer_Genesis_Same_Slot(t *testing.T) {
	var (
		handler     = NewProposerHandler(dutystore.NewDuties[eth2apiv1.ProposerDuty]())
		currentSlot = &SafeValue[phase0.Slot]{}
		dutiesMap   = hashmap.New[phase0.Epoch, []*eth2apiv1.ProposerDuty]()
	)
	currentSlot.Set(phase0.Slot(0))
	scheduler, logger, ticker, timeout, cancel, schedulerPool, startFn := setupSchedulerAndMocks(t, []dutyHandler{handler}, currentSlot, goclient.FarFutureEpoch)
	fetchDutiesCall, executeDutiesCall := setupProposerGenesisDutiesMock(scheduler, dutiesMap)
	startFn()

	dutiesMap.Set(phase0.Epoch(0), []*eth2apiv1.ProposerDuty{
		{
			PubKey:         phase0.BLSPubKey{1, 2, 3},
			Slot:           phase0.Slot(0),
			ValidatorIndex: phase0.ValidatorIndex(1),
		},
	})

	// STEP 1: wait for proposer duties to be fetched and executed at the same slot
	duties, _ := dutiesMap.Get(phase0.Epoch(0))
	expected := expectedExecutedGenesisProposerDuties(handler, duties)
	setExecuteGenesisDutyFunc(scheduler, executeDutiesCall, len(expected))

	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesFetch(t, logger, fetchDutiesCall, executeDutiesCall, timeout)
	waitForGenesisDutiesExecution(t, logger, fetchDutiesCall, executeDutiesCall, timeout, expected)

	// Stop scheduler & wait for graceful exit.
	cancel()
	require.NoError(t, schedulerPool.Wait())
}

func TestScheduler_Proposer_Genesis_Diff_Slots(t *testing.T) {
	var (
		handler     = NewProposerHandler(dutystore.NewDuties[eth2apiv1.ProposerDuty]())
		currentSlot = &SafeValue[phase0.Slot]{}
		dutiesMap   = hashmap.New[phase0.Epoch, []*eth2apiv1.ProposerDuty]()
	)
	currentSlot.Set(phase0.Slot(0))
	scheduler, logger, ticker, timeout, cancel, schedulerPool, startFn := setupSchedulerAndMocks(t, []dutyHandler{handler}, currentSlot, goclient.FarFutureEpoch)
	fetchDutiesCall, executeDutiesCall := setupProposerGenesisDutiesMock(scheduler, dutiesMap)
	startFn()

	dutiesMap.Set(phase0.Epoch(0), []*eth2apiv1.ProposerDuty{
		{
			PubKey:         phase0.BLSPubKey{1, 2, 3},
			Slot:           phase0.Slot(2),
			ValidatorIndex: phase0.ValidatorIndex(1),
		},
	})

	// STEP 1: wait for proposer duties to be fetched
	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesFetch(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 2: wait for no action to be taken
	currentSlot.Set(phase0.Slot(1))
	ticker.Send(currentSlot.Get())
	waitForNoActionGenesis(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 3: wait for proposer duties to be executed
	currentSlot.Set(phase0.Slot(2))
	duties, _ := dutiesMap.Get(phase0.Epoch(0))
	expected := expectedExecutedGenesisProposerDuties(handler, duties)
	setExecuteGenesisDutyFunc(scheduler, executeDutiesCall, len(expected))

	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesExecution(t, logger, fetchDutiesCall, executeDutiesCall, timeout, expected)

	// Stop scheduler & wait for graceful exit.
	cancel()
	require.NoError(t, schedulerPool.Wait())
}

// execute duty after two slots after the indices changed
func TestScheduler_Proposer_Genesis_Indices_Changed(t *testing.T) {
	var (
		handler     = NewProposerHandler(dutystore.NewDuties[eth2apiv1.ProposerDuty]())
		currentSlot = &SafeValue[phase0.Slot]{}
		dutiesMap   = hashmap.New[phase0.Epoch, []*eth2apiv1.ProposerDuty]()
	)
	currentSlot.Set(phase0.Slot(0))
	scheduler, logger, ticker, timeout, cancel, schedulerPool, startFn := setupSchedulerAndMocks(t, []dutyHandler{handler}, currentSlot, goclient.FarFutureEpoch)
	fetchDutiesCall, executeDutiesCall := setupProposerGenesisDutiesMock(scheduler, dutiesMap)
	startFn()

	// STEP 1: wait for no action to be taken
	ticker.Send(currentSlot.Get())
	waitForNoActionGenesis(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 2: wait for no action to be taken
	currentSlot.Set(phase0.Slot(1))
	ticker.Send(currentSlot.Get())
	waitForNoActionGenesis(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 3: trigger a change in active indices
	scheduler.indicesChg <- struct{}{}
	dutiesMap.Set(phase0.Epoch(0), []*eth2apiv1.ProposerDuty{
		{
			PubKey:         phase0.BLSPubKey{1, 2, 3},
			Slot:           phase0.Slot(1),
			ValidatorIndex: phase0.ValidatorIndex(1),
		},
		{
			PubKey:         phase0.BLSPubKey{1, 2, 4},
			Slot:           phase0.Slot(2),
			ValidatorIndex: phase0.ValidatorIndex(2),
		},
		{
			PubKey:         phase0.BLSPubKey{1, 2, 5},
			Slot:           phase0.Slot(3),
			ValidatorIndex: phase0.ValidatorIndex(3),
		},
	})
	// no execution should happen in slot 1
	waitForNoActionGenesis(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 4: wait for proposer duties to be fetched again
	currentSlot.Set(phase0.Slot(2))
	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesFetch(t, logger, fetchDutiesCall, executeDutiesCall, timeout)
	// no execution should happen in slot 2
	waitForNoActionGenesis(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 4: wait for proposer duties to be executed
	currentSlot.Set(phase0.Slot(3))
	duties, _ := dutiesMap.Get(phase0.Epoch(0))
	expected := expectedExecutedGenesisProposerDuties(handler, []*eth2apiv1.ProposerDuty{duties[2]})
	setExecuteGenesisDutyFunc(scheduler, executeDutiesCall, len(expected))

	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesExecution(t, logger, fetchDutiesCall, executeDutiesCall, timeout, expected)

	// Stop scheduler & wait for graceful exit.
	cancel()
	require.NoError(t, schedulerPool.Wait())
}

func TestScheduler_Proposer_Genesis_Multiple_Indices_Changed_Same_Slot(t *testing.T) {
	var (
		handler     = NewProposerHandler(dutystore.NewDuties[eth2apiv1.ProposerDuty]())
		currentSlot = &SafeValue[phase0.Slot]{}
		dutiesMap   = hashmap.New[phase0.Epoch, []*eth2apiv1.ProposerDuty]()
	)
	currentSlot.Set(phase0.Slot(0))
	scheduler, logger, ticker, timeout, cancel, schedulerPool, startFn := setupSchedulerAndMocks(t, []dutyHandler{handler}, currentSlot, goclient.FarFutureEpoch)
	fetchDutiesCall, executeDutiesCall := setupProposerGenesisDutiesMock(scheduler, dutiesMap)
	startFn()

	dutiesMap.Set(phase0.Epoch(0), []*eth2apiv1.ProposerDuty{
		{
			PubKey:         phase0.BLSPubKey{1, 2, 3},
			Slot:           phase0.Slot(2),
			ValidatorIndex: phase0.ValidatorIndex(1),
		},
	})

	// STEP 1: wait for proposer duties to be fetched
	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesFetch(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 2: trigger a change in active indices
	scheduler.indicesChg <- struct{}{}
	waitForNoActionGenesis(t, logger, fetchDutiesCall, executeDutiesCall, timeout)
	duties, _ := dutiesMap.Get(phase0.Epoch(0))
	dutiesMap.Set(phase0.Epoch(0), append(duties, &eth2apiv1.ProposerDuty{
		PubKey:         phase0.BLSPubKey{1, 2, 4},
		Slot:           phase0.Slot(3),
		ValidatorIndex: phase0.ValidatorIndex(2),
	}))

	// STEP 3: trigger a change in active indices in the same slot
	scheduler.indicesChg <- struct{}{}
	waitForNoActionGenesis(t, logger, fetchDutiesCall, executeDutiesCall, timeout)
	duties, _ = dutiesMap.Get(phase0.Epoch(0))
	dutiesMap.Set(phase0.Epoch(0), append(duties, &eth2apiv1.ProposerDuty{
		PubKey:         phase0.BLSPubKey{1, 2, 5},
		Slot:           phase0.Slot(4),
		ValidatorIndex: phase0.ValidatorIndex(3),
	}))

	// STEP 4: wait for proposer duties to be fetched again
	currentSlot.Set(phase0.Slot(1))
	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesFetch(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 5: wait for proposer duties to be executed
	currentSlot.Set(phase0.Slot(2))
	duties, _ = dutiesMap.Get(phase0.Epoch(0))
	expected := expectedExecutedGenesisProposerDuties(handler, []*eth2apiv1.ProposerDuty{duties[0]})
	setExecuteGenesisDutyFunc(scheduler, executeDutiesCall, len(expected))

	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesExecution(t, logger, fetchDutiesCall, executeDutiesCall, timeout, expected)

	// STEP 6: wait for proposer duties to be executed
	currentSlot.Set(phase0.Slot(3))
	duties, _ = dutiesMap.Get(phase0.Epoch(0))
	expected = expectedExecutedGenesisProposerDuties(handler, []*eth2apiv1.ProposerDuty{duties[1]})
	setExecuteGenesisDutyFunc(scheduler, executeDutiesCall, len(expected))

	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesExecution(t, logger, fetchDutiesCall, executeDutiesCall, timeout, expected)

	// STEP 7: wait for proposer duties to be executed
	currentSlot.Set(phase0.Slot(4))
	duties, _ = dutiesMap.Get(phase0.Epoch(0))
	expected = expectedExecutedGenesisProposerDuties(handler, []*eth2apiv1.ProposerDuty{duties[2]})
	setExecuteGenesisDutyFunc(scheduler, executeDutiesCall, len(expected))

	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesExecution(t, logger, fetchDutiesCall, executeDutiesCall, timeout, expected)

	// Stop scheduler & wait for graceful exit.
	cancel()
	require.NoError(t, schedulerPool.Wait())
}

// reorg current dependent root changed
func TestScheduler_Proposer_Genesis_Reorg_Current(t *testing.T) {
	var (
		handler     = NewProposerHandler(dutystore.NewDuties[eth2apiv1.ProposerDuty]())
		currentSlot = &SafeValue[phase0.Slot]{}
		dutiesMap   = hashmap.New[phase0.Epoch, []*eth2apiv1.ProposerDuty]()
	)
	currentSlot.Set(phase0.Slot(34))
	scheduler, logger, ticker, timeout, cancel, schedulerPool, startFn := setupSchedulerAndMocks(t, []dutyHandler{handler}, currentSlot, goclient.FarFutureEpoch)
	fetchDutiesCall, executeDutiesCall := setupProposerGenesisDutiesMock(scheduler, dutiesMap)
	startFn()

	dutiesMap.Set(phase0.Epoch(1), []*eth2apiv1.ProposerDuty{
		{
			PubKey:         phase0.BLSPubKey{1, 2, 3},
			Slot:           phase0.Slot(36),
			ValidatorIndex: phase0.ValidatorIndex(1),
		},
	})

	// STEP 1: wait for proposer duties to be fetched
	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesFetch(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 2: trigger head event
	e := &eth2apiv1.Event{
		Data: &eth2apiv1.HeadEvent{
			Slot:                     currentSlot.Get(),
			CurrentDutyDependentRoot: phase0.Root{0x01},
		},
	}
	scheduler.HandleHeadEvent(logger)(e)
	waitForNoActionGenesis(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 3: Ticker with no action
	currentSlot.Set(phase0.Slot(35))
	ticker.Send(currentSlot.Get())
	waitForNoActionGenesis(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 4: trigger reorg
	e = &eth2apiv1.Event{
		Data: &eth2apiv1.HeadEvent{
			Slot:                     currentSlot.Get(),
			CurrentDutyDependentRoot: phase0.Root{0x02},
		},
	}
	dutiesMap.Set(phase0.Epoch(1), []*eth2apiv1.ProposerDuty{
		{
			PubKey:         phase0.BLSPubKey{1, 2, 3},
			Slot:           phase0.Slot(37),
			ValidatorIndex: phase0.ValidatorIndex(1),
		},
	})
	scheduler.HandleHeadEvent(logger)(e)
	waitForNoActionGenesis(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 5: wait for proposer duties to be fetched again for the current epoch.
	// The first assigned duty should not be executed
	currentSlot.Set(phase0.Slot(36))
	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesFetch(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 7: The second assigned duty should be executed
	currentSlot.Set(phase0.Slot(37))
	duties, _ := dutiesMap.Get(phase0.Epoch(1))
	expected := expectedExecutedGenesisProposerDuties(handler, duties)
	setExecuteGenesisDutyFunc(scheduler, executeDutiesCall, len(expected))

	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesExecution(t, logger, fetchDutiesCall, executeDutiesCall, timeout, expected)

	// Stop scheduler & wait for graceful exit.
	cancel()
	require.NoError(t, schedulerPool.Wait())
}

// reorg current dependent root changed
func TestScheduler_Proposer_Genesis_Reorg_Current_Indices_Changed(t *testing.T) {
	var (
		handler     = NewProposerHandler(dutystore.NewDuties[eth2apiv1.ProposerDuty]())
		currentSlot = &SafeValue[phase0.Slot]{}
		dutiesMap   = hashmap.New[phase0.Epoch, []*eth2apiv1.ProposerDuty]()
	)
	currentSlot.Set(phase0.Slot(34))
	scheduler, logger, ticker, timeout, cancel, schedulerPool, startFn := setupSchedulerAndMocks(t, []dutyHandler{handler}, currentSlot, goclient.FarFutureEpoch)
	fetchDutiesCall, executeDutiesCall := setupProposerGenesisDutiesMock(scheduler, dutiesMap)
	startFn()

	dutiesMap.Set(phase0.Epoch(1), []*eth2apiv1.ProposerDuty{
		{
			PubKey:         phase0.BLSPubKey{1, 2, 3},
			Slot:           phase0.Slot(36),
			ValidatorIndex: phase0.ValidatorIndex(1),
		},
	})

	// STEP 1: wait for proposer duties to be fetched
	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesFetch(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 2: trigger head event
	e := &eth2apiv1.Event{
		Data: &eth2apiv1.HeadEvent{
			Slot:                     currentSlot.Get(),
			CurrentDutyDependentRoot: phase0.Root{0x01},
		},
	}
	scheduler.HandleHeadEvent(logger)(e)
	waitForNoActionGenesis(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 3: Ticker with no action
	currentSlot.Set(phase0.Slot(35))
	ticker.Send(currentSlot.Get())
	waitForNoActionGenesis(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 4: trigger reorg
	e = &eth2apiv1.Event{
		Data: &eth2apiv1.HeadEvent{
			Slot:                     currentSlot.Get(),
			CurrentDutyDependentRoot: phase0.Root{0x02},
		},
	}
	dutiesMap.Set(phase0.Epoch(1), []*eth2apiv1.ProposerDuty{
		{
			PubKey:         phase0.BLSPubKey{1, 2, 3},
			Slot:           phase0.Slot(37),
			ValidatorIndex: phase0.ValidatorIndex(1),
		},
	})
	scheduler.HandleHeadEvent(logger)(e)
	waitForNoActionGenesis(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 5: trigger a change in active indices in the same slot
	scheduler.indicesChg <- struct{}{}
	duties, _ := dutiesMap.Get(phase0.Epoch(1))
	dutiesMap.Set(phase0.Epoch(1), append(duties, &eth2apiv1.ProposerDuty{
		PubKey:         phase0.BLSPubKey{1, 2, 4},
		Slot:           phase0.Slot(38),
		ValidatorIndex: phase0.ValidatorIndex(2),
	}))
	waitForNoActionGenesis(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 6: wait for proposer duties to be fetched again for the current epoch.
	// The first assigned duty should not be executed
	currentSlot.Set(phase0.Slot(36))
	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesFetch(t, logger, fetchDutiesCall, executeDutiesCall, timeout)

	// STEP 7: The second assigned duty should be executed
	currentSlot.Set(phase0.Slot(37))
	duties, _ = dutiesMap.Get(phase0.Epoch(1))
	expected := expectedExecutedGenesisProposerDuties(handler, []*eth2apiv1.ProposerDuty{duties[0]})
	setExecuteGenesisDutyFunc(scheduler, executeDutiesCall, len(expected))

	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesExecution(t, logger, fetchDutiesCall, executeDutiesCall, timeout, expected)

	// STEP 8: The second assigned duty should be executed
	currentSlot.Set(phase0.Slot(38))
	duties, _ = dutiesMap.Get(phase0.Epoch(1))
	expected = expectedExecutedGenesisProposerDuties(handler, []*eth2apiv1.ProposerDuty{duties[1]})
	setExecuteGenesisDutyFunc(scheduler, executeDutiesCall, len(expected))

	ticker.Send(currentSlot.Get())
	waitForGenesisDutiesExecution(t, logger, fetchDutiesCall, executeDutiesCall, timeout, expected)

	// Stop scheduler & wait for graceful exit.
	cancel()
	require.NoError(t, schedulerPool.Wait())
}
