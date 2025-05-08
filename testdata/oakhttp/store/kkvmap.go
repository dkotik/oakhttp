package store

import (
	"context"
	"errors"
	"sync"
	"time"
)

type mapKeyKeyValue struct {
	duration    time.Duration
	recordLimit int
	total       int
	mu          sync.Mutex
	tokens      map[string]map[string]*expiringValue
}

func NewMapKeyKeyValue(withOptions ...Option) (KeyKeyValue, error) {
	options, err := newOptions(withOptions)
	if err != nil {
		return nil, err
	}
	m := &mapKeyKeyValue{
		duration:    options.retainValuesFor,
		recordLimit: options.valueLimit,

		mu:     sync.Mutex{},
		tokens: make(map[string]map[string]*expiringValue),
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

func (m *mapKeyKeyValue) get(ctx context.Context, key1, key2 string) (*expiringValue, error) {
	mapData, ok := m.tokens[key1]
	if !ok {
		return nil, ErrValueNotFound
	}
	userData, ok := mapData[key2]
	if !ok {
		return nil, ErrValueNotFound
	}
	if userData.Expires.Before(time.Now()) {
		m.total--
		delete(mapData, key2)
		if len(mapData) == 0 {
			delete(m.tokens, key1)
		}
		return nil, ErrValueNotFound
	}
	return userData, nil
}

func (m *mapKeyKeyValue) Get(ctx context.Context, key1, key2 []byte) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, err := m.get(ctx, string(key1), string(key2))
	if err != nil {
		return nil, err
	}
	return data.Data, nil
}

func (m *mapKeyKeyValue) set(ctx context.Context, key1, key2 string, value []byte) error {
	mapData, ok := m.tokens[key1]
	if !ok {
		if m.total >= m.recordLimit {
			return ErrFull
		}
		m.tokens[key1] = map[string]*expiringValue{
			key2: &expiringValue{
				Data:    value,
				Expires: time.Now().Add(m.duration),
			},
		}
		m.total++
		return nil
	}

	data, ok := mapData[key2]
	if !ok {
		if m.total >= m.recordLimit {
			return ErrFull
		}
		mapData[key2] = &expiringValue{
			Data:    value,
			Expires: time.Now().Add(m.duration),
		}
		m.total++
		return nil
	}
	data.Data = value
	return nil
}

func (m *mapKeyKeyValue) Set(ctx context.Context, key1, key2, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.set(ctx, string(key1), string(key2), value)
}

func (m *mapKeyKeyValue) Update(ctx context.Context, key1, key2 []byte, update Update) (err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	k1, k2 := string(key1), string(key2)
	data, err := m.get(ctx, k1, k2)
	if err != nil {
		if errors.Is(ErrValueNotFound, err) {
			value, err := update(nil)
			if err != nil {
				return err
			}
			return m.set(ctx, k1, k2, value)
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

func (m *mapKeyKeyValue) Delete(ctx context.Context, key1, key2 []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	k1, k2 := string(key1), string(key2)
	mapData, ok := m.tokens[k1]
	if ok {
		if _, ok = mapData[k2]; ok {
			m.total--
			delete(mapData, k2)
			if len(mapData) == 0 {
				delete(m.tokens, k1)
			}
		}
	}
	return nil
}

func (m *mapKeyKeyValue) RemoveExpired(ctx context.Context, cutoff time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for mapKey, mapData := range m.tokens {
		for existing, userData := range mapData {
			if userData.Expires.Before(cutoff) {
				delete(mapData, existing)
				m.total--
				if len(existing) == 0 {
					delete(m.tokens, mapKey)
				}
			}
		}
	}
}
