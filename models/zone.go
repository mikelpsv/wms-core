package models

const (
	ZoneTypeStorage = iota
	ZoneTypeIncoming
	ZoneTypeOutGoing
	ZoneTypeCustom = 99
)

type Zone struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	WhsId    int    `json:"whs_id"`
	ZoneType int
}

func (z *Zone) GetWhs() *Whs {
	whs := new(Whs)
	whs.FindById(z.WhsId)
	return whs
}
