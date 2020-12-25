package models

import "fmt"

// Размеры ячейки в см/см3
type SpecificSize struct {
	Length       int `json:"length"`
	Width        int `json:"width"`
	Height       int `json:"height"`
	Volume       int `json:"volume"`        // lentgth * width * height
	UsefulVolume int `json:"useful_volume"` // lentgth * width * height * K(0.8)
}

type Cell struct {
	Id        int64        `json:"id"`
	Name      string       `json:"name"`
	WhsId     int          `json:"whs_id"`     // Id склада (может быть именован)
	ZoneId    int          `json:"zone_id"`    // Id зоны назначения (может быть именован)
	PassageId int          `json:"passage_id"` // Id проезда (может быть именован)
	RackId    int          `json:"rack_id"`    // Id стеллажа (может быть именован)
	Floor     int          `json:"floor"`
	Size      SpecificSize `json:"size"`
}

// Представление в виде набора чисел
func (c *Cell) GetNumeric() string {
	return fmt.Sprintf("%01d%02d%02d%02d%02d", c.WhsId, c.ZoneId, c.PassageId, c.RackId, c.Floor)
}

// Человеко-понятное представление
func (c *Cell) GetNumericView() string {
	return fmt.Sprintf("%01d-%02d-%02d-%02d-%02d", c.WhsId, c.ZoneId, c.PassageId, c.RackId, c.Floor)
}

func (c *Cell) AddProduct(product Product, quantity int) {

}

func (c *Cell) RemoveProduct(product Product, quantity int) {

}
