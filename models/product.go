package models

import "database/sql"

// Продукт. Единица хранения.
type Product struct {
	Id           int64          `json:"id"`
	Name         string         `json:"name"`
	Barcodes     map[string]int `json:"barcode"`
	Manufacturer Manufacturer   `json:"manufacturer"`
	Size         SpecificSize   `json:"size"`
}

// Производитель
type Manufacturer struct {
	Id   int64 `json:"id"`
	Name string `json:"name"`
}

type ProductService struct {
	Storage *Storage
}

// Возвращает список штрихкодов продукта
func (ps *ProductService) GetProductBarcodes(productId int64) (map[string]int, error) {
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
	return bMap, nil
}

// Возвращает продукт по внутреннему идентификатору
func (ps *ProductService) FindProductById(productId int64) (*Product, error) {

	sqlCell := "SELECT p.id, p.name, p.manufacturer_id, m.name as manufacturer_name " +
		"FROM products p " +
		"LEFT JOIN manufacturers m ON p.manufacturer_id = m.id " +
		"WHERE p.id = $1"
	row := ps.Storage.Db.QueryRow(sqlCell, productId)
	p := new(Product)
	err := row.Scan(&p.Id, &p.Name, &p.Manufacturer.Id, &p.Manufacturer.Name)
	if err != nil {
		return nil, err
	}

	pBarcodes, err := ps.GetProductBarcodes(p.Id)
	if err != nil {
		return nil, err
	}
	p.Barcodes = pBarcodes
	return p, nil
}

// Возвращает продукт по штрихкоду
func (ps *ProductService) FindProductsByBarcode(barcodeStr string) (*Product, error) {
	var pId int64
	var bcType int
	var bcVal string

	sqlBc := "SELECT product_id, barcode, barcode_type FROM barcodes WHERE barcode = $1"
	row := ps.Storage.Db.QueryRow(sqlBc, barcodeStr)
	err := row.Scan(&pId, &bcVal, &bcType)
	if err != nil {
		return nil, err
	}

	p, err := ps.FindProductById(pId)

	if err != nil {
		return nil, err
	}

	return p, nil
}

// Возвращает список продуктов
func (ps *ProductService) GetProducts() ([]Product, error) {
	sqlProd := "SELECT p.id, p.name, p.manufacturer_id, m.name FROM products p LEFT JOIN manufacturers m ON p.manufacturer_id = m.id"
	rows, err := ps.Storage.Query(sqlProd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prods := make([]Product, 0, 10)
	for rows.Next() {
		p := new(Product)
		err = rows.Scan(&p.Id, &p.Name, &p.Manufacturer.Id, &p.Manufacturer.Name)
		prods = append(prods, *p)
	}
	return prods, nil
}

// Возвращает список производителей
func (ps *ProductService) GetManufacturers() ([]Manufacturer, error) {
	sqlMnf := "SELECT m.id, m.name FROM manufacturers m"
	rows, err := ps.Storage.Query(sqlMnf)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mnfs := make([]Manufacturer, 0, 10)
	for rows.Next() {
		m := new(Manufacturer)
		err = rows.Scan(&m.Id, &m.Name)
		mnfs = append(mnfs, *m)
	}
	return mnfs, nil
}

