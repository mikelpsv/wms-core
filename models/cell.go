package models

import "fmt"

// Размеры ячейки в см/см3
type SpecificSize struct {
	length       int
	width        int
	height       int
	volume       float32 // lentgth * width * height
	usefulVolume float32 // lentgth * width * height * K(0.8)
}

type Cell struct {
	Id           int64        `json:"id"`
	Name         string       `json:"name"`
	WhsId        int          `json:"whs_id"`     // Id склада (может быть именован)
	ZoneId       int          `json:"zone_id"`    // Id зоны назначения (может быть именован)
	PassageId    int          `json:"passage_id"` // Id проезда (может быть именован)
	RackId       int          `json:"rack_id"`    // Id стеллажа (может быть именован)
	Floor        int          `json:"floor"`
	IsSizeFree   bool         `json:"is_size_free"`
	IsWeightFree bool         `json:"is_weight_free"`
	Size         SpecificSize `json:"size"`
}

type CellService struct {
	Storage *Storage
}

func (sz *SpecificSize) SetSize(length, width, height int, kUV float32) {
	sz.volume = float32(length * width * height)
	sz.usefulVolume = sz.volume * kUV
}

func (sz *SpecificSize) GetSize() (int, int, int, float32, float32) {
	return sz.length, sz.width, sz.height, sz.volume, sz.usefulVolume
}

// Представление в виде набора чисел
func (c *Cell) GetNumeric() string {
	return fmt.Sprintf("%01d%02d%02d%02d%02d", c.WhsId, c.ZoneId, c.PassageId, c.RackId, c.Floor)
}

// Человеко-понятное представление
func (c *Cell) GetNumericView() string {
	return fmt.Sprintf("%01d-%02d-%02d-%02d-%02d", c.WhsId, c.ZoneId, c.PassageId, c.RackId, c.Floor)
}

func (cs *CellService) AddProduct(cell *Cell, product *Product, quantity int) {

	fmt.Println(product.Id)
}

