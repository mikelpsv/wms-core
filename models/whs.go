package models

type Whs struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type WhsService struct {
	Storage *Storage
}

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
