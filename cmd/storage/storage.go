package storage

import "sync"

// InMemoryStorage - структура для хранения URL в памяти.
type InMemoryStorage struct {
	data map[string]string
	mu   sync.RWMutex
}

// NewInMemoryStorage - конструктор для InMemoryStorage.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]string),
	}
}

// Save - сохраняет короткий и оригинальный URL в память.
func (s *InMemoryStorage) Save(shortURL, originalURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[shortURL] = originalURL
	return nil
}

// Find - ищет оригинальный URL по короткому.
func (s *InMemoryStorage) Find(shortURL string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	originalURL, exists := s.data[shortURL]
	return originalURL, exists
}
