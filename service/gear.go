package service

import "github.com/mp40/go-htmx-pccs/state"

type GearService struct {
	state equipmentState
}

type equipmentState interface {
	GetEquipment() ([]state.Equipment, error)
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
