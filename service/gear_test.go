package service

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/state"
)

type stubGearState struct {
	firearms []state.Firearm
	err      error
}

func (g *stubGearState) GetFirearms() ([]state.Firearm, error) {
	return g.firearms, g.err
}

func (g *stubGearState) GetFirearmByID(id int) (*state.Firearm, error) {
	return nil, g.err
}

func (g *stubGearState) GetEquipment() ([]state.Equipment, error) {
	panic("method not used in test")
}

func TestGetFirearms(t *testing.T) {
	t.Run("it returns list of parsed firearm data transfer objects", func(t *testing.T) {
		state := &stubGearState{firearms: []state.Firearm{{ID: 1, Name: "TEST FIREARM", AimTimeID: 3}}}
		service := NewGearService(state)

		got, err := service.GetFirearms()
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}

		want := domain.FirearmDTO{
			ID:   1,
			Name: "TEST FIREARM",
		}

		if !cmp.Equal(got[0], want) {
			t.Errorf("unexpected firearm: %v", cmp.Diff(got[0], want))
		}

	})

	t.Run("it returns error on db error", func(t *testing.T) {
		state := &stubGearState{err: fmt.Errorf("sadness")}
		service := NewGearService(state)

		_, err := service.GetFirearms()
		if err == nil {
			t.Errorf("expected error, got <nil>")
		}
	})
}
