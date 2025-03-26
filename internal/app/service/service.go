package service

import (
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/config"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/repository"
	"github.com/SwarrenB/go-musthave-shortener-tpl/internal/app/urlgenerate"
)

type ServiceImpl interface {
	AddingURL(originalURL string, userID string) (string, error)
	GetOriginalURL(shortURL string) (string, error)
	GetURLByUserID(userID string) ([]repository.Record, error)
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

func (s *ShortenerService) AddingURL(originalURL string, userID string) (string, error) {
	shortURL := s.urlGenerator.GenerateURL(originalURL)
	newURL, err := s.repo.AddURL(shortURL, originalURL, userID)

	return newURL, err
}

func (s *ShortenerService) GetOriginalURL(shortURL string) (string, error) {

	originalURL, err := s.repo.GetURL(shortURL)

	if err != nil {
		return "", err
	}

	return originalURL, nil
}

func (s *ShortenerService) GetURLByUserID(userID string) ([]repository.Record, error) {
	return s.repo.GetURLByUserID(userID)
}
