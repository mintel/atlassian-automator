// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mintel/atlassian-automator/pkg/jiraclient (interfaces: Client)

// Package mock_jiraclient is a generated GoMock package.
package mock_jiraclient

import (
	bytes "bytes"
	context "context"
	io "io"
	http "net/http"
	url "net/url"
	reflect "reflect"

	jira "github.com/andygrunwald/go-jira"
	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockClient) Do(arg0 *http.Request, arg1 interface{}) (*jira.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", arg0, arg1)
	ret0, _ := ret[0].(*jira.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockClientMockRecorder) Do(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockClient)(nil).Do), arg0, arg1)
}

// GetBaseURL mocks base method.
func (m *MockClient) GetBaseURL() url.URL {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBaseURL")
	ret0, _ := ret[0].(url.URL)
	return ret0
}

// GetBaseURL indicates an expected call of GetBaseURL.
func (mr *MockClientMockRecorder) GetBaseURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBaseURL", reflect.TypeOf((*MockClient)(nil).GetBaseURL))
}

// NewMultiPartRequest mocks base method.
func (m *MockClient) NewMultiPartRequest(arg0, arg1 string, arg2 *bytes.Buffer) (*http.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewMultiPartRequest", arg0, arg1, arg2)
	ret0, _ := ret[0].(*http.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewMultiPartRequest indicates an expected call of NewMultiPartRequest.
func (mr *MockClientMockRecorder) NewMultiPartRequest(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewMultiPartRequest", reflect.TypeOf((*MockClient)(nil).NewMultiPartRequest), arg0, arg1, arg2)
}

// NewMultiPartRequestWithContext mocks base method.
func (m *MockClient) NewMultiPartRequestWithContext(arg0 context.Context, arg1, arg2 string, arg3 *bytes.Buffer) (*http.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewMultiPartRequestWithContext", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*http.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewMultiPartRequestWithContext indicates an expected call of NewMultiPartRequestWithContext.
func (mr *MockClientMockRecorder) NewMultiPartRequestWithContext(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewMultiPartRequestWithContext", reflect.TypeOf((*MockClient)(nil).NewMultiPartRequestWithContext), arg0, arg1, arg2, arg3)
}

// NewRawRequest mocks base method.
func (m *MockClient) NewRawRequest(arg0, arg1 string, arg2 io.Reader) (*http.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRawRequest", arg0, arg1, arg2)
	ret0, _ := ret[0].(*http.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRawRequest indicates an expected call of NewRawRequest.
func (mr *MockClientMockRecorder) NewRawRequest(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRawRequest", reflect.TypeOf((*MockClient)(nil).NewRawRequest), arg0, arg1, arg2)
}

// NewRawRequestWithContext mocks base method.
func (m *MockClient) NewRawRequestWithContext(arg0 context.Context, arg1, arg2 string, arg3 io.Reader) (*http.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRawRequestWithContext", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*http.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRawRequestWithContext indicates an expected call of NewRawRequestWithContext.
func (mr *MockClientMockRecorder) NewRawRequestWithContext(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRawRequestWithContext", reflect.TypeOf((*MockClient)(nil).NewRawRequestWithContext), arg0, arg1, arg2, arg3)
}

// NewRequest mocks base method.
func (m *MockClient) NewRequest(arg0, arg1 string, arg2 interface{}) (*http.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRequest", arg0, arg1, arg2)
	ret0, _ := ret[0].(*http.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRequest indicates an expected call of NewRequest.
func (mr *MockClientMockRecorder) NewRequest(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRequest", reflect.TypeOf((*MockClient)(nil).NewRequest), arg0, arg1, arg2)
}

// NewRequestWithContext mocks base method.
func (m *MockClient) NewRequestWithContext(arg0 context.Context, arg1, arg2 string, arg3 interface{}) (*http.Request, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewRequestWithContext", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*http.Request)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewRequestWithContext indicates an expected call of NewRequestWithContext.
func (mr *MockClientMockRecorder) NewRequestWithContext(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewRequestWithContext", reflect.TypeOf((*MockClient)(nil).NewRequestWithContext), arg0, arg1, arg2, arg3)
}