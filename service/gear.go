package service

import (
	"fmt"
	"reflect"

	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/state"
)

type GearService struct {
	state equipmentState
}

type equipmentState interface {
	GetEquipment() ([]state.Equipment, error)
	GetFirearms() ([]state.Firearm, error)
	GetFirearmByID(ID int) (*state.Firearm, error)
	GetFirearmAimTimesByID(aimTimeID int) (*state.AimTime, error)
	GetFirearmBallistics(firearmID int) ([]state.Ballistics, error)
	GetFirearmAmmo(firearmID int) ([]state.Ammunition, error)
	GetFirearmMagazines(firearmID int) ([]state.Magazine, error)
}

func NewGearService(state equipmentState) *GearService {
	return &GearService{
		state: state,
	}
}

func (g *GearService) GetEquipment() ([]state.Equipment, error) {
	equipment, err := g.state.GetEquipment()
	return equipment, err
}

func (g *GearService) GetFirearms() ([]domain.FirearmDTO, error) {
	rawFirearms, err := g.state.GetFirearms()
	if err != nil {
		return []domain.FirearmDTO{}, err
	}

	mapped := []domain.FirearmDTO{}

	for _, f := range rawFirearms {
		mapped = append(mapped, mapStateFirearmToDomainFirearm(f))
	}

	return mapped, nil
}

func (g *GearService) GetEnrichedFirearmByID(id int) (*domain.EnrichedFirearmDTO, error) {
	firearm, err := g.state.GetFirearmByID(id)
	if err != nil {
		return nil, err
	}
	if firearm == nil {
		return nil, nil
	}

	aims, err := g.state.GetFirearmAimTimesByID(firearm.AimTimeID)
	if err != nil {
		return nil, err
	}
	if aims == nil {
		return nil, nil
	}

	ballistics, err := g.state.GetFirearmBallistics(id)
	if err != nil {
		return nil, err
	}

	ammo, err := g.state.GetFirearmAmmo(id)
	if err != nil {
		return nil, err
	}

	mags, err := g.state.GetFirearmMagazines(id)
	if err != nil {
		return nil, err
	}

	enriched, err := mapStateToEnrichedFirearm(*firearm, *aims, ballistics, ammo, mags)
	if err != nil {
		// fmt.Printf("err: %v\n", err)
		return nil, err
	}

	return &enriched, nil
}

func mapStateFirearmToDomainFirearm(raw state.Firearm) domain.FirearmDTO {
	mapped := domain.FirearmDTO{
		ID:                 raw.ID,
		Name:               raw.Name,
		Calibre:            raw.Calibre,
		BaseWeight:         raw.BaseWeight,
		Length:             raw.Length,
		LengthStockFolded:  raw.LengthStockFolded,
		InternalCapacity:   raw.InternalCapacity,
		RoundWeight:        raw.RoundWeight,
		ReloadTime:         raw.ReloadTime,
		ManualCycleTime:    raw.ManualCycleTime,
		SemiAutomatic:      raw.SemiAutomatic,
		BurstRounds:        raw.BurstRounds,
		FullAutoRateOfFire: raw.FullAutoRateOfFire,
		KnockDown:          raw.KnockDown,
		SAB:                raw.SAB,
	}
	return mapped
}

func mapStateToEnrichedFirearm(rawFirearm state.Firearm, rawAims state.AimTime, rawBallistics []state.Ballistics, rawAmmo []state.Ammunition, rawMags []state.Magazine) (domain.EnrichedFirearmDTO, error) {
	firearmDTO := mapStateFirearmToDomainFirearm(rawFirearm)

	ballistics, err := mapBallistics(rawBallistics)
	if err != nil {
		return domain.EnrichedFirearmDTO{}, err
	}

	enriched := domain.EnrichedFirearmDTO{
		FirearmDTO:         firearmDTO,
		AimTimes:           mapAimTimesToIntSlice(rawAims),
		RangeBrackets:      []int{10, 20, 40, 70, 100, 200, 300, 400},
		Ballistics:         ballistics,
		Ammunition:         mapAmmo(rawAmmo),
		AvailableMagazines: mapMags(rawMags),
	}

	return enriched, nil
}

func mapAimTimesToIntSlice(rawAims state.AimTime) []int {
	aims := []int{}

	values := reflect.Indirect(reflect.ValueOf(rawAims))
	for i := 1; i < values.NumField(); i++ {
		f := values.Field(i)
		switch {
		case f.Kind() == reflect.Int:
			aims = append(aims, int(f.Int()))
		case f.Kind() == reflect.Ptr && f.Type().Elem().Kind() == reflect.Int:
			if f.IsNil() {
				break
			} else {
				aims = append(aims, int(f.Elem().Int()))
			}
		}

	}
	return aims
}

func mapBallistics(rawBallistics []state.Ballistics) (domain.Ballistics, error) {
	b := domain.Ballistics{}

	for _, r := range rawBallistics {
		b.BalisticAccuracy = append(b.BalisticAccuracy, r.BA)
		b.TimeOfFlight = append(b.TimeOfFlight, r.TOF)

		if r.MA != nil {
			b.MinimumArc = append(b.MinimumArc, *r.MA)
		}
	}

	if len(b.MinimumArc) > 0 && len(b.MinimumArc) != len(b.BalisticAccuracy) {
		return b, fmt.Errorf("unexpected ballistic data")
	}

	return b, nil
}

func mapAmmo(rawAmmo []state.Ammunition) []domain.Ammunition {
	a := []domain.Ammunition{}
	lastType := ""
	index := -1
	for _, r := range rawAmmo {
		if lastType != r.Type {
			a = append(a, domain.Ammunition{Type: r.Type})
			lastType = r.Type
			index++
		}
		a[index].Pen = append(a[index].Pen, r.Pen)
		a[index].DC = append(a[index].DC, r.DC)
	}

	return a
}

func mapMags(rawMags []state.Magazine) []domain.Magazine {
	mags := []domain.Magazine{}
	for _, r := range rawMags {
		m := domain.Magazine{
			ID:       r.ID,
			Name:     r.Name,
			Type:     r.FeedType,
			Capacity: r.Capacity,
			Weight:   r.Weight,
		}
		mags = append(mags, m)
	}
	return mags
}
