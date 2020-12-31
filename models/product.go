package models

import "database/sql"

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


type ProductService struct {
	Storage *Storage
}

func (ps *ProductService) GetProductBarcodes(productId int64) (*map[string]int, error) {
	var bcVal string
	var bcType int

	bMap := make(map[string]int)

	sqlBc := "SELECT barcode, barcode_type FROM barcodes WHERE product_id = $1"
	rows, err := ps.Storage.Db.Query(sqlBc, productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&bcVal, &bcType)
		if err != nil {
			return nil, err
		}
		bMap[bcVal] = bcType
	}

	if len(bMap) == 0 {
		return nil, sql.ErrNoRows
	}
	return &bMap, nil
}
