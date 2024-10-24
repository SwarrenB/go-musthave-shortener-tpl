package repository

import (
	"errors"
	"sync"
)

type URLRepository interface {
	AddURL(shortURL string, originalURL string) error
	GetURL(shortURL string) (string, error)
}

type URLRepositoryImpl struct {
	sync.RWMutex
	URLRepository
	values map[string]string
}

func CreateInMemoryURLRepository() *URLRepositoryImpl {
	return &URLRepositoryImpl{
		values: map[string]string{},
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
