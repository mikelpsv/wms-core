package models

// Склад. Физический склад
// Обязательно содержит минимум 3 зоны Zone{} - приемки, хранения и отгрузки.
// Зоны приемки и отгрузки не может быть более 1, эти зоны являются входом и выходом на складе соответственно
type Whs struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Address        string `json:"address"`
	AcceptanceZone Zone
	ShippingZone   Zone
	StorageZones   []Zone
}

type WhsService struct {
	Storage *Storage
}

/*
	Возвращает склад по идентификатору
*/
func (ws *WhsService) FindWhsById(whsId int64) (*Whs, error) {
	sqlCell := "SELECT id, name FROM whs WHERE id = $1"
	row := ws.Storage.Db.QueryRow(sqlCell, whsId)
	w := new(Whs)
	err := row.Scan(&w.Id, &w.Name)
	if err != nil {
		return nil, err
	}
	return w, nil
}

// Возвращает список складов
func (ws *WhsService) GetWarehouses() ([]Whs, error) {
	sqlProd := "SELECT w.id, w.name FROM whs w"
	rows, err := ws.Storage.Query(sqlProd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	whss := make([]Whs, 0, 10)
	for rows.Next() {
		w := new(Whs)
		err = rows.Scan(&w.Id, &w.Name)
		whss = append(whss, *w)
	}
	return whss, nil
}

/*
	Возвращает список зон склада
*/
func (ws *WhsService) GetZones(whs Whs) ([]Zone, error) {
	return ws.GetZonesByWhsId(whs.Id)
}

/*
	Возвращает список зон склада по его идентификатору
*/
func (ws *WhsService) GetZonesByWhsId(whsId int) ([]Zone, error) {
	sqlZones := "SELECT id, name, zone_type FROM zones WHERE whs_id = $1"

	rows, err := ws.Storage.Db.Query(sqlZones, whsId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]Zone, 0)
	for rows.Next() {
		z := Zone{}
		err := rows.Scan(&z.Id, &z.Name, &z.ZoneType)
		if err != nil {
			return nil, err
		}
		res = append(res, z)
	}
	return res, nil
}

/*

 */
func (ws *WhsService) CreateWarehouse(whs *Whs) {

}
