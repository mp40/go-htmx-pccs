package state

import (
	"database/sql"
)

type Equipment struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Weight float32 `json:"weight_lbs"`
}

type State struct {
	db *sql.DB
}

func NewStateService(db *sql.DB) *State {
	return &State{db: db}
}

func (s *State) GetEquipment() ([]Equipment, error) {
	rows, err := s.db.Query("SELECT * FROM equipment")
	if err != nil {
		return []Equipment{}, err
	}

	var equipment []Equipment

	for rows.Next() {
		e := Equipment{}
		err := rows.Scan(&e.ID, &e.Name, &e.Weight)
		if err != nil {
			return []Equipment{}, err
		}

		equipment = append(equipment, e)
	}
	return equipment, err
}
