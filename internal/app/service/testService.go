package service

import "github.com/stretchr/testify/mock"

type TestService struct {
	ServiceImpl
	mock.Mock
}

func CreateTestService() *TestService {
	return &TestService{}
}

func (m *TestService) AddingURL(originalURL string) (string, error) {
	args := m.Called(originalURL)
	return args.String(0), args.Error(1)
}

func (m *TestService) GetOriginalURL(shortURL string) (string, error) {
	args := m.Called(shortURL)
	return args.String(0), args.Error(1)
}
