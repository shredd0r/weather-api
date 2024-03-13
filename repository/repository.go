package repository

type Repository[T any] interface {
	Get(cityName string) (*T, error)
	Save(cityName string, entity T) error
	Delete(cityName string) error
}
