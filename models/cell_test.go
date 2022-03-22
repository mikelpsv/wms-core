package models

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"log"
	"testing"
)

func NewMock() (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return db, mock
}

func TestCell_GetNumericView(t *testing.T) {
	c := new(Cell)
	c.WhsId = 1
	c.ZoneId = 2
	c.PassageId = 3
	c.RackId = 10
	c.Floor = 4
	strView := c.GetNumericView()

	if strView != "1-02-03-10-04" {
		t.Errorf("invalid view %s", strView)
	}
}

func TestCell_GetNumeric(t *testing.T) {
	c := new(Cell)
	c.WhsId = 1
	c.ZoneId = 2
	c.PassageId = 3
	c.RackId = 10
	c.Floor = 4
	strView := c.GetNumeric()

	if strView != "102031004" {
		t.Errorf("invalid view %s", strView)
	}
}
