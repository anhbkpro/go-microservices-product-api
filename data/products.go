package data

import (
	"fmt"
)

// ErrProductNotFound is an error raised when a product can not be found in the database
var ErrProductNotFound = fmt.Errorf("product not found")

// Product defines the structure for an API product
// swagger:model
type Product struct {
	// the id for the product
	//
	// required: false
	// min: 1
	ID int `json:"id"`

	// the name for this product
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this product
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the product
	//
	// required: true
	// min: 0.01
	Price float32 `json:"price" validate:"gt=0"`

	// the SKU for the product
	//
	// required: true
	// pattern: [a-z]+-[a-z]+-[a-z]+
	SKU string `json:"sku" validate:"required,sku"`
}

// Products defines a slice of Product
type Products []*Product

// GetProducts returns all products from the database
func GetProducts() Products {
	return productList
}

// GetProductByID returns a single product which matches the id from the
// database.
// If a product is not found this function returns a ProductNotFound error
func GetProductByID(id int) (*Product, error) {
	i := findIndexByProductId(id)
	if i == -1 {
		return nil, ErrProductNotFound
	}

	return productList[i], nil
}

// UpdateProduct replaces a product in the database with the given
// item.
// If a product with the given id does not exist in the database
// this function returns a ProductNotFound error
func UpdateProduct(p *Product) error {
	i := findIndexByProductId(p.ID)
	if i == -1 {
		return ErrProductNotFound
	}

	productList[i] = p

	return nil
}

// AddProduct adds a new product to the database
func AddProduct(p *Product) {
	p.ID = getNextID()
	productList = append(productList, p)
}

// DeleteProduct deletes a product from the database
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

// findIndex finds the index of a product in the database
// returns -1 when no product can be found
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
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Frothy milky coffee",
		Price:       1.99,
		SKU:         "ghk451",
	},
}
