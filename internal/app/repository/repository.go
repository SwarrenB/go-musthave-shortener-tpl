package repository

import (
	"errors"
	"sync"
)

type URLRepository interface {
	AddURL(shortURL string, originalURL string, userID string) (string, error)
	GetURL(shortURL string) (string, error)

	CreateURLRepository() (*URLRepositoryState, error)
	RestoreURLRepository(m *URLRepositoryState) error
	GetURLByUserID(userID string) ([]Record, error)
}

type URLRepositoryImpl struct {
	repo sync.Map
}

func CreateInMemoryURLRepository() *URLRepositoryImpl {
	return &URLRepositoryImpl{
		repo: sync.Map{},
	}
}

func (ms *URLRepositoryImpl) AddURL(shortURL string, originalURL string, userID string) (string, error) {

	_, ok := ms.repo.Load(shortURL)
	if ok {
		var existingURL string
		ms.repo.Range(func(key, value any) bool {
			if value.(Record).OriginalURL == originalURL {
				existingURL = key.(string)
				return false
			}
			return true
		})
		return existingURL, errors.New("this URL already exists")
	} else {
		ms.repo.Store(
			shortURL,
			Record{
				ID:          0,
				ShortURL:    shortURL,
				OriginalURL: originalURL,
				UserID:      userID,
			})
		return shortURL, nil
	}
}

func (ms *URLRepositoryImpl) GetURL(shortURL string) (string, error) {

	value, ok := ms.repo.Load(shortURL)
	if !ok {
		return "", errors.New("this URL was not found")
	}
	return value.(Record).OriginalURL, nil
}

func (ms *URLRepositoryImpl) deepCopyValues() map[string]Record {
	copy := make(map[string]Record)
	ms.repo.Range(func(key, value any) bool {
		copy[key.(string)] = value.(Record)
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

func (ms *URLRepositoryImpl) GetURLByUserID(userID string) ([]Record, error) {
	var results []Record
	ms.repo.Range(func(key, value any) bool {
		if value.(Record).UserID == userID {
			results = append(results, value.(Record))
			return false
		}
		return true
	})
	return results, nil
}
