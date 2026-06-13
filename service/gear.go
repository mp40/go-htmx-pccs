package service

import "github.com/mp40/go-htmx-pccs/state"

type GearService struct {
	state State
}

type State interface {
	GetEquipment() ([]state.Equipment, error)
}

func NewGearService(state State) *GearService {
	return &GearService{
		state: state,
	}
}

func (g *GearService) GetEquipment() ([]state.Equipment, error) {
	equipment, err := g.state.GetEquipment()
	return equipment, err
}
