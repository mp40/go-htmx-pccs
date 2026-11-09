package state

import (
	"database/sql"
	"fmt"
	"strings"
)

type Firearm struct {
	ID                 int      `json:"id"`
	Name               string   `json:"name"`
	BaseWeight         float32  `json:"base_weight"`
	Length             float32  `json:"length"`
	LengthStockFolded  *float32 `json:"length_stock_folded"`
	InternalCapacity   *int     `json:"internal_capacity"`
	RoundWeight        *float32 `json:"round_weight"`
	ReloadTime         int      `json:"reload_time"`
	ManualCycleTime    *int     `json:"manual_cycle_time"`
	SemiAutomatic      bool     `json:"semi_automatic"`
	BurstRounds        *int     `json:"burst_rounds"`
	FullAutoRateOfFire *int     `json:"full_auto_rof"`
	KnockDown          int      `json:"knock_down"`
	SAB                int      `json:"sab"`
	AimTimeID          int      `json:"aim_time_id"`
}

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

func (s *State) GetFirearms() ([]Firearm, error) {
	rows, err := s.db.Query("SELECT id, name, base_weight, length, length_stock_folded, internal_capacity, round_weight, reload_time, manual_cycle_time, semi_automatic, burst_rounds, full_auto_rof, knock_down, sab, aim_time_id FROM firearm")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	firearms := []Firearm{}

	for rows.Next() {
		f := Firearm{}
		err := rows.Scan(&f.ID, &f.Name, &f.BaseWeight, &f.Length, &f.LengthStockFolded, &f.InternalCapacity, &f.RoundWeight, &f.ReloadTime, &f.ManualCycleTime, &f.SemiAutomatic, &f.BurstRounds, &f.FullAutoRateOfFire, &f.KnockDown, &f.SAB, &f.AimTimeID)
		if err != nil {
			return nil, err
		}

		firearms = append(firearms, f)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return firearms, nil
}

func (s *State) GetEquipment() ([]Equipment, error) {
	rows, err := s.db.Query("SELECT id, name, weight_lbs FROM equipment")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	equipment := []Equipment{}

	for rows.Next() {
		e := Equipment{}
		err := rows.Scan(&e.ID, &e.Name, &e.Weight)
		if err != nil {
			return nil, err
		}

		equipment = append(equipment, e)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return equipment, nil
}

func (s *State) GetUniforms() ([]Uniform, error) {
	rows, err := s.db.Query("SELECT id, name, weight_lbs FROM uniform")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	uniforms := []Uniform{}

	for rows.Next() {
		u := Uniform{}
		err := rows.Scan(&u.ID, &u.Name, &u.Weight)
		if err != nil {
			return nil, err
		}

		uniforms = append(uniforms, u)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return uniforms, nil
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

	return &uniform, nil
}

func (s *State) GetEquipmentByIDs(equipmentIDs []int) ([]Equipment, error) {
	if len(equipmentIDs) == 0 {
		return []Equipment{}, nil
	}

	rows, err := s.db.Query(fmt.Sprintf("SELECT id, name, weight_lbs FROM equipment WHERE id IN (%s)", createPlaceholdersFromSliceLength(len(equipmentIDs))), convertIntSliceToQueryVariadicArgs(equipmentIDs)...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	equipment := []Equipment{}

	for rows.Next() {
		e := Equipment{}
		err := rows.Scan(&e.ID, &e.Name, &e.Weight)
		if err != nil {
			return nil, err
		}

		equipment = append(equipment, e)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return equipment, nil
}

func createPlaceholdersFromSliceLength(length int) string {
	p := ""
	for i := 0; i < length; i++ {
		p = p + "?,"
	}

	return strings.TrimSuffix(p, ",")
}

func convertIntSliceToQueryVariadicArgs(intSlice []int) []any {
	converted := make([]any, len(intSlice))
	for i, v := range intSlice {
		converted[i] = v
	}

	return converted
}
