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

	_, err := s.Get(Cell{Id: 2, WhsId: 1, ZoneId: 1}, &prod32, 180, nil)
	if err != nil {
		fmt.Println(err)
	}
	_, err = s.Get(Cell{Id: 2, WhsId: 1, ZoneId: 1}, &prod32, 30, nil)
	if err != nil {
		fmt.Println(err)
	}

}

func TestStorage_FindWhsById(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(1, "test 1")

	mock.ExpectQuery("^SELECT (.+) FROM whs").
		WillReturnRows(rows)
	s := new(Storage)
	s.Db = db
	w, err := s.FindWhsById(1)
	if err != nil {
		t.Error(err)
	}
	if w == nil {
		t.Error(errors.New("whs is nil"))
	}

	rows = sqlmock.NewRows([]string{"id", "name"})
	mock.ExpectQuery("^SELECT (.+) FROM whs").
		WillReturnRows(rows)

	w, err = s.FindWhsById(999)

	if err != sql.ErrNoRows {
		t.Error(err, "error must be sql.ErrNoRows")
	}
	if err == nil {
		t.Error(errors.New("no whs - no error"))
	}
	if w != nil {
		t.Error(errors.New("whs is not nil"))
	}

}

func TestStorage_FindZoneById(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "whs_id", "zone_type"})
	rows.AddRow(1, "test 1", 1, 2)

	mock.ExpectQuery("^SELECT (.+) FROM zones").
		WillReturnRows(rows)
	s := new(Storage)
	s.Db = db
	z, err := s.FindZoneById(1)
	if err != nil {
		t.Error(err)
	}
	if z == nil {
		t.Error(errors.New("cell is nil"))
	}

	rows = sqlmock.NewRows([]string{"id", "name", "whs_id", "zone_type"})
	mock.ExpectQuery("^SELECT (.+) FROM zones").
		WillReturnRows(rows)

	z, err = s.FindZoneById(999)

	if err != sql.ErrNoRows {
		t.Error(err, "error must be sql.ErrNoRows")
	}
	if err == nil {
		t.Error(errors.New("no zone - no error"))
	}
	if z != nil {
		t.Error(errors.New("zone is not nil"))
	}

}

func TestStorage_FindCellById(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "whs_id", "zone_id", "passage_id", "rack_id", "floor"})
	rows.AddRow(1, "test 1", 1, 1, 2, 3, 1)

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

func TestStorage_FindProductById(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	rowsBc := sqlmock.NewRows([]string{"barcode", "barcode_type"})
	rowsBc.AddRow("123456789", 1)

	rows := sqlmock.NewRows([]string{"id", "name", "manufacturer_id", "manufacturer_name"})
	rows.AddRow(1, "test 1", 1, "Pfizer")

	mock.ExpectQuery("^SELECT (.+) FROM products").
		WillReturnRows(rows)

	mock.ExpectQuery("^SELECT (.+) FROM barcodes").
		WillReturnRows(rowsBc)

	s := new(Storage)
	s.Db = db
	p, err := s.FindProductById(1)
	if err != nil {
		t.Error(err)
	}
	if p == nil {
		t.Error(errors.New("product is nil"))
	}

	rowsBc = sqlmock.NewRows([]string{"barcode", "barcode_type"})

	rows = sqlmock.NewRows([]string{"id", "name", "manufacturer_id", "manufacturer_name"})
	mock.ExpectQuery("^SELECT (.+) FROM products").
		WillReturnRows(rows)

	mock.ExpectQuery("^SELECT (.+) FROM barcodes").
		WillReturnRows(rowsBc)

	p, err = s.FindProductById(999)

	if err != sql.ErrNoRows {
		t.Error(err, "error must be sql.ErrNoRows")
	}
	if err == nil {
		t.Error(errors.New("no product - no error"))
	}
	if p != nil {
		t.Error(errors.New("product is not nil"))
	}

}

func TestStorage_FindProductsByBarcode(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	bc := "123456789456"
	// не нашли штрихкод
	rowsBc := sqlmock.NewRows([]string{"product_id", "barcode", "barcode_type"})
	mock.ExpectQuery("^SELECT (.+) FROM barcodes").
		WillReturnRows(rowsBc)

	// до этого набора не должно дойти
	rows := sqlmock.NewRows([]string{"id", "name", "manufacturer_id"})
	rows.AddRow(10, "Тест продукт", 1)
	mock.ExpectQuery("^SELECT (.+) FROM products").
		WillReturnRows(rows)

	s := new(Storage)
	s.Db = db
	p, err := s.FindProductsByBarcode(bc)
	if err != sql.ErrNoRows {
		t.Error("err must be sql.ErrNoRows")
	}

	if p != nil {
		t.Error(errors.New("product must be nil"))
	}

	db, mock = NewMock()
	defer db.Close()

	s = new(Storage)
	s.Db = db

	// нашли штрихкод, но не нашли товар. ошибка странная, но...
	rowsBc = sqlmock.NewRows([]string{"product_id", "barcode", "barcode_type"})
	rowsBc.AddRow(10, bc, 1)
	mock.ExpectQuery("^SELECT (.+) FROM barcodes").
		WillReturnRows(rowsBc)

	rows = sqlmock.NewRows([]string{"id", "name", "manufacturer_id"})
	//rows.AddRow(10, "Тест продукт", 1)
	mock.ExpectQuery("^SELECT (.+) FROM products").
		WillReturnRows(rows)

	p, err = s.FindProductsByBarcode(bc)
	if err != sql.ErrNoRows {
		t.Error("err must be sql.ErrNoRows")
	}
	if p != nil {
		t.Error(errors.New("product must be nil"))
	}

	db, mock = NewMock()
	defer db.Close()
	s = new(Storage)
	s.Db = db

	// нашли штрихкод, нашли товар
	rowsBc = sqlmock.NewRows([]string{"product_id", "barcode", "barcode_type"})
	rowsBc.AddRow(10, bc, 1)
	mock.ExpectQuery("^SELECT (.+) FROM barcodes").
		WillReturnRows(rowsBc)

	rows = sqlmock.NewRows([]string{"id", "name", "manufacturer_id", "manufacturer_name"})
	rows.AddRow(10, "Тест продукт", 1, "производитель")
	mock.ExpectQuery("^SELECT (.+) FROM products").
		WillReturnRows(rows)

	// все штрихкоды для товара
	rowsBcs := sqlmock.NewRows([]string{"barcode", "barcode_type"})
	rowsBcs.AddRow(bc, 1)
	rowsBcs.AddRow("45324523454235", 2)
	rowsBcs.AddRow("65745674567456", 3)

	mock.ExpectQuery("^SELECT (.+) FROM barcodes").
		WillReturnRows(rowsBcs)

	p, err = s.FindProductsByBarcode(bc)
	if err != nil {
		t.Error(err)
	}
	if p == nil {
		t.Error(errors.New("product must not be nil"))
	}
	if p.Barcodes[bc] != 1 {
		t.Error("main barcode must be in result found")
	}

}

func TestStorage_GetProductBarcodes(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	rowsBc := sqlmock.NewRows([]string{"barcode", "barcode_type"})

	mock.ExpectQuery("^SELECT (.+) FROM barcodes").
		WillReturnRows(rowsBc)

	s := new(Storage)
	s.Db = db
	mBc, err := s.GetProductBarcodes(10)
	if err != sql.ErrNoRows {
		t.Error(err, "err must be sql.ErrNoRows")
	}
	if mBc != nil {
		t.Error("result must be nil")
	}

	rowsBc.AddRow("12345678902", 1)
	rowsBc.AddRow("123456789032", 2)
	rowsBc.AddRow("463456789032", 2)

	s = new(Storage)
	s.Db = db

	mBc, err = s.GetProductBarcodes(10)
	if err == sql.ErrNoRows {
		t.Error(err, "err must not be sql.ErrNoRows")
	}

	if mBc != nil && len(*mBc) != 3 {
		t.Error("wrong number of rows count")
	}
}
