package services

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"url-shortener/internal/models"
	"url-shortener/internal/storage"
	"url-shortener/pkg/utils"
)

const (
	defaultCodeLength = 6
	maxCodeRetries    = 16
)

var (
	ErrInvalidURL         = errors.New("invalid URL")
	ErrNotFound           = errors.New("URL not found")
	ErrCouldNotCreateCode = errors.New("could not generate unique short code")
)

type URLService struct {
	store        storage.StoreInterface
	generateCode func(int) string
}

func New(store storage.StoreInterface) *URLService {
	return &URLService{
		store:        store,
		generateCode: utils.GenerateCode,
	}
}

func newWithGenerator(store storage.StoreInterface, generator func(int) string) *URLService {
	return &URLService{
		store:        store,
		generateCode: generator,
	}
}

func (s *URLService) Shorten(longURL string) (*models.URL, error) {
	normalizedURL, err := normalizeURL(longURL)
	if err != nil {
		return nil, ErrInvalidURL
	}

	if existing, found := s.store.GetByLongURL(normalizedURL); found {
		return existing, nil
	}

	for attempt := 0; attempt < maxCodeRetries; attempt++ {
		code := s.generateCode(defaultCodeLength)
		if _, found := s.store.Get(code); found {
			continue
		}
		return s.store.Set(code, normalizedURL), nil
	}

	for attempt := 0; attempt < maxCodeRetries; attempt++ {
		code := s.generateCode(defaultCodeLength + 1)
		if _, found := s.store.Get(code); found {
			continue
		}
		return s.store.Set(code, normalizedURL), nil
	}

	return nil, ErrCouldNotCreateCode
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

func normalizeURL(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", fmt.Errorf("%w: empty URL", ErrInvalidURL)
	}

	if !strings.HasPrefix(strings.ToLower(trimmed), "http://") && !strings.HasPrefix(strings.ToLower(trimmed), "https://") {
		trimmed = "https://" + trimmed
	}

	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("%w: unsupported scheme", ErrInvalidURL)
	}

	if parsed.Host == "" {
		return "", fmt.Errorf("%w: missing host", ErrInvalidURL)
	}

	return parsed.String(), nil
}
