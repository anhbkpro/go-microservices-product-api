package handlers

import (
	"github.com/anhbkpro/go-microservices-product-api/data"
	"net/http"
)

func (p *Products) Create(rw http.ResponseWriter, r *http.Request) {
	// fetch the product from the context
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	p.l.Println("[DEBUG] inserting product: %#v\n", prod)
	data.AddProduct(&prod)
}
