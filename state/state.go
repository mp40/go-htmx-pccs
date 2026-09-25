package state

import (
	"database/sql"
)

type Equipment struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Weight float32 `json:"weight_lbs"`
}

type Uniform struct {
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

func (s *State) GetUniforms() ([]Uniform, error) {
	rows, err := s.db.Query("SELECT * FROM uniform")
	if err != nil {
		return []Uniform{}, err
	}

	var uniforms []Uniform

	for rows.Next() {
		u := Uniform{}
		err := rows.Scan(&u.ID, &u.Name, &u.Weight)
		if err != nil {
			return []Uniform{}, err
		}

		uniforms = append(uniforms, u)
	}
	return uniforms, err
}

func (s *State) GetUniformByID(uniformID int) (*Uniform, error) {
	uniform := Uniform{}
	err := s.db.QueryRow("SELECT id, name, weight_lbs FROM uniform WHERE id = ?", uniformID).
		Scan(&uniform.ID, &uniform.Name, &uniform.Weight)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &uniform, err
}
