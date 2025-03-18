package repository

import (
	"errors"
	"sync"
)

type URLRepository interface {
	AddURL(shortURL string, originalURL string) error
	GetURL(shortURL string) (string, error)
	GetExistingURL(originalURL string) string

	CreateURLRepository() (*URLRepositoryState, error)
	RestoreURLRepository(m *URLRepositoryState) error
}

type URLRepositoryImpl struct {
	repo sync.Map
}

func CreateInMemoryURLRepository() *URLRepositoryImpl {
	return &URLRepositoryImpl{
		repo: sync.Map{},
	}
}

func (ms *URLRepositoryImpl) AddURL(shortURL string, originalURL string) error {
	_, ok := ms.repo.LoadOrStore(shortURL, originalURL)
	if ok {
		return errors.New("this URL already exists")
	} else {
		return nil
	}
}

func (ms *URLRepositoryImpl) GetURL(shortURL string) (string, error) {

	value, ok := ms.repo.Load(shortURL)
	if !ok {
		return "", errors.New("this URL was not found")
	}
	return value.(string), nil
}

func (ms *URLRepositoryImpl) GetExistingURL(originalURL string) string {
	var foundKey string
	ms.repo.Range(func(key, value any) bool {
		if value == originalURL {
			foundKey = key.(string)
			return false
		}
		return true
	})
	return foundKey
}

func (ms *URLRepositoryImpl) deepCopyValues() map[string]string {
	copy := make(map[string]string)
	ms.repo.Range(func(key, value any) bool {
		copy[key.(string)] = value.(string)
		return true
	})

	return copy
}

func (ms *URLRepositoryImpl) CreateURLRepository() (*URLRepositoryState, error) {
	return CreateURLRepositoryState(ms.deepCopyValues()), nil
}

func (ms *URLRepositoryImpl) RestoreURLRepository(m *URLRepositoryState) error {
	for k, v := range m.GetURLRepositoryState() {
		ms.repo.Store(k, v)
	}
	return nil
}
