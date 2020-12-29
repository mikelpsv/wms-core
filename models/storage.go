package models

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	Db *sql.DB
}

func (s *Storage) FindWhsById(whsId int64) (*Whs, error) {
	sqlCell := "SELECT id, name FROM whs WHERE id = $1"
	row := s.Db.QueryRow(sqlCell, whsId)
	w := new(Whs)
	err := row.Scan(&w.Id, &w.Name)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Storage) FindZoneById(zoneId int64) (*Zone, error) {
	sqlCell := "SELECT id, name, whs_id, zone_type FROM zones WHERE id = $1"
	row := s.Db.QueryRow(sqlCell, zoneId)
	z := new(Zone)
	err := row.Scan(&z.Id, &z.Name, &z.WhsId, &z.ZoneType)
	if err != nil {
		return nil, err
	}
	return z, nil
}

func (s *Storage) FindCellById(cellId int64) (*Cell, error) {
	sqlCell := "SELECT id, name, whs_id, zone_id, passage_id, rack_id, floor FROM cells WHERE id = $1"
	row := s.Db.QueryRow(sqlCell, cellId)
	c := new(Cell)
	err := row.Scan(&c.Id, &c.Name, &c.WhsId, &c.ZoneId, &c.PassageId, &c.RackId, &c.Floor)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Storage) FindProductById(productId int64) (*Product, error) {
	sqlCell := "SELECT p.id, p.name, p.manufacturer_id, m.name as manufacturer_name " +
		"FROM products p " +
		"LEFT JOIN manufacturers m ON p.manufacturer_id = m.id " +
		"WHERE id = $1"
	row := s.Db.QueryRow(sqlCell, productId)
	p := new(Product)
	err := row.Scan(&p.Id, &p.Name, &p.Manufacturer.Id, &p.Manufacturer.Name)
	if err != nil {
		return nil, err
	}

	pBarcodes, err := s.GetProductBarcodes(p.Id)
	if err != nil {
		return nil, err
	}
	p.Barcodes = *pBarcodes
	return p, nil
}

func (s *Storage) FindProductsByBarcode(barcodeStr string) (*Product, error) {
	var pId int64
	var bcType int
	var bcVal string

	sqlBc := "SELECT product_id, barcode, barcode_type FROM barcodes WHERE barcode = $1"
	row := s.Db.QueryRow(sqlBc, barcodeStr)
	err := row.Scan(&pId, &bcVal, &bcType)
	if err != nil {
		return nil, err
	}

	p, err := s.FindProductById(pId)

	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Storage) GetProductBarcodes(productId int64) (*map[string]int, error) {
	var bcVal string
	var bcType int

	bMap := make(map[string]int)

	sqlBc := "SELECT barcode, barcode_type FROM barcodes WHERE product_id = $1"
	rows, err := s.Db.Query(sqlBc, productId)
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

func (s *Storage) Init(host, dbname, dbuser, dbpass string) error {
	var err error
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", host, dbname, dbuser, dbpass)
	s.Db, err = sql.Open("postgres", connStr)
	if err != nil {
		//log.Fatal(err)
		return err
	}

	err = s.Db.Ping()
	if err != nil {
		//log.Fatal(err)
		return err
	}
	return nil
}

func (s *Storage) Put(cell Cell, prod IProduct, quantity int, tx *sql.Tx) (int, error) {
	var err error
	sql := fmt.Sprintf("INSERT INTO storage%d (zone_id, cell_id, prod_id, quantity) VALUES ($1, $2, $3, $4)", cell.WhsId)
	if tx != nil {
		_, err = tx.Exec(sql, cell.ZoneId, cell.Id, prod.GetProductId(), quantity)
	} else {
		_, err = s.Db.Exec(sql, cell.ZoneId, cell.Id, prod.GetProductId(), quantity)
	}
	if err != nil {
		return quantity, err
	}
	return quantity, nil
}

func (s *Storage) Get(cell Cell, prod IProduct, quantity int, tx *sql.Tx) (int, error) {
	var err error

	if tx == nil {
		tx, err = s.Db.Begin()
		if err != nil {
			// не смогли начать транзакцию
			return 0, err
		}
	}

	sqlInsert := fmt.Sprintf("INSERT INTO storage%d (zone_id, cell_id, prod_id, quantity) VALUES ($1, $2, $3, $4)", cell.WhsId)
	_, err = tx.Exec(sqlInsert, cell.ZoneId, cell.Id, prod.GetProductId(), -1*quantity)
	if err != nil {
		return 0, err
	}

	sqlQuant := fmt.Sprintf("SELECT SUM(quantity) AS quantity "+
		"FROM storage%d WHERE zone_id = $1 AND cell_id = $2 AND prod_id = $3 "+
		"GROUP BY zone_id, cell_id, prod_id "+
		"HAVING SUM(quantity) < 0", cell.WhsId)
	rows, err := tx.Query(sqlQuant, cell.ZoneId, cell.Id, prod.GetProductId())
	if err != nil {
		// ошибка контроля
		return 0, err
	}
	defer rows.Close()
	// мы должны получить пустой запрос
	if rows.Next() {
		err = tx.Rollback()
		if err != nil {
			// ошибка отката... все очень плохо
			return 0, err
		}
		return 0, err
	}
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return quantity, nil
}

func (s *Storage) Quantity(whsId int, cell Cell, tx *sql.Tx) (*map[int]int, error) {
	var zoneId, cellId, prodId, quantity int
	res := make(map[int]int)

	sqlQuantity := fmt.Sprintf("SELECT zone_id, cell_id, prod_id, SUM(quantity) AS quantity "+
		"FROM storage%d WHERE zone_id = $1 AND cell_id = $2 "+
		"GROUP BY zone_id, cell_id, prod_id "+
		"HAVING SUM(quantity) <> 0 %s", whsId, "")

	var err error
	var rows *sql.Rows

	if tx != nil {
		rows, err = tx.Query(sqlQuantity, cell.ZoneId, cell.Id)
	} else {
		rows, err = s.Db.Query(sqlQuantity, cell.ZoneId, cell.Id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&zoneId, &cellId, &prodId, &quantity)
		if err != nil {
			return nil, err
		}
		res[prodId] = quantity
	}
	return &res, nil
}

func (s *Storage) Move(cellSrc, cellDst Cell, prod IProduct, quantity int) error {
	// TODO: cellSrc.WhsId <> cellDst.WhsId - веременной разрыв или виртуальное перемещение

	_, err := s.Get(cellSrc, prod, quantity, nil)
	if err != nil {
		return err
	}
	_, err = s.Put(cellDst, prod, quantity, nil)
	if err == nil {
		return err
	}
	return nil
}
