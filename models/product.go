package models

type Product struct {
	Id      int64        `json:"id"`
	Name    string       `json:"name"`
	Barcode string       `json:"barcode"`
	Size    SpecificSize `json:"size"`
}
