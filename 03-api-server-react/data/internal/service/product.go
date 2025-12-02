package service

import (
	"fmt"
	"time"

	"github.com/example/api-server/internal/model"
)

var products = []model.Product{
	{
		ID:          1,
		Name:        "Laptop",
		Description: "High-performance laptop",
		Price:       1299.99,
		Stock:       10,
		CreatedAt:   time.Now(),
	},
	{
		ID:          2,
		Name:        "Mouse",
		Description: "Wireless mouse",
		Price:       29.99,
		Stock:       50,
		CreatedAt:   time.Now(),
	},
}

func GetAllProducts() []model.Product {
	return products
}

func GetProductByID(id int) (*model.Product, error) {
	for i := range products {
		if products[i].ID == id {
			return &products[i], nil
		}
	}
	return nil, fmt.Errorf("product not found")
}

func CreateProduct(req model.CreateProductRequest) model.Product {
	newProduct := model.Product{
		ID:          len(products) + 1,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CreatedAt:   time.Now(),
	}
	products = append(products, newProduct)
	return newProduct
}

func UpdateProductStock(id int, stock int) error {
	for i := range products {
		if products[i].ID == id {
			products[i].Stock = stock
			return nil
		}
	}
	return fmt.Errorf("product not found")
}
