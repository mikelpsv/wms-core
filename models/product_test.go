package models

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestProduct_GetProductId(t *testing.T) {

	p := new(Product)
	p.Id = 30
	if p.Id != 30 {
		t.Error("get product_id fail")
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
	ps := s.GetProductService()
	p, err := ps.FindProductById(1)
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

	p, err = ps.FindProductById(999)

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
	ps := s.GetProductService()

	p, err := ps.FindProductsByBarcode(bc)
	if err != sql.ErrNoRows {
		t.Error("err must be sql.ErrNoRows")
	}

	if p != nil {
		t.Error(errors.New("product must be nil"))
	}

	//////////////////////////////////

	db, mock = NewMock()
	defer db.Close()

	s = new(Storage)
	s.Db = db
	ps = s.GetProductService()

	// нашли штрихкод, но не нашли товар. ошибка странная, но...
	rowsBc = sqlmock.NewRows([]string{"product_id", "barcode", "barcode_type"})
	rowsBc.AddRow(10, bc, 1)
	mock.ExpectQuery("^SELECT (.+) FROM barcodes").
		WillReturnRows(rowsBc)

	rows = sqlmock.NewRows([]string{"id", "name", "manufacturer_id"})
	//rows.AddRow(10, "Тест продукт", 1)
	mock.ExpectQuery("^SELECT (.+) FROM products").
		WillReturnRows(rows)

	p, err = ps.FindProductsByBarcode(bc)
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
	ps = s.GetProductService()

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

	p, err = ps.FindProductsByBarcode(bc)
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
	ps := s.GetProductService()
	mBc, err := ps.GetProductBarcodes(10)
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

	mBc, err = ps.GetProductBarcodes(10)
	if err == sql.ErrNoRows {
		t.Error(err, "err must not be sql.ErrNoRows")
	}

	if mBc != nil && len(mBc) != 3 {
		t.Error("wrong number of rows count")
	}
}
