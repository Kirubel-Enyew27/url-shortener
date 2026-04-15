package services

import (
	"errors"
	"url-shortener/internal/models"
	"url-shortener/internal/storage"
	"url-shortener/pkg/utils"
)

var (
	ErrInvalidURL = errors.New("invalid URL")
	ErrNotFound   = errors.New("URL not found")
)

type URLService struct {
	store storage.StoreInterface
}

func New(store storage.StoreInterface) *URLService {
	return &URLService{store: store}
}

func (s *URLService) Shorten(longURL string) (*models.URL, error) {
	if longURL == "" {
		return nil, ErrInvalidURL
	}

	if existing, found := s.store.GetByLongURL(longURL); found {
		return existing, nil
	}

	code := utils.GenerateCode(6)
	return s.store.Set(code, longURL), nil
}

func (s *URLService) Get(code string) (*models.URL, error) {
	url, found := s.store.Get(code)
	if !found {
		return nil, ErrNotFound
	}

	return url, nil
}

func (s *URLService) GetAll() []*models.URL {
	return s.store.GetAll()
}

func (s *URLService) Resolve(code string) (string, error) {
	url, found := s.store.Get(code)
	if !found {
		return "", ErrNotFound
	}

	s.store.IncrementClicks(code)
	return url.LongURL, nil
}
