package models

type Product struct {
	Id      int64        `json:"id"`
	Name    string       `json:"name"`
	Barcode string       `json:"barcode"`
	Size    SpecificSize `json:"size"`
}

func (p *Product) GetProductId() int64 {
	return p.Id
}

type IProduct interface {
	GetProductId() int64
}
