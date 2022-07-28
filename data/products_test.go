package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name:  "Anh",
		Price: 1.99,
		SKU:   "aaa-xxx-vvv",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
