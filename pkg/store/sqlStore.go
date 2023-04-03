package store

import (
	"clase19/internal/domain"
	"database/sql"
	"fmt"
	"time"
)

type sqlStore struct {
	DB *sql.DB
}

// NewSqlStore crea un nuevo store de products
func NewSqlStore(db *sql.DB) Store {
	return &sqlStore{
		DB: db,
	}
}

// GetAll devuelve todos los productos
func (s *sqlStore) GetAll() ([]domain.Product, error) {
	var products []domain.Product

	query := "SELECT * FROM products"
	rows, err := s.DB.Query(query)
	if err != nil {
		return []domain.Product{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var productReturn domain.Product
		err = rows.Scan(&productReturn.Id, &productReturn.Name, &productReturn.Quantity, &productReturn.CodeValue, &productReturn.IsPublished, &productReturn.Expiration, &productReturn.Price)
		if err != nil {
			return []domain.Product{}, err
		}
		products = append(products, productReturn)
	}
	if err = rows.Err(); err != nil {
		return []domain.Product{}, err
	}
	return products, nil
}

// GetOne devuelve un producto por su id
func (s *sqlStore) GetOne(id int) (domain.Product, error) {
	var productReturn domain.Product

	query := "SELECT * FROM products WHERE id = ?;"
	row := s.DB.QueryRow(query, id)
	err := row.Scan(&productReturn.Id, &productReturn.Name, &productReturn.Quantity, &productReturn.CodeValue, &productReturn.IsPublished, &productReturn.Expiration, &productReturn.Price)

	if err != nil {
		return domain.Product{}, err
	}
	return productReturn, nil
}

// AddOne agrega un nuevo producto
func (s *sqlStore) AddOne(product domain.Product) (domain.Product, error) {
	stmt, err := s.DB.Prepare("INSERT INTO products(name, quantity, code_value, is_published, expiration, price) VALUES( ?, ?, ?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return domain.Product{}, err
	}
	defer stmt.Close()
	var result sql.Result
	date, err := time.Parse("2006-01-02", product.Expiration)
	if err != nil {
		fmt.Println("aca 2")
		return domain.Product{}, err
	}
	result, err = stmt.Exec(product.Name, product.Quantity, product.CodeValue, product.IsPublished, date, product.Price)
	if err != nil {
		fmt.Println("aca 3")
		return domain.Product{}, err
	}
	insertedId, _ := result.LastInsertId()
	product.Id = int(insertedId)
	return product, nil
}

// UpdateOne actualiza un producto
func (s *sqlStore) UpdateOne(product domain.Product) error {
	p, err := s.GetOne(product.Id)
	productUpdated, err := s.completeEmptyAttributes(p, product)
	if err != nil {
		return err
	}
	stmt, err := s.DB.Prepare("UPDATE products(Name, Quantity, CodeValue, IsPublished, Expiration, Price) VALUES( ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	date, err := time.Parse("2006-01-02", productUpdated.Expiration)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(productUpdated.Name, productUpdated.Quantity, productUpdated.CodeValue, productUpdated.IsPublished, date, productUpdated.Price)
	if err != nil {
		return err
	}
	return nil
}

// DeleteOne elimina un producto
func (s *sqlStore) DeleteOne(id int) error {
	stmt := "DELETE FROM products WHERE id = ?"
	_, err := s.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

// completeEmptyAttributes compara dos productos y se queda con los campos diferentes
func (s *sqlStore) completeEmptyAttributes(product domain.Product, updatedProduct domain.Product) (domain.Product, error) {
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
