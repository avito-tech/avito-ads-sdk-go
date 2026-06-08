package avito

import (
	"sync"
	"time"
)

// accessToken — токен доступа с абсолютным временем истечения.
type accessToken struct {
	Value     string
	ExpiresAt time.Time
	TokenType string
}

// authorizationHeader возвращает готовое значение заголовка Authorization.
func (t accessToken) authorizationHeader() string {
	typ := t.TokenType
	if typ == "" {
		typ = "Bearer"
	}
	return typ + " " + t.Value
}

// isExpired сообщает, истёк ли токен с учётом запаса leeway.
func (t accessToken) isExpired(leeway time.Duration, now time.Time) bool {
	return !now.Before(t.ExpiresAt.Add(-leeway))
}

// TokenStorage — интерфейс хранилища токена доступа.
//
// Реализуйте его (например, поверх Redis), чтобы переиспользовать токен между
// перезапусками процесса, и передайте через WithTokenStorage.
type TokenStorage interface {
	// Get возвращает токен по ключу и признак его наличия.
	Get(key string) (token string, expiresAt time.Time, ok bool)
	// Set сохраняет токен с временем истечения.
	Set(key, token string, expiresAt time.Time)
	// Delete удаляет токен по ключу.
	Delete(key string)
}

// MemoryTokenStorage — потокобезопасное хранилище токена в памяти (по умолчанию).
type MemoryTokenStorage struct {
	mu      sync.RWMutex
	token   string
	expires time.Time
	present bool
	key     string
}

// NewMemoryTokenStorage создаёт пустое хранилище в памяти.
func NewMemoryTokenStorage() *MemoryTokenStorage {
	return &MemoryTokenStorage{}
}

// Get реализует TokenStorage.
func (s *MemoryTokenStorage) Get(key string) (string, time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.present || s.key != key {
		return "", time.Time{}, false
	}
	return s.token, s.expires, true
}

// Set реализует TokenStorage.
func (s *MemoryTokenStorage) Set(key, token string, expiresAt time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.key = key
	s.token = token
	s.expires = expiresAt
	s.present = true
}

// Delete реализует TokenStorage.
func (s *MemoryTokenStorage) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.key == key {
		s.present = false
		s.token = ""
	}
}
