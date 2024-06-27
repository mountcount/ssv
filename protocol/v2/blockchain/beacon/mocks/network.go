// Code generated by MockGen. DO NOT EDIT.
// Source: ./network.go
//
// Generated by this command:
//
//	mockgen -package=mocks -destination=./mocks/network.go -source=./network.go
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	phase0 "github.com/attestantio/go-eth2-client/spec/phase0"
	types "github.com/ssvlabs/ssv-spec/types"
	beacon "github.com/ssvlabs/ssv/protocol/v2/blockchain/beacon"
	gomock "go.uber.org/mock/gomock"
)

// MockBeaconNetwork is a mock of BeaconNetwork interface.
type MockBeaconNetwork struct {
	ctrl     *gomock.Controller
	recorder *MockBeaconNetworkMockRecorder
}

// MockBeaconNetworkMockRecorder is the mock recorder for MockBeaconNetwork.
type MockBeaconNetworkMockRecorder struct {
	mock *MockBeaconNetwork
}

// NewMockBeaconNetwork creates a new mock instance.
func NewMockBeaconNetwork(ctrl *gomock.Controller) *MockBeaconNetwork {
	mock := &MockBeaconNetwork{ctrl: ctrl}
	mock.recorder = &MockBeaconNetworkMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBeaconNetwork) EXPECT() *MockBeaconNetworkMockRecorder {
	return m.recorder
}

// EpochStartTime mocks base method.
func (m *MockBeaconNetwork) EpochStartTime(epoch phase0.Epoch) time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EpochStartTime", epoch)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// EpochStartTime indicates an expected call of EpochStartTime.
func (mr *MockBeaconNetworkMockRecorder) EpochStartTime(epoch any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EpochStartTime", reflect.TypeOf((*MockBeaconNetwork)(nil).EpochStartTime), epoch)
}

// EpochsPerSyncCommitteePeriod mocks base method.
func (m *MockBeaconNetwork) EpochsPerSyncCommitteePeriod() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EpochsPerSyncCommitteePeriod")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// EpochsPerSyncCommitteePeriod indicates an expected call of EpochsPerSyncCommitteePeriod.
func (mr *MockBeaconNetworkMockRecorder) EpochsPerSyncCommitteePeriod() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EpochsPerSyncCommitteePeriod", reflect.TypeOf((*MockBeaconNetwork)(nil).EpochsPerSyncCommitteePeriod))
}

// EstimatedCurrentEpoch mocks base method.
func (m *MockBeaconNetwork) EstimatedCurrentEpoch() phase0.Epoch {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EstimatedCurrentEpoch")
	ret0, _ := ret[0].(phase0.Epoch)
	return ret0
}

// EstimatedCurrentEpoch indicates an expected call of EstimatedCurrentEpoch.
func (mr *MockBeaconNetworkMockRecorder) EstimatedCurrentEpoch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EstimatedCurrentEpoch", reflect.TypeOf((*MockBeaconNetwork)(nil).EstimatedCurrentEpoch))
}

// EstimatedCurrentSlot mocks base method.
func (m *MockBeaconNetwork) EstimatedCurrentSlot() phase0.Slot {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EstimatedCurrentSlot")
	ret0, _ := ret[0].(phase0.Slot)
	return ret0
}

// EstimatedCurrentSlot indicates an expected call of EstimatedCurrentSlot.
func (mr *MockBeaconNetworkMockRecorder) EstimatedCurrentSlot() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EstimatedCurrentSlot", reflect.TypeOf((*MockBeaconNetwork)(nil).EstimatedCurrentSlot))
}

// EstimatedEpochAtSlot mocks base method.
func (m *MockBeaconNetwork) EstimatedEpochAtSlot(slot phase0.Slot) phase0.Epoch {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EstimatedEpochAtSlot", slot)
	ret0, _ := ret[0].(phase0.Epoch)
	return ret0
}

// EstimatedEpochAtSlot indicates an expected call of EstimatedEpochAtSlot.
func (mr *MockBeaconNetworkMockRecorder) EstimatedEpochAtSlot(slot any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EstimatedEpochAtSlot", reflect.TypeOf((*MockBeaconNetwork)(nil).EstimatedEpochAtSlot), slot)
}

// EstimatedSlotAtTime mocks base method.
func (m *MockBeaconNetwork) EstimatedSlotAtTime(time int64) phase0.Slot {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EstimatedSlotAtTime", time)
	ret0, _ := ret[0].(phase0.Slot)
	return ret0
}

// EstimatedSlotAtTime indicates an expected call of EstimatedSlotAtTime.
func (mr *MockBeaconNetworkMockRecorder) EstimatedSlotAtTime(time any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EstimatedSlotAtTime", reflect.TypeOf((*MockBeaconNetwork)(nil).EstimatedSlotAtTime), time)
}

// EstimatedSyncCommitteePeriodAtEpoch mocks base method.
func (m *MockBeaconNetwork) EstimatedSyncCommitteePeriodAtEpoch(epoch phase0.Epoch) uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EstimatedSyncCommitteePeriodAtEpoch", epoch)
	ret0, _ := ret[0].(uint64)
	return ret0
}

// EstimatedSyncCommitteePeriodAtEpoch indicates an expected call of EstimatedSyncCommitteePeriodAtEpoch.
func (mr *MockBeaconNetworkMockRecorder) EstimatedSyncCommitteePeriodAtEpoch(epoch any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EstimatedSyncCommitteePeriodAtEpoch", reflect.TypeOf((*MockBeaconNetwork)(nil).EstimatedSyncCommitteePeriodAtEpoch), epoch)
}

// EstimatedTimeAtSlot mocks base method.
func (m *MockBeaconNetwork) EstimatedTimeAtSlot(slot phase0.Slot) int64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EstimatedTimeAtSlot", slot)
	ret0, _ := ret[0].(int64)
	return ret0
}

// EstimatedTimeAtSlot indicates an expected call of EstimatedTimeAtSlot.
func (mr *MockBeaconNetworkMockRecorder) EstimatedTimeAtSlot(slot any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EstimatedTimeAtSlot", reflect.TypeOf((*MockBeaconNetwork)(nil).EstimatedTimeAtSlot), slot)
}

// FirstEpochOfSyncPeriod mocks base method.
func (m *MockBeaconNetwork) FirstEpochOfSyncPeriod(period uint64) phase0.Epoch {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FirstEpochOfSyncPeriod", period)
	ret0, _ := ret[0].(phase0.Epoch)
	return ret0
}

// FirstEpochOfSyncPeriod indicates an expected call of FirstEpochOfSyncPeriod.
func (mr *MockBeaconNetworkMockRecorder) FirstEpochOfSyncPeriod(period any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FirstEpochOfSyncPeriod", reflect.TypeOf((*MockBeaconNetwork)(nil).FirstEpochOfSyncPeriod), period)
}

// FirstSlotAtEpoch mocks base method.
func (m *MockBeaconNetwork) FirstSlotAtEpoch(epoch phase0.Epoch) phase0.Slot {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FirstSlotAtEpoch", epoch)
	ret0, _ := ret[0].(phase0.Slot)
	return ret0
}

// FirstSlotAtEpoch indicates an expected call of FirstSlotAtEpoch.
func (mr *MockBeaconNetworkMockRecorder) FirstSlotAtEpoch(epoch any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FirstSlotAtEpoch", reflect.TypeOf((*MockBeaconNetwork)(nil).FirstSlotAtEpoch), epoch)
}

// ForkVersion mocks base method.
func (m *MockBeaconNetwork) ForkVersion() [4]byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ForkVersion")
	ret0, _ := ret[0].([4]byte)
	return ret0
}

// ForkVersion indicates an expected call of ForkVersion.
func (mr *MockBeaconNetworkMockRecorder) ForkVersion() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ForkVersion", reflect.TypeOf((*MockBeaconNetwork)(nil).ForkVersion))
}

// GetBeaconNetwork mocks base method.
func (m *MockBeaconNetwork) GetBeaconNetwork() types.BeaconNetwork {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBeaconNetwork")
	ret0, _ := ret[0].(types.BeaconNetwork)
	return ret0
}

// GetBeaconNetwork indicates an expected call of GetBeaconNetwork.
func (mr *MockBeaconNetworkMockRecorder) GetBeaconNetwork() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBeaconNetwork", reflect.TypeOf((*MockBeaconNetwork)(nil).GetBeaconNetwork))
}

// GetEpochFirstSlot mocks base method.
func (m *MockBeaconNetwork) GetEpochFirstSlot(epoch phase0.Epoch) phase0.Slot {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEpochFirstSlot", epoch)
	ret0, _ := ret[0].(phase0.Slot)
	return ret0
}

// GetEpochFirstSlot indicates an expected call of GetEpochFirstSlot.
func (mr *MockBeaconNetworkMockRecorder) GetEpochFirstSlot(epoch any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEpochFirstSlot", reflect.TypeOf((*MockBeaconNetwork)(nil).GetEpochFirstSlot), epoch)
}

// GetNetwork mocks base method.
func (m *MockBeaconNetwork) GetNetwork() beacon.Network {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetwork")
	ret0, _ := ret[0].(beacon.Network)
	return ret0
}

// GetNetwork indicates an expected call of GetNetwork.
func (mr *MockBeaconNetworkMockRecorder) GetNetwork() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetwork", reflect.TypeOf((*MockBeaconNetwork)(nil).GetNetwork))
}

// GetSlotEndTime mocks base method.
func (m *MockBeaconNetwork) GetSlotEndTime(slot phase0.Slot) time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSlotEndTime", slot)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// GetSlotEndTime indicates an expected call of GetSlotEndTime.
func (mr *MockBeaconNetworkMockRecorder) GetSlotEndTime(slot any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSlotEndTime", reflect.TypeOf((*MockBeaconNetwork)(nil).GetSlotEndTime), slot)
}

// GetSlotStartTime mocks base method.
func (m *MockBeaconNetwork) GetSlotStartTime(slot phase0.Slot) time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSlotStartTime", slot)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// GetSlotStartTime indicates an expected call of GetSlotStartTime.
func (mr *MockBeaconNetworkMockRecorder) GetSlotStartTime(slot any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSlotStartTime", reflect.TypeOf((*MockBeaconNetwork)(nil).GetSlotStartTime), slot)
}

// IsFirstSlotOfEpoch mocks base method.
func (m *MockBeaconNetwork) IsFirstSlotOfEpoch(slot phase0.Slot) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsFirstSlotOfEpoch", slot)
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsFirstSlotOfEpoch indicates an expected call of IsFirstSlotOfEpoch.
func (mr *MockBeaconNetworkMockRecorder) IsFirstSlotOfEpoch(slot any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsFirstSlotOfEpoch", reflect.TypeOf((*MockBeaconNetwork)(nil).IsFirstSlotOfEpoch), slot)
}

// LastSlotOfSyncPeriod mocks base method.
func (m *MockBeaconNetwork) LastSlotOfSyncPeriod(period uint64) phase0.Slot {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastSlotOfSyncPeriod", period)
	ret0, _ := ret[0].(phase0.Slot)
	return ret0
}

// LastSlotOfSyncPeriod indicates an expected call of LastSlotOfSyncPeriod.
func (mr *MockBeaconNetworkMockRecorder) LastSlotOfSyncPeriod(period any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastSlotOfSyncPeriod", reflect.TypeOf((*MockBeaconNetwork)(nil).LastSlotOfSyncPeriod), period)
}

// MinGenesisTime mocks base method.
func (m *MockBeaconNetwork) MinGenesisTime() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MinGenesisTime")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// MinGenesisTime indicates an expected call of MinGenesisTime.
func (mr *MockBeaconNetworkMockRecorder) MinGenesisTime() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MinGenesisTime", reflect.TypeOf((*MockBeaconNetwork)(nil).MinGenesisTime))
}

// SlotDurationSec mocks base method.
func (m *MockBeaconNetwork) SlotDurationSec() time.Duration {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SlotDurationSec")
	ret0, _ := ret[0].(time.Duration)
	return ret0
}

// SlotDurationSec indicates an expected call of SlotDurationSec.
func (mr *MockBeaconNetworkMockRecorder) SlotDurationSec() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SlotDurationSec", reflect.TypeOf((*MockBeaconNetwork)(nil).SlotDurationSec))
}

// SlotsPerEpoch mocks base method.
func (m *MockBeaconNetwork) SlotsPerEpoch() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SlotsPerEpoch")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// SlotsPerEpoch indicates an expected call of SlotsPerEpoch.
func (mr *MockBeaconNetworkMockRecorder) SlotsPerEpoch() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SlotsPerEpoch", reflect.TypeOf((*MockBeaconNetwork)(nil).SlotsPerEpoch))
}
