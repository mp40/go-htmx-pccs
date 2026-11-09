package service

import (
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
	return nil, nil
}

func mapStateFirearmToDomainFirearm(raw state.Firearm) domain.FirearmDTO {
	mapped := domain.FirearmDTO{
		ID:                 raw.ID,
		Name:               raw.Name,
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
