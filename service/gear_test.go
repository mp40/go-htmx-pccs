package service

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mp40/go-htmx-pccs/domain"
	"github.com/mp40/go-htmx-pccs/state"
)

type stubGearState struct {
	firearms   []state.Firearm
	firearm    *state.Firearm
	aimTime    *state.AimTime
	ballistics []state.Ballistics
	ammunition []state.Ammunition
	magazines  []state.Magazine
	err        error
}

func (g *stubGearState) GetFirearms() ([]state.Firearm, error) {
	return g.firearms, g.err
}

func (g *stubGearState) GetFirearmByID(id int) (*state.Firearm, error) {
	return g.firearm, g.err
}

func (g *stubGearState) GetFirearmAimTimesByID(aimTimeID int) (*state.AimTime, error) {
	return g.aimTime, g.err
}

func (g *stubGearState) GetFirearmBallistics(firearmID int) ([]state.Ballistics, error) {
	return g.ballistics, g.err
}

func (g *stubGearState) GetFirearmAmmo(firearmID int) ([]state.Ammunition, error) {
	return g.ammunition, g.err
}

func (g *stubGearState) GetFirearmMagazines(firearmID int) ([]state.Magazine, error) {
	return g.magazines, g.err
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

func TestGetEnrichedFirearmByID(t *testing.T) {
	t.Run("it returns an enriched firearm on success", func(t *testing.T) {
		state := &stubGearState{
			firearm: &state.Firearm{Name: "TEST FIREARM", Calibre: "6x9", BaseWeight: 7.7},
			aimTime: &state.AimTime{
				AimTime1:  20,
				AimTime2:  15,
				AimTime3:  14,
				AimTime4:  13,
				AimTime5:  intToPointer(12),
				AimTime6:  intToPointer(11),
				AimTime7:  intToPointer(10),
				AimTime8:  intToPointer(9),
				AimTime9:  intToPointer(8),
				AimTime10: nil,
				AimTime11: nil,
				AimTime12: nil,
			},
			ballistics: []state.Ballistics{
				{FirearmID: 1, Range: 10, BA: 61, TOF: 0, MA: floatToPointer(0.3)},
				{FirearmID: 1, Range: 20, BA: 60, TOF: 0, MA: floatToPointer(0.6)},
				{FirearmID: 1, Range: 40, BA: 59, TOF: 0, MA: floatToPointer(1)},
				{FirearmID: 1, Range: 70, BA: 55, TOF: 1, MA: floatToPointer(2)},
				{FirearmID: 1, Range: 100, BA: 50, TOF: 2, MA: floatToPointer(3)},
				{FirearmID: 1, Range: 200, BA: 45, TOF: 3, MA: floatToPointer(6)},
				{FirearmID: 1, Range: 300, BA: 40, TOF: 4, MA: floatToPointer(9)},
				{FirearmID: 1, Range: 400, BA: 30, TOF: 5, MA: floatToPointer(12)},
			},
			ammunition: []state.Ammunition{
				{Type: "FMJ", Range: 10, Pen: 15, DC: 6},
				{Type: "FMJ", Range: 20, Pen: 15, DC: 6},
				{Type: "FMJ", Range: 40, Pen: 14, DC: 6},
				{Type: "FMJ", Range: 70, Pen: 12, DC: 5},
				{Type: "FMJ", Range: 100, Pen: 11, DC: 5},
				{Type: "FMJ", Range: 200, Pen: 7, DC: 4},
				{Type: "FMJ", Range: 300, Pen: 4.6, DC: 3},
				{Type: "FMJ", Range: 400, Pen: 3, DC: 2},
				{Type: "JHP", Range: 10, Pen: 15, DC: 8},
				{Type: "JHP", Range: 20, Pen: 14, DC: 8},
				{Type: "JHP", Range: 40, Pen: 13, DC: 7},
				{Type: "JHP", Range: 70, Pen: 10, DC: 7},
				{Type: "JHP", Range: 100, Pen: 6.7, DC: 7},
				{Type: "JHP", Range: 200, Pen: 4.4, DC: 6},
				{Type: "JHP", Range: 300, Pen: 3.3, DC: 5},
				{Type: "JHP", Range: 400, Pen: 2, DC: 3},
			},
			magazines: []state.Magazine{
				{Name: "STANAG", FeedType: "Mag", Capacity: 20, Weight: 0.7},
				{Name: "STANAG", FeedType: "Mag", Capacity: 30, Weight: 1},
			},
		}
		service := NewGearService(state)

		got, err := service.GetEnrichedFirearmByID(1)
		if err != nil {
			t.Errorf("unexpected error, got %v", err)
		}
		if got == nil {
			t.Errorf("firearm is <nil>")
		}

		want := domain.EnrichedFirearmDTO{
			FirearmDTO: domain.FirearmDTO{
				Name:       "TEST FIREARM",
				Calibre:    "6x9",
				BaseWeight: 7.7,
			},
			AimTimes:      []int{20, 15, 14, 13, 12, 11, 10, 9, 8},
			RangeBrackets: []int{10, 20, 40, 70, 100, 200, 300, 400},
			Ballistics: domain.Ballistics{
				BalisticAccuracy: []int{61, 60, 59, 55, 50, 45, 40, 30},
				TimeOfFlight:     []int{0, 0, 0, 1, 2, 3, 4, 5},
				MinimumArc:       []float32{0.3, 0.6, 1, 2, 3, 6, 9, 12},
			},
			Ammunition: []domain.Ammunition{
				{
					Type: "FMJ",
					Pen:  []float32{15, 15, 14, 12, 11, 7, 4.6, 3},
					DC:   []int{6, 6, 6, 5, 5, 4, 3, 2},
				},
				{
					Type: "JHP",
					Pen:  []float32{15, 14, 13, 10, 6.7, 4.4, 3.3, 2},
					DC:   []int{8, 8, 7, 7, 7, 6, 5, 3},
				},
			},
			AvailableMagazines: []domain.Magazine{
				{Name: "STANAG", Type: "Mag", Capacity: 20, Weight: 0.7},
				{Name: "STANAG", Type: "Mag", Capacity: 30, Weight: 1},
			},
		}

		if !cmp.Equal(*got, want) {
			t.Errorf("unexpected enriched stats: %v", cmp.Diff(*got, want))
		}
	})
}

func TestMapAimTimesToIntSlice(t *testing.T) {
	t.Run("it maps manidtory aims with all nil optional", func(t *testing.T) {
		aims := state.AimTime{
			AimTime1: 20,
			AimTime2: 15,
			AimTime3: 10,
			AimTime4: 5,
		}

		got := mapAimTimesToIntSlice(aims)

		want := []int{20, 15, 10, 5}

		if !cmp.Equal(got, want) {
			t.Errorf("unexpected result: %v", cmp.Diff(got, want))
		}
	})

	t.Run("it maps manidtory aims with partial optional non nil", func(t *testing.T) {
		aims := state.AimTime{
			AimTime1: 20,
			AimTime2: 15,
			AimTime3: 10,
			AimTime4: 5,
			AimTime5: intToPointer(4),
			AimTime6: intToPointer(3),
			AimTime7: intToPointer(2),
		}

		got := mapAimTimesToIntSlice(aims)

		want := []int{20, 15, 10, 5, 4, 3, 2}

		if !cmp.Equal(got, want) {
			t.Errorf("unexpected result: %v", cmp.Diff(got, want))
		}
	})

	t.Run("it maps manidtory aims with alll optional non nil", func(t *testing.T) {
		aims := state.AimTime{
			AimTime1:  -10,
			AimTime2:  -5,
			AimTime3:  0,
			AimTime4:  5,
			AimTime5:  intToPointer(10),
			AimTime6:  intToPointer(11),
			AimTime7:  intToPointer(12),
			AimTime8:  intToPointer(13),
			AimTime9:  intToPointer(14),
			AimTime10: intToPointer(15),
			AimTime11: intToPointer(16),
			AimTime12: intToPointer(17),
		}

		got := mapAimTimesToIntSlice(aims)

		want := []int{-10, -5, 0, 5, 10, 11, 12, 13, 14, 15, 16, 17}

		if !cmp.Equal(got, want) {
			t.Errorf("unexpected result: %v", cmp.Diff(got, want))
		}
	})
}

func TestMapAmmo(t *testing.T) {
	t.Run("it should map one ammo type", func(t *testing.T) {
		state := []state.Ammunition{
			{Type: "FMJ", Pen: 3, DC: 5},
			{Type: "FMJ", Pen: 2, DC: 4},
			{Type: "FMJ", Pen: 1, DC: 3},
		}

		got := mapAmmo(state)

		want := []domain.Ammunition{
			{
				Type: "FMJ",
				Pen:  []float32{3, 2, 1},
				DC:   []int{5, 4, 3},
			},
		}

		if !cmp.Equal(got, want) {
			t.Errorf("unexpected result: %v", cmp.Diff(got, want))
		}
	})

	t.Run("it should map multiple ammo types", func(t *testing.T) {
		state := []state.Ammunition{
			{Type: "FMJ", Pen: 3, DC: 5},
			{Type: "FMJ", Pen: 2, DC: 4},
			{Type: "FMJ", Pen: 1, DC: 3},
			{Type: "JHP", Pen: 2, DC: 7},
			{Type: "JHP", Pen: 1, DC: 6},
			{Type: "JHP", Pen: 1, DC: 5},
			{Type: "AP", Pen: 5, DC: 3},
			{Type: "AP", Pen: 4, DC: 2},
			{Type: "AP", Pen: 3, DC: 1},
		}

		got := mapAmmo(state)

		want := []domain.Ammunition{
			{
				Type: "FMJ",
				Pen:  []float32{3, 2, 1},
				DC:   []int{5, 4, 3},
			},
			{
				Type: "JHP",
				Pen:  []float32{2, 1, 1},
				DC:   []int{7, 6, 5},
			},
			{
				Type: "AP",
				Pen:  []float32{5, 4, 3},
				DC:   []int{3, 2, 1},
			},
		}

		if !cmp.Equal(got, want) {
			t.Errorf("unexpected result: %v", cmp.Diff(got, want))
		}
	})
}

func floatToPointer(f float32) *float32 {
	return &f
}

func intToPointer(i int) *int {
	return &i
}
