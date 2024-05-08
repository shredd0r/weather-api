package repository

import (
	"errors"
	"github.com/google/uuid"
	"weather-api/dto"
	"weather-api/entity"
)

var (
	ErrNotFound = errors.New("volume not found")
)

type Repository[T entity.Entity] interface {
	FindById(Id uuid.UUID) (*T, error)
	GetAll() (*T, error)
	Save(entity T) error
	Delete(Id uuid.UUID) error
	DeleteAll(Id ...uuid.UUID) error
}

type WeatherRepository[T any] interface {
	Repository[T]
	FindByCityNameAndForecaster(cityName string, forecaster dto.WeatherForecaster) (*T, error)
	FindAllByCityNameAndForecaster(cityName string, forecaster dto.WeatherForecaster) ([]*T, error)
}

type LocalityRepository interface {
	Repository[entity.LocalityEntity]
	FindByCityName(cityName string) (*entity.LocalityEntity, error)
}
