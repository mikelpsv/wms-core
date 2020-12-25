package models

type Whs struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (w *Whs) FindById(id int) {

}

func (w *Whs) Create() {

}

func (w *Whs) Delete() {

}
