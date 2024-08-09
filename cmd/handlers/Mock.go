package handlers

import (
	"sync"
)

type MockRepository struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		data: make(map[string]string),
	}
}

// Save сохраняет shortURL и originalURL в хранилище
// и возвращает ошибку, если что-то пошло не так (в данном случае всегда возвращаем nil)
func (m *MockRepository) Save(shortURL, originalURL string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[shortURL] = originalURL
	return nil // Возвращаем nil, поскольку ошибок нет
}

// Find ищет originalURL по shortURL и возвращает его
// Также возвращает bool, указывающий на наличие URL в хранилище
func (m *MockRepository) Find(shortURL string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	originalURL, exists := m.data[shortURL]
	return originalURL, exists
}
