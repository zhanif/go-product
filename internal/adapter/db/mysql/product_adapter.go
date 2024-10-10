package mysql

import (
	"database/sql"
	"fmt"
	"go-product/internal/core/domain"
	"go-product/internal/core/port"

	"github.com/google/uuid"
)

type ProductRepositoryImpl struct {
	db *sql.DB
}

var _ port.ProductRepository = &ProductRepositoryImpl{}

func NewProductRepository(db *sql.DB) port.ProductRepository {
	return &ProductRepositoryImpl{db: db}
}

func (r *ProductRepositoryImpl) FindById(id string) (*domain.Product, error) {
	var product domain.Product

	row := r.db.QueryRow("SELECT id, name, price, stock FROM products WHERE id = ?", id)

	if err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock); err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepositoryImpl) FindAll() ([]domain.Product, error) {
	rows, err := r.db.Query("SELECT id, name, price, stock FROM products")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []domain.Product{}
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepositoryImpl) Save(product domain.Product) error {
	id := uuid.New().String()
	_, err := r.db.Exec("INSERT INTO products (id, name, price, stock) VALUES (?, ?, ?, ?)", id, product.Name, product.Price, product.Stock)
	return err
}

func (r *ProductRepositoryImpl) Update(id string, product domain.Product) error {
	_, err := r.db.Exec("UPDATE products SET name = ?, price = ?, stock = ? WHERE id = ?", product.Name, product.Price, product.Stock, id)
	return err
}

func (r *ProductRepositoryImpl) Destroy(id string) error {
	result, err := r.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		return err
	}

	if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
		return fmt.Errorf("product with ID %s not found", id)
	}

	return nil
}
