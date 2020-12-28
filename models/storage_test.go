package models

import (
	"fmt"
	"testing"
)

func TestStorage_Init(t *testing.T) {
	s := new(Storage)
	s.Init("localhost", "wmsdb", "devuser", "devuser")

	_, err := s.Get(Cell{Id: 2, WhsId: 1, ZoneId: 1}, Product{Id: 32}, 180, nil)
	if err != nil {
		fmt.Println(err)
	}
	_, err = s.Get(Cell{Id: 2, WhsId: 1, ZoneId: 1}, Product{Id: 32}, 30, nil)
	if err != nil {
		fmt.Println(err)
	}

}

func TestStorage_Init2(t *testing.T) {

}
