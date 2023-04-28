package oakbotswat

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrCacheFull = errors.New("there are too many cache records")

type Cache interface {
	GetToken(ctx context.Context, key string) (string, bool, error)
	SetToken(ctx context.Context, key, value string) error
}

type UserData struct {
	Data    string
	Expires time.Time
}

type MapCache struct {
	duration    time.Duration
	recordLimit int
	cleanLimit  int
	mu          sync.Mutex
	tokens      map[string]*UserData
}

func NewMapCache(d time.Duration, recordLimit int) *MapCache {
	if d < time.Second {
		d = time.Second
	}
	if recordLimit == 0 {
		recordLimit = 1
	}

	return &MapCache{
		duration:    d,
		recordLimit: recordLimit,
		cleanLimit:  recordLimit / 4,
		mu:          sync.Mutex{},
		tokens:      make(map[string]*UserData),
	}
}

func (m *MapCache) GetToken(ctx context.Context, key string) (string, bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	userData, ok := m.tokens[key]
	if !ok {
		return "", false, nil
	}
	if userData.Expires.Before(time.Now()) {
		delete(m.tokens, key)
		return "", false, nil
	}
	return userData.Data, true, nil
}

func (m *MapCache) SetToken(ctx context.Context, key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	total := len(m.tokens)

	t := time.Now()
	if len(m.tokens) >= m.cleanLimit {
		for existing, userData := range m.tokens {
			if userData.Expires.Before(t) {
				delete(m.tokens, existing)
			}
		}
		total = len(m.tokens)
	}
	if total >= m.recordLimit {
		return ErrCacheFull
	}

	m.tokens[key] = &UserData{
		Data:    value,
		Expires: t.Add(m.duration),
	}
	return nil
}
