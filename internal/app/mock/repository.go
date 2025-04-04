// Code generated by MockGen. DO NOT EDIT.
// Source: C:\Users\dasa1020\go\go-musthave-shortener-tpl\internal\app\repository\repository.go

// Package mock_repository is a generated GoMock package.
package mock

import (
	reflect "reflect"

	repository "github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/repository"
	gomock "github.com/golang/mock/gomock"
)

// MockURLRepository is a mock of URLRepository interface.
type MockURLRepository struct {
	ctrl     *gomock.Controller
	recorder *MockURLRepositoryMockRecorder
}

// MockURLRepositoryMockRecorder is the mock recorder for MockURLRepository.
type MockURLRepositoryMockRecorder struct {
	mock *MockURLRepository
}

// NewMockURLRepository creates a new mock instance.
func NewMockURLRepository(ctrl *gomock.Controller) *MockURLRepository {
	mock := &MockURLRepository{ctrl: ctrl}
	mock.recorder = &MockURLRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLRepository) EXPECT() *MockURLRepositoryMockRecorder {
	return m.recorder
}

// AddURL mocks base method.
func (m *MockURLRepository) AddURL(shortURL, originalURL string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddURL", shortURL, originalURL)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddURL indicates an expected call of AddURL.
func (mr *MockURLRepositoryMockRecorder) AddURL(shortURL, originalURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddURL", reflect.TypeOf((*MockURLRepository)(nil).AddURL), shortURL, originalURL)
}

// CreateURLRepository mocks base method.
func (m *MockURLRepository) CreateURLRepository() (*repository.URLRepositoryState, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateURLRepository")
	ret0, _ := ret[0].(*repository.URLRepositoryState)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateURLRepository indicates an expected call of CreateURLRepository.
func (mr *MockURLRepositoryMockRecorder) CreateURLRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateURLRepository", reflect.TypeOf((*MockURLRepository)(nil).CreateURLRepository))
}

// GetURL mocks base method.
func (m *MockURLRepository) GetURL(shortURL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetURL", shortURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetURL indicates an expected call of GetURL.
func (mr *MockURLRepositoryMockRecorder) GetURL(shortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURL", reflect.TypeOf((*MockURLRepository)(nil).GetURL), shortURL)
}

// RestoreURLRepository mocks base method.
func (m_2 *MockURLRepository) RestoreURLRepository(m *repository.URLRepositoryState) error {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "RestoreURLRepository", m)
	ret0, _ := ret[0].(error)
	return ret0
}

// RestoreURLRepository indicates an expected call of RestoreURLRepository.
func (mr *MockURLRepositoryMockRecorder) RestoreURLRepository(m interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RestoreURLRepository", reflect.TypeOf((*MockURLRepository)(nil).RestoreURLRepository), m)
}
