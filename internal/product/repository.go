package product

import (
	"errors"
	"fmt"

	"clase18/internal/domain"
	"clase18/pkg/store"
)

type Repository interface {
	GetAll() []domain.Product
	GetByID(id int) (domain.Product, error)
	SearchPriceGt(price float64) []domain.Product
	ConsumerPrice(listIdsInt []int) ([]domain.Product, float64, error)
	Create(p domain.Product) (domain.Product, error)
	UpdateProductPut(id int, p domain.Product) (domain.Product, error)
	UpdateProductPatch(id int, updatedProduct domain.Product) (domain.Product, error)
	Delete(id int) error
}

type repository struct {
	storage store.Store
}

// NewRepository crea un nuevo repositorio
func NewRepository(storage store.Store) Repository {
	return &repository{storage}
}

// GetAll devuelve todos los productos
func (r *repository) GetAll() []domain.Product {
	products, err := r.storage.GetAll()
	if err != nil {
		return []domain.Product{}
	}
	return products
}

// GetByID busca un producto por su id
func (r *repository) GetByID(id int) (domain.Product, error) {
	product, err := r.storage.GetOne(id)
	if err != nil {
		return domain.Product{}, errors.New(fmt.Sprintf("product %d not found", id))
	}
	return product, nil
}

// SearchPriceGt busca productos por precio mayor o igual que el precio dado
func (r *repository) SearchPriceGt(price float64) []domain.Product {
	var products []domain.Product
	list, err := r.storage.GetAll()
	if err != nil {
		return products
	}
	for _, product := range list {
		if product.Price > price {
			products = append(products, product)
		}
	}
	return products
}

// ConsumerPrice devuelve el precio de una lista de productos
func (r *repository) ConsumerPrice(listIdsInt []int) ([]domain.Product, float64, error) {
	cant := 0
	price := 0.0
	var products []domain.Product
	for _, id := range listIdsInt {
		flag := true
		for k, p := range products {
			if id == p.Id {
				if p.Quantity <= 0 {
					return []domain.Product{}, 0, errors.New(fmt.Sprintf("product(%d) stock not available", p.Id))
				}
				products[k].Quantity -= 1
				price += p.Price
				cant++
				flag = false
				break
			}
		}
		if flag {
			product, err := r.GetByID(id)
			if err != nil {
				return []domain.Product{}, 0, err
			}
			err = validProduct(product)
			if err != nil {
				return []domain.Product{}, 0, err
			}
			product.Quantity -= 1
			products = append(products, product)
			price += product.Price
			cant++
		}
	}
	if cant <= 10 {
		price *= 1.21
	} else if cant > 10 && cant < 20 {
		price *= 1.17
	} else {
		price *= 1.15
	}
	return products, price, nil
}

// Create agrega un nuevo producto
func (r *repository) Create(p domain.Product) (domain.Product, error) {
	if !r.validateCodeValue(p.CodeValue) {
		return domain.Product{}, errors.New("code value already exists")
	}
	err := r.storage.AddOne(p)
	if err != nil {
		return domain.Product{}, errors.New("error creating product")
	}
	return p, nil
}

// validateCodeValue valida que el codigo no exista en la lista de productos
func (r *repository) validateCodeValue(codeValue string) bool {
	list, err := r.storage.GetAll()
	if err != nil {
		return false
	}
	for _, product := range list {
		if product.CodeValue == codeValue {
			return false
		}
	}
	return true
}

// Update actualiza un producto
func (r *repository) UpdateProductPut(id int, p domain.Product) (domain.Product, error) {
	if !r.validateCodeValue(p.CodeValue) {
		return domain.Product{}, errors.New("code value already exists")
	}
	err := r.storage.UpdatePut(p)
	if err != nil {
		return domain.Product{}, errors.New("error updating product")
	}
	return p, nil
}

// UpdateProduct actualiza los atributos que cambiaron
func (r *repository) UpdateProductPatch(id int, updatedProduct domain.Product) (domain.Product, error) {
	if !r.validateCodeValue(updatedProduct.CodeValue) {
		return domain.Product{}, errors.New("code value already exists")
	}
	err := r.storage.UpdatePut(updatedProduct)
	if err != nil {
		return domain.Product{}, errors.New("error updating product")
	}
	return updatedProduct, nil
}

// Delete busca un producto por su id y lo elimina
func (r *repository) Delete(id int) error {
	err := r.storage.DeleteOne(id)
	if err != nil {
		return err
	}
	return nil
}

// validProduct comprueba si un producto cumple con los requisitos para ser comprado
func validProduct(product domain.Product) error {
	if product.Quantity <= 0 {
		return errors.New(fmt.Sprintf("product(%d) stock not available", product.Id))
	}
	if !product.IsPublished {
		return errors.New(fmt.Sprintf("product(%d) is not published", product.Id))
	}
	return nil
}

// loadProducts carga los productos desde un archivo json
// func loadProducts(path string, list *[]domain.Product) {
// 	file, err := os.ReadFile(path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = json.Unmarshal([]byte(file), &list)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// completeEmptyAttributes compara dos productos y se queda con los campos diferentes
// func (r *repository) completeEmptyAttributes(product domain.Product, updatedProduct domain.Product) (domain.Product, error) {
// 	p := product
// 	if updatedProduct.Name != "" {
// 		p.Name = updatedProduct.Name
// 	}
// 	if updatedProduct.Quantity != 0 {
// 		p.Quantity = updatedProduct.Quantity
// 	}
// 	if updatedProduct.CodeValue != "" {
// 		if !r.validateCodeValue(updatedProduct.CodeValue) {
// 			return domain.Product{}, errors.New("code value already exists")
// 		}
// 		p.CodeValue = updatedProduct.CodeValue
// 	}
// 	if updatedProduct.IsPublished {
// 		p.IsPublished = updatedProduct.IsPublished
// 	}
// 	if updatedProduct.Expiration != "" {
// 		p.Expiration = updatedProduct.Expiration
// 	}
// 	if updatedProduct.Price != 0.0 {
// 		p.Price = updatedProduct.Price
// 	}
// 	return p, nil
// }