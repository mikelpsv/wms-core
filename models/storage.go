package models

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	Db *sql.DB
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

func (s *Storage) GetProductService() *ProductService {
	ps := new(ProductService)
	ps.Storage = s
	return ps
}

func (s *Storage) GetWhsService() *WhsService {
	ws := new(WhsService)
	ws.Storage = s
	return ws
}

func (s *Storage) GetZoneService() *ZoneService {
	zs := new(ZoneService)
	zs.Storage = s
	return zs
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

func (s *Storage) Put(cell Cell, prod IProduct, quantity int, tx *sql.Tx) (int, error) {
	var err error
	sqlIns := fmt.Sprintf("INSERT INTO storage%d (zone_id, cell_id, prod_id, quantity) VALUES ($1, $2, $3, $4)", cell.WhsId)
	if tx != nil {
		_, err = tx.Exec(sqlIns, cell.ZoneId, cell.Id, prod.GetProductId(), quantity)
	} else {
		_, err = s.Db.Exec(sqlIns, cell.ZoneId, cell.Id, prod.GetProductId(), quantity)
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
