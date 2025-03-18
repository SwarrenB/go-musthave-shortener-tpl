package service

import (
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/config"
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
	config       *config.Config
}

func CreateShortenerService(
	repo repository.URLRepository,
	gen urlgenerate.URLGenerator,
	config *config.Config,
) *ShortenerService {
	return &ShortenerService{
		repo:         repo,
		urlGenerator: gen,
		config:       config,
	}
}

func (s *ShortenerService) AddingURL(originalURL string) (string, error) {
	var shortURL string
	var err error
	shortURL = s.urlGenerator.GenerateURL(originalURL)

	err = s.repo.AddURL(shortURL, originalURL)

	if err != nil {
		existingURL := s.repo.GetExistingURL(originalURL)
		return existingURL, err
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
