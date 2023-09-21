// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/confluentinc/cli/v3/pkg/flink/types (interfaces: ResultFetcherInterface)
//
// Generated by this command:
//
//	mockgen -destination result_fetcher_mock.go -package=mock github.com/confluentinc/cli/v3/pkg/flink/types ResultFetcherInterface
//
// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"
	time "time"

	types "github.com/confluentinc/cli/v3/pkg/flink/types"
	gomock "go.uber.org/mock/gomock"
)

// MockResultFetcherInterface is a mock of ResultFetcherInterface interface.
type MockResultFetcherInterface struct {
	ctrl     *gomock.Controller
	recorder *MockResultFetcherInterfaceMockRecorder
}

// MockResultFetcherInterfaceMockRecorder is the mock recorder for MockResultFetcherInterface.
type MockResultFetcherInterfaceMockRecorder struct {
	mock *MockResultFetcherInterface
}

// NewMockResultFetcherInterface creates a new mock instance.
func NewMockResultFetcherInterface(ctrl *gomock.Controller) *MockResultFetcherInterface {
	mock := &MockResultFetcherInterface{ctrl: ctrl}
	mock.recorder = &MockResultFetcherInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResultFetcherInterface) EXPECT() *MockResultFetcherInterfaceMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockResultFetcherInterface) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockResultFetcherInterfaceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockResultFetcherInterface)(nil).Close))
}

// GetLastRefreshTimestamp mocks base method.
func (m *MockResultFetcherInterface) GetLastRefreshTimestamp() *time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastRefreshTimestamp")
	ret0, _ := ret[0].(*time.Time)
	return ret0
}

// GetLastRefreshTimestamp indicates an expected call of GetLastRefreshTimestamp.
func (mr *MockResultFetcherInterfaceMockRecorder) GetLastRefreshTimestamp() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastRefreshTimestamp", reflect.TypeOf((*MockResultFetcherInterface)(nil).GetLastRefreshTimestamp))
}

// GetMaterializedStatementResults mocks base method.
func (m *MockResultFetcherInterface) GetMaterializedStatementResults() *types.MaterializedStatementResults {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMaterializedStatementResults")
	ret0, _ := ret[0].(*types.MaterializedStatementResults)
	return ret0
}

// GetMaterializedStatementResults indicates an expected call of GetMaterializedStatementResults.
func (mr *MockResultFetcherInterfaceMockRecorder) GetMaterializedStatementResults() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMaterializedStatementResults", reflect.TypeOf((*MockResultFetcherInterface)(nil).GetMaterializedStatementResults))
}

// GetRefreshState mocks base method.
func (m *MockResultFetcherInterface) GetRefreshState() types.RefreshState {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRefreshState")
	ret0, _ := ret[0].(types.RefreshState)
	return ret0
}

// GetRefreshState indicates an expected call of GetRefreshState.
func (mr *MockResultFetcherInterfaceMockRecorder) GetRefreshState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRefreshState", reflect.TypeOf((*MockResultFetcherInterface)(nil).GetRefreshState))
}

// GetStatement mocks base method.
func (m *MockResultFetcherInterface) GetStatement() types.ProcessedStatement {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatement")
	ret0, _ := ret[0].(types.ProcessedStatement)
	return ret0
}

// GetStatement indicates an expected call of GetStatement.
func (mr *MockResultFetcherInterfaceMockRecorder) GetStatement() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatement", reflect.TypeOf((*MockResultFetcherInterface)(nil).GetStatement))
}

// Init mocks base method.
func (m *MockResultFetcherInterface) Init(arg0 types.ProcessedStatement) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Init", arg0)
}

// Init indicates an expected call of Init.
func (mr *MockResultFetcherInterfaceMockRecorder) Init(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockResultFetcherInterface)(nil).Init), arg0)
}

// IsRefreshRunning mocks base method.
func (m *MockResultFetcherInterface) IsRefreshRunning() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRefreshRunning")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsRefreshRunning indicates an expected call of IsRefreshRunning.
func (mr *MockResultFetcherInterfaceMockRecorder) IsRefreshRunning() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRefreshRunning", reflect.TypeOf((*MockResultFetcherInterface)(nil).IsRefreshRunning))
}

// IsTableMode mocks base method.
func (m *MockResultFetcherInterface) IsTableMode() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsTableMode")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsTableMode indicates an expected call of IsTableMode.
func (mr *MockResultFetcherInterfaceMockRecorder) IsTableMode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsTableMode", reflect.TypeOf((*MockResultFetcherInterface)(nil).IsTableMode))
}

// SetRefreshCallback mocks base method.
func (m *MockResultFetcherInterface) SetRefreshCallback(arg0 func()) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetRefreshCallback", arg0)
}

// SetRefreshCallback indicates an expected call of SetRefreshCallback.
func (mr *MockResultFetcherInterfaceMockRecorder) SetRefreshCallback(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetRefreshCallback", reflect.TypeOf((*MockResultFetcherInterface)(nil).SetRefreshCallback), arg0)
}

// ToggleRefresh mocks base method.
func (m *MockResultFetcherInterface) ToggleRefresh() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ToggleRefresh")
}

// ToggleRefresh indicates an expected call of ToggleRefresh.
func (mr *MockResultFetcherInterfaceMockRecorder) ToggleRefresh() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToggleRefresh", reflect.TypeOf((*MockResultFetcherInterface)(nil).ToggleRefresh))
}

// ToggleTableMode mocks base method.
func (m *MockResultFetcherInterface) ToggleTableMode() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ToggleTableMode")
}

// ToggleTableMode indicates an expected call of ToggleTableMode.
func (mr *MockResultFetcherInterfaceMockRecorder) ToggleTableMode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToggleTableMode", reflect.TypeOf((*MockResultFetcherInterface)(nil).ToggleTableMode))
}
