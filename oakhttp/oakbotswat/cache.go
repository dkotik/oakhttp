package oakbotswat

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrCacheFull = errors.New("there are too many cache records")

type Cache interface {
	GetToken(context.Context, string) (bool, error)
	SetToken(context.Context, string) error
}

type MapCache struct {
	duration    time.Duration
	recordLimit int
	cleanLimit  int
	mu          sync.Mutex
	tokens      map[string]time.Time
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
		tokens:      make(map[string]time.Time),
	}
}

func (m *MapCache) GetToken(ctx context.Context, token string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	expires, ok := m.tokens[token]
	if !ok {
		return false, nil
	}
	return expires.After(time.Now()), nil
}

func (m *MapCache) SetToken(ctx context.Context, token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	total := len(m.tokens)

	t := time.Now()
	if len(m.tokens) >= m.cleanLimit {
		for existing, expires := range m.tokens {
			if expires.Before(t) {
				delete(m.tokens, existing)
			}
		}
		total = len(m.tokens)
	}
	if total >= m.recordLimit {
		return ErrCacheFull
	}

	m.tokens[token] = t.Add(m.duration)
	return nil
}
