package data

import (
	"fmt"
	"time"
)

var ErrProductNotFound = fmt.Errorf("product not found")

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}

type Products []*Product

func GetProducts() Products {
	return productList
}

func GetProductByID(id int) (*Product, error) {
	i := findIndexByProductId(id)
	if i == -1 {
		return nil, ErrProductNotFound
	}

	return productList[i], nil
}

func UpdateProduct(p *Product) error {
	i := findIndexByProductId(p.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	productList[i] = p

	return nil
}

func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

func DeleteProduct(id int) error {
	i := findIndexByProductId(id)
	if i == -1 {
		return ErrProductNotFound
	}

	productList = append(productList[:i], productList[i+1:]...)

	return nil
}

func getNextID() int {
	p := productList[len(productList)-1]
	return p.ID + 1
}

func findIndexByProductId(id int) int {
	for i, p := range productList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc123",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Frothy milky coffee",
		Price:       1.99,
		SKU:         "dfgh451",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
