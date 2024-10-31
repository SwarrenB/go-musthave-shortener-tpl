package service

import (
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/repository"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/urlgenerate"
)

type ServiceImpl interface {
	AddingURL(originalURL string) (string, error)
	GetOriginalURL(shortURL string) (string, error)
}

type ShortenerService struct {
	repo         repository.URLRepository
	urlGenerator urlgenerate.URLGenerator
}

func CreateShortenerService() *ShortenerService {
	return &ShortenerService{
		repo:         repository.CreateInMemoryURLRepository(),
		urlGenerator: urlgenerate.CreateURLGenerator(),
	}
}

func (s *ShortenerService) AddingURL(originalURL string) (string, error) {
	var shortURL string
	var err error
	shortURL = s.urlGenerator.GenerateURL(originalURL)

	err = s.repo.AddURL(shortURL, originalURL)

	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (s *ShortenerService) GetOriginalURL(shortURL string) (string, error) {

	originalURL, err := s.repo.GetURL(shortURL)

	if err != nil {
		return "", err
	}

	return originalURL, nil
}
