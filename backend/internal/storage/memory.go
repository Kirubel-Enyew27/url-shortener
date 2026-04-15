package storage

import (
	"sync"
	"url-shortener/internal/models"
)

type StoreInterface interface {
	Set(code, longURL string) *models.URL
	Get(code string) (*models.URL, bool)
	GetByLongURL(longURL string) (*models.URL, bool)
	GetAll() []*models.URL
	IncrementClicks(code string)
}

type MemoryStore struct {
	mu    sync.RWMutex
	urls  map[string]*models.URL
	index map[string]string
}

func New() *MemoryStore {
	return &MemoryStore{
		urls:  make(map[string]*models.URL),
		index: make(map[string]string),
	}
}

func (s *MemoryStore) Set(code, longURL string) *models.URL {
	s.mu.Lock()
	defer s.mu.Unlock()

	url := &models.URL{
		LongURL:   longURL,
		ShortCode: code,
		Clicks:    0,
	}

	s.urls[code] = url
	s.index[longURL] = code
	return url
}

func (s *MemoryStore) Get(code string) (*models.URL, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.urls[code]
	return url, exists
}

func (s *MemoryStore) GetByLongURL(longURL string) (*models.URL, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	code, exists := s.index[longURL]
	if !exists {
		return nil, false
	}

	url, exists := s.urls[code]
	return url, exists
}

func (s *MemoryStore) GetAll() []*models.URL {
	s.mu.RLock()
	defer s.mu.RUnlock()

	urls := make([]*models.URL, 0, len(s.urls))
	for _, url := range urls {
		urls = append(urls, url)
	}

	return urls
}

func (s *MemoryStore) IncrementClicks(code string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if url, exists := s.urls[code]; exists {
		url.Clicks++
	}
}
