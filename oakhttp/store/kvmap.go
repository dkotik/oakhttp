package store

import (
	"context"
	"errors"
	"sync"
	"time"
)

type mapKeyValue struct {
	duration    time.Duration
	recordLimit int
	mu          sync.Mutex
	tokens      map[string]*expiringValue
}

func NewMapKeyValue(withOptions ...Option) (KeyValue, error) {
	options, err := newOptions(withOptions)
	if err != nil {
		return nil, err
	}
	m := &mapKeyValue{
		duration:    options.retainValuesFor,
		recordLimit: options.valueLimit,

		mu:     sync.Mutex{},
		tokens: make(map[string]*expiringValue),
	}
	go func(ctx context.Context, frequency time.Duration) {
		t := time.NewTicker(frequency)
		var cutoff time.Time
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				return
			case cutoff = <-t.C:
				m.RemoveExpired(ctx, cutoff)
			}
		}
	}(options.removalContext, options.removeExpiredValuesEvery)
	return m, nil
}

func (m *mapKeyValue) get(ctx context.Context, key string) (*expiringValue, error) {
	userData, ok := m.tokens[key]
	if !ok {
		return nil, ErrValueNotFound
	}
	if userData.Expires.Before(time.Now()) {
		delete(m.tokens, key)
		return nil, ErrValueNotFound
	}
	return userData, nil
}

func (m *mapKeyValue) Get(ctx context.Context, key []byte) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, err := m.get(ctx, string(key))
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

func (m *mapKeyValue) set(ctx context.Context, key string, value []byte) error {
	data, err := m.get(ctx, key)
	if err != nil {
		if errors.Is(ErrValueNotFound, err) {
			if len(m.tokens) >= m.recordLimit {
				return ErrFull
			}
			m.tokens[key] = &expiringValue{
				Data:    value,
				Expires: time.Now().Add(m.duration),
			}
			return nil
		}
		return err
	}
	data.Data = value
	return nil
}

func (m *mapKeyValue) Set(ctx context.Context, key, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.set(ctx, string(key), value)
}

func (m *mapKeyValue) Update(ctx context.Context, key []byte, update Update) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	k := string(key)

	data, err := m.get(ctx, k)
	if err != nil {
		if errors.Is(ErrValueNotFound, err) {
			value, err := update(nil)
			if err != nil {
				return err
			}
			m.tokens[k] = &expiringValue{
				Data:    value,
				Expires: time.Now().Add(m.duration),
			}
		}
		return err
	}
	value, err := update(data.Data)
	if err != nil {
		return err
	}
	data.Data = value
	return nil
}

func (m *mapKeyValue) Delete(ctx context.Context, key []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.tokens, string(key))
	return nil
}

func (m *mapKeyValue) RemoveExpired(ctx context.Context, cutoff time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for existing, userData := range m.tokens {
		if userData.Expires.Before(cutoff) {
			delete(m.tokens, existing)
		}
	}
}
