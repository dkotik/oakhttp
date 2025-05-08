package secrets

import (
	"context"
	"fmt"
)

type fileRepository struct {
	values map[string]interface{}
}

func (f *fileRepository) SyncMap(name string) func(context.Context) (map[string]interface{}, error) {
	return func(context.Context) (map[string]interface{}, error) {
		value, ok := f.values[name]
		if !ok {
			return nil, ErrSecretNotFound
		}
		m, ok := value.(map[string]interface{})
		if !ok {
			return nil, ErrSecretNotFound
		}
		return m, nil
	}
}

func (f *fileRepository) SyncString(name string) func(context.Context) (string, error) {
	return func(context.Context) (string, error) {
		value, ok := f.values[name]
		if !ok {
			return "", ErrSecretNotFound
		}
		s := fmt.Sprintf("%s", value)
		if s == "" {
			return "", ErrSecretNotFound
		}
		return s, nil
	}

}

func (f *fileRepository) SyncUint(name string) func(context.Context) (uint, error) {
	return func(context.Context) (uint, error) {
		value, ok := f.values[name]
		if !ok {
			return 0, ErrSecretNotFound
		}
		m, ok := value.(uint)
		if !ok {
			return 0, ErrSecretNotFound
		}
		return m, nil
	}
}

func (f *fileRepository) SyncInt(name string) func(context.Context) (int, error) {
	return func(context.Context) (int, error) {
		value, ok := f.values[name]
		if !ok {
			return 0, ErrSecretNotFound
		}
		m, ok := value.(int)
		if !ok {
			return 0, ErrSecretNotFound
		}
		return m, nil
	}
}

// func FromFileJSON(p string) (Repository, error) {
// 	data, err := os.ReadFile(p)
// 	if err != nil {
// 		return nil, err
// 	}
// 	values := make(map[string]interface{})
// 	if err = json.Unmarshal(data, &values); err != nil {
// 		return nil, err
// 	}
// 	return &fileRepository{values}, nil
// }
