package repository

import (
	"errors"
	"sync"
)

type URLRepository interface {
	AddURL(shortURL string, originalURL string) error
	GetURL(shortURL string) (string, error)

	CreateURLRepository() (*URLRepositoryState, error)
	RestoreURLRepository(m *URLRepositoryState) error
}

type URLRepositoryImpl struct {
	sync.RWMutex
	values map[string]string
}

func CreateInMemoryURLRepository() *URLRepositoryImpl {
	return &URLRepositoryImpl{
		RWMutex: sync.RWMutex{},
		values:  map[string]string{},
	}
}

func (ms *URLRepositoryImpl) AddURL(shortURL string, originalURL string) error {
	ms.Lock()
	defer ms.Unlock()
	_, ok := ms.values[shortURL]
	if ok {
		return errors.New("this URL exists in vocabulary")
	}
	ms.values[shortURL] = originalURL
	return nil
}

func (ms *URLRepositoryImpl) GetURL(shortURL string) (string, error) {
	ms.Lock()
	defer ms.Unlock()

	value, ok := ms.values[shortURL]
	if !ok {
		return "", errors.New("this URL was not found")

	}
	return value, nil
}

func (ms *URLRepositoryImpl) deepCopyValues() map[string]string {
	copy := make(map[string]string)
	for k, v := range ms.values {
		copy[k] = v
	}

	return copy
}

func (ms *URLRepositoryImpl) CreateURLRepository() (*URLRepositoryState, error) {
	ms.Lock()
	defer ms.Unlock()

	return CreateURLRepositoryState(ms.deepCopyValues()), nil
}

func (ms *URLRepositoryImpl) RestoreURLRepository(m *URLRepositoryState) error {
	ms.Lock()
	defer ms.Unlock()

	copy := make(map[string]string)
	for k, v := range m.GetURLRepositoryState() {
		copy[k] = v
	}

	ms.values = copy
	return nil
}
