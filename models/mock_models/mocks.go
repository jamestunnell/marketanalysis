// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jamestunnell/marketanalysis/models (interfaces: Block,Param,Constraint,Input,Output)

// Package mock_models is a generated GoMock package.
package mock_models

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
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
func (m *MockBlock) GetInputs() models.Inputs {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInputs")
	ret0, _ := ret[0].(models.Inputs)
	return ret0
}

// GetInputs indicates an expected call of GetInputs.
func (mr *MockBlockMockRecorder) GetInputs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInputs", reflect.TypeOf((*MockBlock)(nil).GetInputs))
}

// GetOutputs mocks base method.
func (m *MockBlock) GetOutputs() models.Outputs {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOutputs")
	ret0, _ := ret[0].(models.Outputs)
	return ret0
}

// GetOutputs indicates an expected call of GetOutputs.
func (mr *MockBlockMockRecorder) GetOutputs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOutputs", reflect.TypeOf((*MockBlock)(nil).GetOutputs))
}

// GetParams mocks base method.
func (m *MockBlock) GetParams() models.Params {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetParams")
	ret0, _ := ret[0].(models.Params)
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

// Constraints mocks base method.
func (m *MockParam) Constraints() []models.Constraint {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Constraints")
	ret0, _ := ret[0].([]models.Constraint)
	return ret0
}

// Constraints indicates an expected call of Constraints.
func (mr *MockParamMockRecorder) Constraints() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Constraints", reflect.TypeOf((*MockParam)(nil).Constraints))
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

// LoadVal mocks base method.
func (m *MockParam) LoadVal(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoadVal", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// LoadVal indicates an expected call of LoadVal.
func (mr *MockParamMockRecorder) LoadVal(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoadVal", reflect.TypeOf((*MockParam)(nil).LoadVal), arg0)
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

// StoreVal mocks base method.
func (m *MockParam) StoreVal() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreVal")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StoreVal indicates an expected call of StoreVal.
func (mr *MockParamMockRecorder) StoreVal() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreVal", reflect.TypeOf((*MockParam)(nil).StoreVal))
}

// Type mocks base method.
func (m *MockParam) Type() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(string)
	return ret0
}

// Type indicates an expected call of Type.
func (mr *MockParamMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockParam)(nil).Type))
}

// MockConstraint is a mock of Constraint interface.
type MockConstraint struct {
	ctrl     *gomock.Controller
	recorder *MockConstraintMockRecorder
}

// MockConstraintMockRecorder is the mock recorder for MockConstraint.
type MockConstraintMockRecorder struct {
	mock *MockConstraint
}

// NewMockConstraint creates a new mock instance.
func NewMockConstraint(ctrl *gomock.Controller) *MockConstraint {
	mock := &MockConstraint{ctrl: ctrl}
	mock.recorder = &MockConstraintMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConstraint) EXPECT() *MockConstraintMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockConstraint) Check(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockConstraintMockRecorder) Check(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockConstraint)(nil).Check), arg0)
}

// Type mocks base method.
func (m *MockConstraint) Type() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(string)
	return ret0
}

// Type indicates an expected call of Type.
func (mr *MockConstraintMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockConstraint)(nil).Type))
}

// ValueBounds mocks base method.
func (m *MockConstraint) ValueBounds() []interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValueBounds")
	ret0, _ := ret[0].([]interface{})
	return ret0
}

// ValueBounds indicates an expected call of ValueBounds.
func (mr *MockConstraintMockRecorder) ValueBounds() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValueBounds", reflect.TypeOf((*MockConstraint)(nil).ValueBounds))
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
func (m *MockOutput) Connect(arg0 models.Input) error {
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
