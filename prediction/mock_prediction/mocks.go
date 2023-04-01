// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jamestunnell/marketanalysis/prediction (interfaces: Predictor)

// Package mock_prediction is a generated GoMock package.
package mock_prediction

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	prediction "github.com/jamestunnell/marketanalysis/prediction"
)

// MockPredictor is a mock of Predictor interface.
type MockPredictor struct {
	ctrl     *gomock.Controller
	recorder *MockPredictorMockRecorder
}

// MockPredictorMockRecorder is the mock recorder for MockPredictor.
type MockPredictorMockRecorder struct {
	mock *MockPredictor
}

// NewMockPredictor creates a new mock instance.
func NewMockPredictor(ctrl *gomock.Controller) *MockPredictor {
	mock := &MockPredictor{ctrl: ctrl}
	mock.recorder = &MockPredictorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPredictor) EXPECT() *MockPredictorMockRecorder {
	return m.recorder
}

// Train mocks base method.
func (m *MockPredictor) Train(arg0 []*prediction.TrainingElem) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Train", arg0)
}

// Train indicates an expected call of Train.
func (mr *MockPredictorMockRecorder) Train(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Train", reflect.TypeOf((*MockPredictor)(nil).Train), arg0)
}