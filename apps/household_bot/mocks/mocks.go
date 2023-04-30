// Code generated by MockGen. DO NOT EDIT.
// Source: internal/telegram/handler/handler.go

// Package mock_handler is a generated GoMock package.
package mock_handler

import (
	context "context"
	reflect "reflect"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gomock "github.com/golang/mock/gomock"
)

// MockRateProvider is a mock of RateProvider interface.
type MockRateProvider struct {
	ctrl     *gomock.Controller
	recorder *MockRateProviderMockRecorder
}

// MockRateProviderMockRecorder is the mock recorder for MockRateProvider.
type MockRateProviderMockRecorder struct {
	mock *MockRateProvider
}

// NewMockRateProvider creates a new mock instance.
func NewMockRateProvider(ctrl *gomock.Controller) *MockRateProvider {
	mock := &MockRateProvider{ctrl: ctrl}
	mock.recorder = &MockRateProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRateProvider) EXPECT() *MockRateProviderMockRecorder {
	return m.recorder
}

// GetYuanRate mocks base method.
func (m *MockRateProvider) GetYuanRate(ctx context.Context) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetYuanRate", ctx)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetYuanRate indicates an expected call of GetYuanRate.
func (mr *MockRateProviderMockRecorder) GetYuanRate(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetYuanRate", reflect.TypeOf((*MockRateProvider)(nil).GetYuanRate), ctx)
}

// MockBot is a mock of Bot interface.
type MockBot struct {
	ctrl     *gomock.Controller
	recorder *MockBotMockRecorder
}

// MockBotMockRecorder is the mock recorder for MockBot.
type MockBotMockRecorder struct {
	mock *MockBot
}

// NewMockBot creates a new mock instance.
func NewMockBot(ctrl *gomock.Controller) *MockBot {
	mock := &MockBot{ctrl: ctrl}
	mock.recorder = &MockBotMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBot) EXPECT() *MockBotMockRecorder {
	return m.recorder
}

// CleanRequest mocks base method.
func (m *MockBot) CleanRequest(c tgbotapi.Chattable) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CleanRequest", c)
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanRequest indicates an expected call of CleanRequest.
func (mr *MockBotMockRecorder) CleanRequest(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanRequest", reflect.TypeOf((*MockBot)(nil).CleanRequest), c)
}

// Send mocks base method.
func (m *MockBot) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", c)
	ret0, _ := ret[0].(tgbotapi.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Send indicates an expected call of Send.
func (mr *MockBotMockRecorder) Send(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockBot)(nil).Send), c)
}

// SendMediaGroup mocks base method.
func (m *MockBot) SendMediaGroup(c tgbotapi.MediaGroupConfig) ([]tgbotapi.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMediaGroup", c)
	ret0, _ := ret[0].([]tgbotapi.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SendMediaGroup indicates an expected call of SendMediaGroup.
func (mr *MockBotMockRecorder) SendMediaGroup(c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMediaGroup", reflect.TypeOf((*MockBot)(nil).SendMediaGroup), c)
}