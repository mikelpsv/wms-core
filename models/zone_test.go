package models

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestStorage_FindZoneById(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name", "whs_id", "zone_type"})
	rows.AddRow(1, "test 1", 1, 2)

	mock.ExpectQuery("^SELECT (.+) FROM zones").
		WillReturnRows(rows)
	s := new(Storage)
	zs := s.GetZoneService()

	s.Db = db
	z, err := zs.FindZoneById(1)
	if err != nil {
		t.Error(err)
	}
	if z == nil {
		t.Error(errors.New("cell is nil"))
	}

	rows = sqlmock.NewRows([]string{"id", "name", "whs_id", "zone_type"})
	mock.ExpectQuery("^SELECT (.+) FROM zones").
		WillReturnRows(rows)

	z, err = zs.FindZoneById(999)

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
