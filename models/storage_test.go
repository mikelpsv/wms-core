package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestStorage_Init(t *testing.T) {
	s := new(Storage)
	s.Init("localhost", "wmsdb", "devuser", "devuser")

	prod32 := Product{
		Id:       32,
		Name:     "tedst",
		Barcodes: make(map[string]int),
		Size:     SpecificSize{},
	}

	c := Cell{Id: 2, WhsId: 1, ZoneId: 1}
	_, err := s.Get(&c, &prod32, 180, nil)
	if err != nil {
		fmt.Println(err)
	}
	_, err = s.Get(&c, &prod32, 30, nil)
	if err != nil {
		fmt.Println(err)
	}

}

func TestStorage_FindCellById(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "whs_id", "zone_id", "passage_id", "rack_id", "floor",
		"sz_length", "sz_width", "sz_height", "sz_volume", "sz_if_volume", "sz_weight",
		"not_allowed_in", "not_allowed_out", "is_service", "is_size_free", "is_weight_free"})
	rows.AddRow(1, "test 1", 1, 1, 2, 3, 1, 2, 2, 2, 2, 2, 2, false, false, false, false, false)

	mock.ExpectQuery("^SELECT (.+) FROM cells").
		WillReturnRows(rows)
	s := new(Storage)
	s.Db = db
	c, err := s.FindCellById(1)
	if err != nil {
		t.Error(err)
	}
	if c == nil {
		t.Error(errors.New("cell is nil"))
	}

	rows = sqlmock.NewRows([]string{"id", "name", "whs_id", "zone_id", "passage_id", "rack_id", "floor"})
	mock.ExpectQuery("^SELECT (.+) FROM cells").
		WillReturnRows(rows)

	c, err = s.FindCellById(999)

	if err != sql.ErrNoRows {
		t.Error(err, "error must be sql.ErrNoRows")
	}
	if err == nil {
		t.Error(errors.New("no cell - no error"))
	}
	if c != nil {
		t.Error(errors.New("cell is not nil"))
	}

}
