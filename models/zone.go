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
	ZoneType int    `json:"zone_type"`
}

type ZoneService struct {
	Storage *Storage
}

func (zs *ZoneService) FindZoneById(zoneId int64) (*Zone, error) {
	sqlCell := "SELECT id, name, whs_id, zone_type FROM zones WHERE id = $1"
	row := zs.Storage.Db.QueryRow(sqlCell, zoneId)
	z := new(Zone)
	err := row.Scan(&z.Id, &z.Name, &z.WhsId, &z.ZoneType)
	if err != nil {
		return nil, err
	}
	return z, nil
}
