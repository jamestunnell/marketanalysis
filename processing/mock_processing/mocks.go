// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jamestunnell/marketanalysis/processing (interfaces: Processor,Source,Element)

// Package mock_processing is a generated GoMock package.
package mock_processing

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/jamestunnell/marketanalysis/models"
)

// MockProcessor is a mock of Processor interface.
type MockProcessor struct {
	ctrl     *gomock.Controller
	recorder *MockProcessorMockRecorder
}

// MockProcessorMockRecorder is the mock recorder for MockProcessor.
type MockProcessorMockRecorder struct {
	mock *MockProcessor
}

// NewMockProcessor creates a new mock instance.
func NewMockProcessor(ctrl *gomock.Controller) *MockProcessor {
	mock := &MockProcessor{ctrl: ctrl}
	mock.recorder = &MockProcessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProcessor) EXPECT() *MockProcessorMockRecorder {
	return m.recorder
}

// Initialize mocks base method.
func (m *MockProcessor) Initialize() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize")
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize.
func (mr *MockProcessorMockRecorder) Initialize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockProcessor)(nil).Initialize))
}

// Output mocks base method.
func (m *MockProcessor) Output() float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Output")
	ret0, _ := ret[0].(float64)
	return ret0
}

// Output indicates an expected call of Output.
func (mr *MockProcessorMockRecorder) Output() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Output", reflect.TypeOf((*MockProcessor)(nil).Output))
}

// Params mocks base method.
func (m *MockProcessor) Params() models.Params {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Params")
	ret0, _ := ret[0].(models.Params)
	return ret0
}

// Params indicates an expected call of Params.
func (mr *MockProcessorMockRecorder) Params() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Params", reflect.TypeOf((*MockProcessor)(nil).Params))
}

// Type mocks base method.
func (m *MockProcessor) Type() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(string)
	return ret0
}

// Type indicates an expected call of Type.
func (mr *MockProcessorMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockProcessor)(nil).Type))
}

// Update mocks base method.
func (m *MockProcessor) Update(arg0 float64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Update", arg0)
}

// Update indicates an expected call of Update.
func (mr *MockProcessorMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockProcessor)(nil).Update), arg0)
}

// WarmUp mocks base method.
func (m *MockProcessor) WarmUp(arg0 []float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WarmUp", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WarmUp indicates an expected call of WarmUp.
func (mr *MockProcessorMockRecorder) WarmUp(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WarmUp", reflect.TypeOf((*MockProcessor)(nil).WarmUp), arg0)
}

// WarmupPeriod mocks base method.
func (m *MockProcessor) WarmupPeriod() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WarmupPeriod")
	ret0, _ := ret[0].(int)
	return ret0
}

// WarmupPeriod indicates an expected call of WarmupPeriod.
func (mr *MockProcessorMockRecorder) WarmupPeriod() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WarmupPeriod", reflect.TypeOf((*MockProcessor)(nil).WarmupPeriod))
}

// MockSource is a mock of Source interface.
type MockSource struct {
	ctrl     *gomock.Controller
	recorder *MockSourceMockRecorder
}

// MockSourceMockRecorder is the mock recorder for MockSource.
type MockSourceMockRecorder struct {
	mock *MockSource
}

// NewMockSource creates a new mock instance.
func NewMockSource(ctrl *gomock.Controller) *MockSource {
	mock := &MockSource{ctrl: ctrl}
	mock.recorder = &MockSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSource) EXPECT() *MockSourceMockRecorder {
	return m.recorder
}

// Initialize mocks base method.
func (m *MockSource) Initialize() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize")
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize.
func (mr *MockSourceMockRecorder) Initialize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockSource)(nil).Initialize))
}

// Output mocks base method.
func (m *MockSource) Output() float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Output")
	ret0, _ := ret[0].(float64)
	return ret0
}

// Output indicates an expected call of Output.
func (mr *MockSourceMockRecorder) Output() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Output", reflect.TypeOf((*MockSource)(nil).Output))
}

// Params mocks base method.
func (m *MockSource) Params() models.Params {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Params")
	ret0, _ := ret[0].(models.Params)
	return ret0
}

// Params indicates an expected call of Params.
func (mr *MockSourceMockRecorder) Params() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Params", reflect.TypeOf((*MockSource)(nil).Params))
}

// Type mocks base method.
func (m *MockSource) Type() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(string)
	return ret0
}

// Type indicates an expected call of Type.
func (mr *MockSourceMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockSource)(nil).Type))
}

// Update mocks base method.
func (m *MockSource) Update(arg0 *models.Bar) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Update", arg0)
}

// Update indicates an expected call of Update.
func (mr *MockSourceMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSource)(nil).Update), arg0)
}

// WarmUp mocks base method.
func (m *MockSource) WarmUp(arg0 models.Bars) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WarmUp", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WarmUp indicates an expected call of WarmUp.
func (mr *MockSourceMockRecorder) WarmUp(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WarmUp", reflect.TypeOf((*MockSource)(nil).WarmUp), arg0)
}

// WarmupPeriod mocks base method.
func (m *MockSource) WarmupPeriod() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WarmupPeriod")
	ret0, _ := ret[0].(int)
	return ret0
}

// WarmupPeriod indicates an expected call of WarmupPeriod.
func (mr *MockSourceMockRecorder) WarmupPeriod() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WarmupPeriod", reflect.TypeOf((*MockSource)(nil).WarmupPeriod))
}

// MockElement is a mock of Element interface.
type MockElement struct {
	ctrl     *gomock.Controller
	recorder *MockElementMockRecorder
}

// MockElementMockRecorder is the mock recorder for MockElement.
type MockElementMockRecorder struct {
	mock *MockElement
}

// NewMockElement creates a new mock instance.
func NewMockElement(ctrl *gomock.Controller) *MockElement {
	mock := &MockElement{ctrl: ctrl}
	mock.recorder = &MockElementMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockElement) EXPECT() *MockElementMockRecorder {
	return m.recorder
}

// Initialize mocks base method.
func (m *MockElement) Initialize() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Initialize")
	ret0, _ := ret[0].(error)
	return ret0
}

// Initialize indicates an expected call of Initialize.
func (mr *MockElementMockRecorder) Initialize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Initialize", reflect.TypeOf((*MockElement)(nil).Initialize))
}

// Output mocks base method.
func (m *MockElement) Output() float64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Output")
	ret0, _ := ret[0].(float64)
	return ret0
}

// Output indicates an expected call of Output.
func (mr *MockElementMockRecorder) Output() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Output", reflect.TypeOf((*MockElement)(nil).Output))
}

// Params mocks base method.
func (m *MockElement) Params() models.Params {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Params")
	ret0, _ := ret[0].(models.Params)
	return ret0
}

// Params indicates an expected call of Params.
func (mr *MockElementMockRecorder) Params() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Params", reflect.TypeOf((*MockElement)(nil).Params))
}

// Type mocks base method.
func (m *MockElement) Type() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Type")
	ret0, _ := ret[0].(string)
	return ret0
}

// Type indicates an expected call of Type.
func (mr *MockElementMockRecorder) Type() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Type", reflect.TypeOf((*MockElement)(nil).Type))
}

// WarmupPeriod mocks base method.
func (m *MockElement) WarmupPeriod() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WarmupPeriod")
	ret0, _ := ret[0].(int)
	return ret0
}

// WarmupPeriod indicates an expected call of WarmupPeriod.
func (mr *MockElementMockRecorder) WarmupPeriod() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WarmupPeriod", reflect.TypeOf((*MockElement)(nil).WarmupPeriod))
}
