package models

import "testing"

func TestProduct_GetProductId(t *testing.T) {

	p := new(Product)
	p.Id = 30
	if p.GetProductId() != 30 {
		t.Error("get product_id fail")
	}
}
