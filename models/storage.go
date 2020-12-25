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

func (s *Storage) Put(whsId int, cell Cell, prod Product, quantity int, tx *sql.Tx) (int, error) {
	var err error
	sql := fmt.Sprintf("INSERT INTO storage%d (zone_id, cell_id, prod_id, quantity) VALUES ($1, $2, $3, $4)", whsId)
	if tx != nil {
		_, err = tx.Exec(sql, cell.ZoneId, cell.Id, prod.Id, quantity)
	} else {
		_, err = s.Db.Exec(sql, cell.ZoneId, cell.Id, prod.Id, quantity)
	}
	if err != nil {
		return quantity, err
	}
	return quantity, nil
}

func (s *Storage) Get(whsId int, cell Cell, prod Product, quantity int, tx *sql.Tx) (int, error) {
	var err error

	if tx != nil {
		var quant int

		sqlCell := fmt.Sprintf("SELECT * FROM storage%d WHERE zone_id = $1 AND cell_id = $2 AND prod_id = $3 FOR SHARE", whsId)
		_, err = tx.Query(sqlCell, cell.ZoneId, cell.Id, prod.Id)
		if err != nil {
			//
		}
		// удалось наложить блокировку

		sqlQuantity := fmt.Sprintf("SELECT SUM(quantity) AS quantity "+
			"FROM storage%d WHERE zone_id = $1 AND cell_id = $2 AND prod_id = $3 "+
			"GROUP BY zone_id, cell_id, prod_id "+
			"HAVING SUM(quantity) <> 0", whsId)

		rows, err := tx.Query(sqlQuantity, cell.ZoneId, cell.Id, prod.Id)
		if err != nil {

		}
		defer rows.Close()
		rows.Next()
		rows.Scan(&quant)

		if err != nil {
			tx.Rollback()
			return quantity, err
		}
		if quant < quantity {
			tx.Rollback()
		}
	}

	sql := fmt.Sprintf("INSERT INTO storage%d (zone_id, cell_id, prod_id, quantity) VALUES ($1, $2, $3, $4)", whsId)
	if tx != nil {
		_, err = tx.Exec(sql, cell.ZoneId, cell.Id, prod.Id, -1*quantity)
	} else {
		_, err = s.Db.Exec(sql, cell.ZoneId, cell.Id, prod.Id, -1*quantity)
	}
	if err != nil {
		return quantity, err
	}

	return quantity, nil
}

func (s *Storage) Quantity(whsId int, cell Cell, tx *sql.Tx, forUpdate string) (*map[int]int, error) {
	var zoneId, cellId, prodId, quantity int
	res := make(map[int]int)

	sqlQuantity := fmt.Sprintf("SELECT zone_id, cell_id, prod_id, SUM(quantity) AS quantity "+
		"FROM storage%d WHERE zone_id = $1 AND cell_id = $2 "+
		"GROUP BY zone_id, cell_id, prod_id "+
		"HAVING SUM(quantity) <> 0 %s", whsId, "")
	rows, err := s.Db.Query(sqlQuantity, cell.ZoneId, cell.Id)
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

func (s *Storage) Move(whsId int, cellSrc, cellDst Cell, prod Product, quantity int) error {

	Tx, err := s.Db.Begin()
	if err != nil {
		return err
	}
	_, err = s.Get(whsId, cellSrc, prod, quantity, Tx)
	if err != nil {
		Tx.Rollback()
		return err
	}
	_, err = s.Put(whsId, cellDst, prod, quantity, Tx)
	if err == nil {
		Tx.Rollback()
		return err
	}

	err = Tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
