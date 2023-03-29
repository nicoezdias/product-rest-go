package product

import (
	"errors"
	"fmt"

	"clase18/internal/domain"
)

type Service interface {
	GetAll() ([]domain.Product, error)
	GetByID(id int) (domain.Product, error)
	SearchPriceGt(price float64) ([]domain.Product, error)
	ConsumerPrice(listIdsInt []int) ([]domain.Product, float64, error)
	Create(p domain.Product) (domain.Product, error)
	UpdateProductPut(id int, p domain.Product) (domain.Product, error)
	UpdateProductPatch(id int, updatedProduct domain.Product) (domain.Product, error)
	Delete(id int) error
}

type service struct {
	r Repository
}

// NewService crea un nuevo servicio
func NewService(r Repository) Service {
	return &service{r}
}

// GetAll devuelve todos los productos
func (s *service) GetAll() ([]domain.Product, error) {
	l := s.r.GetAll()
	return l, nil
}

// GetByID busca un producto por su id
func (s *service) GetByID(id int) (domain.Product, error) {
	p, err := s.r.GetByID(id)
	if err != nil {
		return domain.Product{}, err
	}
	return p, nil
}

// SearchPriceGt busca productos por precio mayor que el precio dado
func (s *service) SearchPriceGt(price float64) ([]domain.Product, error) {
	l := s.r.SearchPriceGt(price)
	if len(l) == 0 {
		return []domain.Product{}, errors.New("no products found")
	}
	return l, nil
}

// ConsumerPrice devuelve el precio de una lista de productos
func (s *service) ConsumerPrice(listIdsInt []int) ([]domain.Product, float64, error) {
	products, price, err := s.r.ConsumerPrice(listIdsInt)
	if err != nil {
		return products, price, err
	}
	return products, price, nil
}

// Create agrega un nuevo producto
func (s *service) Create(p domain.Product) (domain.Product, error) {
	p, err := s.r.Create(p)
	if err != nil {
		return domain.Product{}, err
	}
	return p, nil
}

// Update actualiza un producto
func (s *service) UpdateProductPut(id int, u domain.Product) (domain.Product, error) {
	p, err := s.r.UpdateProductPut(id, u)
	if err != nil {
		return domain.Product{}, err
	}
	return p, nil
}

// UpdateProduct actualiza los atributos que cambiaron
func (s *service) UpdateProductPatch(id int, updatedProduct domain.Product) (domain.Product, error) {
	p, err := s.r.UpdateProductPatch(id, updatedProduct)
	if err != nil {
		fmt.Println(err)
		return domain.Product{}, err
	}
	return p, nil
}

// Delete busca un producto por su id y lo elimina
func (s *service) Delete(id int) error {
	err := s.r.Delete(id)
	if err != nil {
		return err
	}
	return nil
}
