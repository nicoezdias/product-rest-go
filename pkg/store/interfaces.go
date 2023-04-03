package store

import "clase19/internal/domain"

type Store interface {
	GetAll() ([]domain.Product, error)
	GetOne(id int) (domain.Product, error)
	AddOne(product domain.Product) (domain.Product, error)
	UpdateOne(product domain.Product) error
	DeleteOne(id int) error
}
