package models

const (
	BarcodeTypeEAN13 = iota
	BarcodeTypeEAN8
	BarcodeTypeEAN14
	BarcodeTypeCode128
)

type Manufacturer struct {
	Id   int
	Name string
}

type Product struct {
	Id           int64          `json:"id"`
	Name         string         `json:"name"`
	Barcodes     map[string]int `json:"barcode"`
	Manufacturer Manufacturer   `json:"manufacturer"`
	Size         SpecificSize   `json:"size"`
}

func (p *Product) GetProductId() int64 {
	return p.Id
}

type IProduct interface {
	GetProductId() int64
}
