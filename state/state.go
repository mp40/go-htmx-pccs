package state

import (
	"database/sql"
	"fmt"
	"strings"
)

type Firearm struct {
	ID                 int      `json:"id"`
	Name               string   `json:"name"`
	Calibre            string   `json:"string"`
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

type AimTime struct {
	ID        int  `json:"id"`
	AimTime1  int  `json:"at_1"`
	AimTime2  int  `json:"at_2"`
	AimTime3  int  `json:"at_3"`
	AimTime4  int  `json:"at_4"`
	AimTime5  *int `json:"at_5"`
	AimTime6  *int `json:"at_6"`
	AimTime7  *int `json:"at_7"`
	AimTime8  *int `json:"at_8"`
	AimTime9  *int `json:"at_9"`
	AimTime10 *int `json:"at_10"`
	AimTime11 *int `json:"at_11"`
	AimTime12 *int `json:"at_12"`
}

type Ballistics struct {
	FirearmID int      `json:"firearm_id"`
	Range     int      `json:"range_hexes"`
	BA        int      `json:"ba"`
	TOF       int      `json:"tof"`
	MA        *float32 `json:"ma"`
}

type Ammunition struct {
	ID                  int      `json:"id"`
	FirearmID           int      `json:"firearm_id"`
	Range               int      `json:"range_hexes"`
	Type                string   `json:"type"`
	Pen                 float32  `json:"pen"`
	DC                  int      `json:"dc"`
	SALM                int      `json:"salm"`
	BasePelletHitChance *int     `json:"bphc"`
	PelletHits          *int     `json:"bphc_hits"`
	PelletRadius        *float32 `json:"pr"`
}

type Magazine struct {
	ID        int     `json:"magazine_id"`
	FirearmID int     `json:"firearm_id"`
	Name      string  `json:"name"`
	FeedType  string  `json:"feed_type"`
	Capacity  int     `json:"capacity"`
	Weight    float32 `json:"weight"`
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
	rows, err := s.db.Query("SELECT id, name, calibre, base_weight, length, length_stock_folded, internal_capacity, round_weight, reload_time, manual_cycle_time, semi_automatic, burst_rounds, full_auto_rof, knock_down, sab, aim_time_id FROM firearm")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	firearms := []Firearm{}

	for rows.Next() {
		f := Firearm{}
		err := rows.Scan(&f.ID, &f.Name, &f.Calibre, &f.BaseWeight, &f.Length, &f.LengthStockFolded, &f.InternalCapacity, &f.RoundWeight, &f.ReloadTime, &f.ManualCycleTime, &f.SemiAutomatic, &f.BurstRounds, &f.FullAutoRateOfFire, &f.KnockDown, &f.SAB, &f.AimTimeID)
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

func (s *State) GetFirearmByID(firearmID int) (*Firearm, error) {
	f := Firearm{}
	err := s.db.QueryRow("SELECT id, name, calibre, base_weight, length, length_stock_folded, internal_capacity, round_weight, reload_time, manual_cycle_time, semi_automatic, burst_rounds, full_auto_rof, knock_down, sab, aim_time_id FROM firearm WHERE id = ?", firearmID).
		Scan(&f.ID, &f.Name, &f.Calibre, &f.BaseWeight, &f.Length, &f.LengthStockFolded, &f.InternalCapacity, &f.RoundWeight, &f.ReloadTime, &f.ManualCycleTime, &f.SemiAutomatic, &f.BurstRounds, &f.FullAutoRateOfFire, &f.KnockDown, &f.SAB, &f.AimTimeID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &f, nil
}

func (s *State) GetFirearmAimTimesByID(aimTimeID int) (*AimTime, error) {
	at := AimTime{}
	err := s.db.QueryRow("SELECT id, at_1, at_2, at_3, at_4, at_5, at_6, at_7, at_8, at_9, at_10, at_11, at_12 FROM aim_time WHERE id = ?", aimTimeID).
		Scan(&at.ID, &at.AimTime1, &at.AimTime2, &at.AimTime3, &at.AimTime4, &at.AimTime5, &at.AimTime6, &at.AimTime7, &at.AimTime8, &at.AimTime9, &at.AimTime10, &at.AimTime11, &at.AimTime12)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &at, nil
}

func (s *State) GetFirearmBallistics(firearmID int) ([]Ballistics, error) {
	rows, err := s.db.Query("SELECT firearm_id, range_hexes, ba, tof, ma FROM firearm_ballistics_range WHERE firearm_id = ? ORDER BY range_hexes ASC", firearmID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	defer rows.Close()

	ballistics := []Ballistics{}
	for rows.Next() {
		b := Ballistics{}
		err := rows.Scan(&b.FirearmID, &b.Range, &b.MA, &b.BA, &b.TOF)
		if err != nil {
			return nil, err
		}
		ballistics = append(ballistics, b)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return ballistics, nil
}

func (s *State) GetFirearmAmmo(firearmID int) ([]Ammunition, error) {
	rows, err := s.db.Query("SELECT fa.id, fa.firearm_id, fa.type, r.range_hexes, r.pen, r.dc, r.salm, r.bphc, r.bphc_hits, r.pr FROM firearm_ammo fa INNER JOIN firearm_ammo_range r ON fa.id = r.firearm_ammo_id WHERE fa.firearm_id = ?", firearmID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ammo := []Ammunition{}

	for rows.Next() {
		a := Ammunition{}
		err := rows.Scan(&a.ID, &a.FirearmID, &a.Type, &a.Range, &a.Pen, &a.DC, &a.SALM, &a.BasePelletHitChance, &a.PelletHits, &a.PelletRadius)
		if err != nil {
			return nil, err
		}

		ammo = append(ammo, a)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return ammo, nil
}

func (s *State) GetFirearmMagazines(firearmID int) ([]Magazine, error) {
	rows, err := s.db.Query("SELECT fm.magazine_id, fm.firearm_id, m.name, m.feed_type, m.capacity, m.weight FROM firearm_magazine fm INNER JOIN magazine m ON fm.magazine_id = m.id WHERE fm.firearm_id = ?", firearmID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mags := []Magazine{}

	for rows.Next() {
		m := Magazine{}
		err := rows.Scan(&m.ID, &m.FirearmID, &m.Name, &m.FeedType, &m.Capacity, &m.Weight)
		if err != nil {
			return nil, err
		}

		mags = append(mags, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return mags, nil
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
