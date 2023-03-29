package store

import (
	"encoding/json"
	"errors"
	"os"

	"clase18/internal/domain"
)

type Store interface {
	GetAll() ([]domain.Product, error)
	GetOne(id int) (domain.Product, error)
	AddOne(product domain.Product) error
	UpdatePut(product domain.Product) error
	UpdatePatch(product domain.Product) error
	DeleteOne(id int) error
	saveProducts(products []domain.Product) error
	loadProducts() ([]domain.Product, error)
}

type jsonStore struct {
	pathToFile string
}

// loadProducts carga los productos desde un archivo json
func (s *jsonStore) loadProducts() ([]domain.Product, error) {
	var products []domain.Product
	file, err := os.ReadFile(s.pathToFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(file), &products)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// saveProducts guarda los productos en un archivo json
func (s *jsonStore) saveProducts(products []domain.Product) error {
	bytes, err := json.Marshal(products)
	if err != nil {
		return err
	}
	return os.WriteFile(s.pathToFile, bytes, 0644)
}

// NewJsonStore crea un nuevo store de products
func NewStore(path string) Store {
	return &jsonStore{
		pathToFile: path,
	}
}

// GetAll devuelve todos los productos
func (s *jsonStore) GetAll() ([]domain.Product, error) {
	products, err := s.loadProducts()
	if err != nil {
		return nil, err
	}
	return products, nil
}

// GetOne devuelve un producto por su id
func (s *jsonStore) GetOne(id int) (domain.Product, error) {
	products, err := s.loadProducts()
	if err != nil {
		return domain.Product{}, err
	}
	for _, product := range products {
		if product.Id == id {
			return product, nil
		}
	}
	return domain.Product{}, errors.New("product not found")
}

// AddOne agrega un nuevo producto
func (s *jsonStore) AddOne(product domain.Product) error {
	products, err := s.loadProducts()
	if err != nil {
		return err
	}
	product.Id = len(products) + 1
	products = append(products, product)
	return s.saveProducts(products)
}

// UpdateOne actualiza un producto
func (s *jsonStore) UpdatePut(product domain.Product) error {
	products, err := s.loadProducts()
	if err != nil {
		return err
	}
	for i, p := range products {
		if p.Id == product.Id {
			products[i] = product
			return s.saveProducts(products)
		}
	}
	return errors.New("product not found")
}

// UpdateOne actualiza un producto
func (s *jsonStore) UpdatePatch(product domain.Product) error {
	products, err := s.loadProducts()
	if err != nil {
		return err
	}
	for i, p := range products {
		if p.Id == product.Id {
			productUpdated, err := s.completeEmptyAttributes(p, product)
			if err != nil {
				return err
			}
			products[i] = productUpdated
			return s.saveProducts(products)
		}
	}
	return errors.New("product not found")
}

// DeleteOne elimina un producto
func (s *jsonStore) DeleteOne(id int) error {
	products, err := s.loadProducts()
	if err != nil {
		return err
	}
	for i, p := range products {
		if p.Id == id {
			products = append(products[:i], products[i+1:]...)
			return s.saveProducts(products)
		}
	}
	return errors.New("product not found")
}

// completeEmptyAttributes compara dos productos y se queda con los campos diferentes
func (s *jsonStore) completeEmptyAttributes(product domain.Product, updatedProduct domain.Product) (domain.Product, error) {
	p := product
	if updatedProduct.Name != "" {
		p.Name = updatedProduct.Name
	}
	if updatedProduct.Quantity != 0 {
		p.Quantity = updatedProduct.Quantity
	}
	if updatedProduct.CodeValue != "" {
		p.CodeValue = updatedProduct.CodeValue
	}
	if updatedProduct.IsPublished {
		p.IsPublished = updatedProduct.IsPublished
	}
	if updatedProduct.Expiration != "" {
		p.Expiration = updatedProduct.Expiration
	}
	if updatedProduct.Price != 0.0 {
		p.Price = updatedProduct.Price
	}
	return p, nil
}
