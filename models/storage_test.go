package models

import (
	"fmt"
	"testing"
)

func TestStorage_Init(t *testing.T) {
	s := new(Storage)
	s.Init("localhost", "wmsdb", "devuser", "devuser")

	//_, err := s.Put(1, Cell{Id: 2}, Product{Id: 34}, 40, nil)
	mapProd, err := s.Quantity(1, Cell{Id: 2}, nil, "")

	tx, err := s.Db.Begin()
	_, err = s.Get(1, Cell{Id: 2}, Product{Id: 32}, 201, tx)
	if err != nil {
		tx.Commit()
	}

	mapProd, err = s.Quantity(1, Cell{Id: 2}, nil, "")
	//s.Move(1, Cell{Id: 2}, Cell{Id: 3}, Product{Id: 34}, 18)

	fmt.Println(mapProd)
	fmt.Println(err)
}
