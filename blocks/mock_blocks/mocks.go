// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jamestunnell/marketanalysis/blocks (interfaces: Block,Param,Input,Output)

// Package mock_blocks is a generated GoMock package.
package mock_blocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	blocks "github.com/jamestunnell/marketanalysis/blocks"
	models "github.com/jamestunnell/marketanalysis/models"
)

// MockBlock is a mock of Block interface.
type MockBlock struct {
	ctrl     *gomock.Controller
	recorder *MockBlockMockRecorder
}

// MockBlockMockRecorder is the mock recorder for MockBlock.
type MockBlockMockRecorder struct {
	mock *MockBlock
}

// NewMockBlock creates a new mock instance.
func NewMockBlock(ctrl *gomock.Controller) *MockBlock {
	mock := &MockBlock{ctrl: ctrl}
	mock.recorder = &MockBlockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlock) EXPECT() *MockBlockMockRecorder {
	return m.recorder
}

// GetDescription mocks base method.
func (m *MockBlock) GetDescription() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDescription")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetDescription indicates an expected call of GetDescription.
func (mr *MockBlockMockRecorder) GetDescription() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDescription", reflect.TypeOf((*MockBlock)(nil).GetDescription))
}

// GetInputs mocks base method.
func (m *MockBlock) GetInputs() blocks.Inputs {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInputs")
	ret0, _ := ret[0].(blocks.Inputs)
	return ret0
}

// GetInputs indicates an expected call of GetInputs.
func (mr *MockBlockMockRecorder) GetInputs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInputs", reflect.TypeOf((*MockBlock)(nil).GetInputs))
}

// GetOutputs mocks base method.
func (m *MockBlock) GetOutputs() blocks.Outputs {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOutputs")
	ret0, _ := ret[0].(blocks.Outputs)
	return ret0
}

// GetOutputs indicates an expected call of GetOutputs.
func (mr *MockBlockMockRecorder) GetOutputs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOutputs", reflect.TypeOf((*MockBlock)(nil).GetOutputs))
}

// GetParams mocks base method.
func (m *MockBlock) GetParams() blocks.Params {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetParams")
	ret0, _ := ret[0].(blocks.Params)
	return ret0
}

// GetParams indicates an expected call of GetParams.
func (mr *MockBlockMockRecorder) GetParams() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetParams", reflect.TypeOf((*MockBlock)(nil).GetParams))
}

// GetType mocks base method.
func (m *MockBlock) GetType() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetType")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetType indicates an expected call of GetType.
func (mr *MockBlockMockRecorder) GetType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetType", reflect.TypeOf((*MockBlock)(nil).GetType))
}

// GetWarmupPeriod mocks base method.
func (m *MockBlock) GetWarmupPeriod() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWarmupPeriod")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetWarmupPeriod indicates an expected call of GetWarmupPeriod.
func (mr *MockBlockMockRecorder) GetWarmupPeriod() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWarmupPeriod", reflect.TypeOf((*MockBlock)(nil).GetWarmupPeriod))
}

// Init mocks base method.
func (m *MockBlock) Init() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init")
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockBlockMockRecorder) Init() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockBlock)(nil).Init))
}

// IsWarm mocks base method.
func (m *MockBlock) IsWarm() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsWarm")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsWarm indicates an expected call of IsWarm.
func (mr *MockBlockMockRecorder) IsWarm() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsWarm", reflect.TypeOf((*MockBlock)(nil).IsWarm))
}

// Update mocks base method.
func (m *MockBlock) Update(arg0 *models.Bar) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Update", arg0)
}

// Update indicates an expected call of Update.
func (mr *MockBlockMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockBlock)(nil).Update), arg0)
}

// MockParam is a mock of Param interface.
type MockParam struct {
	ctrl     *gomock.Controller
	recorder *MockParamMockRecorder
}

// MockParamMockRecorder is the mock recorder for MockParam.
type MockParamMockRecorder struct {
	mock *MockParam
}

// NewMockParam creates a new mock instance.
func NewMockParam(ctrl *gomock.Controller) *MockParam {
	mock := &MockParam{ctrl: ctrl}
	mock.recorder = &MockParamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockParam) EXPECT() *MockParamMockRecorder {
	return m.recorder
}

// GetDefault mocks base method.
func (m *MockParam) GetDefault() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDefault")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// GetDefault indicates an expected call of GetDefault.
func (mr *MockParamMockRecorder) GetDefault() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDefault", reflect.TypeOf((*MockParam)(nil).GetDefault))
}

// GetSchema mocks base method.
func (m *MockParam) GetSchema() map[string]interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSchema")
	ret0, _ := ret[0].(map[string]interface{})
	return ret0
}

// GetSchema indicates an expected call of GetSchema.
func (mr *MockParamMockRecorder) GetSchema() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSchema", reflect.TypeOf((*MockParam)(nil).GetSchema))
}

// GetVal mocks base method.
func (m *MockParam) GetVal() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVal")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// GetVal indicates an expected call of GetVal.
func (mr *MockParamMockRecorder) GetVal() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVal", reflect.TypeOf((*MockParam)(nil).GetVal))
}

// SetVal mocks base method.
func (m *MockParam) SetVal(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetVal", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetVal indicates an expected call of SetVal.
func (mr *MockParamMockRecorder) SetVal(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetVal", reflect.TypeOf((*MockParam)(nil).SetVal), arg0)
}

// MockInput is a mock of Input interface.
type MockInput struct {
	ctrl     *gomock.Controller
	recorder *MockInputMockRecorder
}

// MockInputMockRecorder is the mock recorder for MockInput.
type MockInputMockRecorder struct {
	mock *MockInput
}

// NewMockInput creates a new mock instance.
func NewMockInput(ctrl *gomock.Controller) *MockInput {
	mock := &MockInput{ctrl: ctrl}
	mock.recorder = &MockInputMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInput) EXPECT() *MockInputMockRecorder {
	return m.recorder
}

// Connect mocks base method.
func (m *MockInput) Connect(arg0 blocks.Output) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect.
func (mr *MockInputMockRecorder) Connect(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockInput)(nil).Connect), arg0)
}

// GetType mocks base method.
func (m *MockInput) GetType() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetType")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetType indicates an expected call of GetType.
func (mr *MockInputMockRecorder) GetType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetType", reflect.TypeOf((*MockInput)(nil).GetType))
}

// IsConnected mocks base method.
func (m *MockInput) IsConnected() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsConnected")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsConnected indicates an expected call of IsConnected.
func (mr *MockInputMockRecorder) IsConnected() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsConnected", reflect.TypeOf((*MockInput)(nil).IsConnected))
}

// IsValueSet mocks base method.
func (m *MockInput) IsValueSet() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsValueSet")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsValueSet indicates an expected call of IsValueSet.
func (mr *MockInputMockRecorder) IsValueSet() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsValueSet", reflect.TypeOf((*MockInput)(nil).IsValueSet))
}

// MockOutput is a mock of Output interface.
type MockOutput struct {
	ctrl     *gomock.Controller
	recorder *MockOutputMockRecorder
}

// MockOutputMockRecorder is the mock recorder for MockOutput.
type MockOutputMockRecorder struct {
	mock *MockOutput
}

// NewMockOutput creates a new mock instance.
func NewMockOutput(ctrl *gomock.Controller) *MockOutput {
	mock := &MockOutput{ctrl: ctrl}
	mock.recorder = &MockOutputMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOutput) EXPECT() *MockOutputMockRecorder {
	return m.recorder
}

// Connect mocks base method.
func (m *MockOutput) Connect(arg0 blocks.Input) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect.
func (mr *MockOutputMockRecorder) Connect(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockOutput)(nil).Connect), arg0)
}

// DisconnectAll mocks base method.
func (m *MockOutput) DisconnectAll() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "DisconnectAll")
}

// DisconnectAll indicates an expected call of DisconnectAll.
func (mr *MockOutputMockRecorder) DisconnectAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DisconnectAll", reflect.TypeOf((*MockOutput)(nil).DisconnectAll))
}

// GetConnected mocks base method.
func (m *MockOutput) GetConnected() []blocks.Input {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnected")
	ret0, _ := ret[0].([]blocks.Input)
	return ret0
}

// GetConnected indicates an expected call of GetConnected.
func (mr *MockOutputMockRecorder) GetConnected() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnected", reflect.TypeOf((*MockOutput)(nil).GetConnected))
}

// GetType mocks base method.
func (m *MockOutput) GetType() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetType")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetType indicates an expected call of GetType.
func (mr *MockOutputMockRecorder) GetType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetType", reflect.TypeOf((*MockOutput)(nil).GetType))
}

// IsConnected mocks base method.
func (m *MockOutput) IsConnected() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsConnected")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsConnected indicates an expected call of IsConnected.
func (mr *MockOutputMockRecorder) IsConnected() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsConnected", reflect.TypeOf((*MockOutput)(nil).IsConnected))
}

// IsValueSet mocks base method.
func (m *MockOutput) IsValueSet() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsValueSet")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsValueSet indicates an expected call of IsValueSet.
func (mr *MockOutputMockRecorder) IsValueSet() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsValueSet", reflect.TypeOf((*MockOutput)(nil).IsValueSet))
}