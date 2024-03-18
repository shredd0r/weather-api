package repository

import "errors"

var (
	ErrNotFound = errors.New("volume not found")
)

type Repository[T any] interface {
	Get(key string) (*T, error)
	Save(key string, entity T) error
	Delete(key string) error
}
