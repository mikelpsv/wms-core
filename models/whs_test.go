package models

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
)

func TestStorage_FindWhsById(t *testing.T) {
	db, mock := NewMock()
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "name"})
	rows.AddRow(1, "test 1")

	mock.ExpectQuery("^SELECT (.+) FROM whs").
		WillReturnRows(rows)
	s := new(Storage)
	s.Db = db
	ws := s.GetWhsService()

	w, err := ws.FindWhsById(1)
	if err != nil {
		t.Error(err)
	}
	if w == nil {
		t.Error(errors.New("whs is nil"))
	}

	rows = sqlmock.NewRows([]string{"id", "name"})
	mock.ExpectQuery("^SELECT (.+) FROM whs").
		WillReturnRows(rows)

	w, err = ws.FindWhsById(999)

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
