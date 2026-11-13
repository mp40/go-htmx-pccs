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
				BallisticAccuracy: []int{61, 60, 59, 55, 50, 45, 40, 30},
				TimeOfFlight:      []int{0, 0, 0, 1, 2, 3, 4, 5},
				MinimumArc:        []float32{0.3, 0.6, 1, 2, 3, 6, 9, 12},
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

	t.Run("it returns an enriched shotgun on success", func(t *testing.T) {
		state := &stubGearState{
			firearm: &state.Firearm{Name: "TEST SHOTGUN"},
			aimTime: &state.AimTime{AimTime1: 20, AimTime2: 15, AimTime3: 10, AimTime4: 5},
			ballistics: []state.Ballistics{
				{Range: 1, BA: 61, TOF: 0},
				{Range: 2, BA: 52, TOF: 0},
				{Range: 4, BA: 42, TOF: 0},
				{Range: 6, BA: 36, TOF: 0},
				{Range: 8, BA: 32, TOF: 0},
				{Range: 10, BA: 29, TOF: 0},
				{Range: 15, BA: 23, TOF: 1},
				{Range: 20, BA: 19, TOF: 1},
				{Range: 30, BA: 13, TOF: 1},
				{Range: 40, BA: 9, TOF: 2},
				{Range: 80, BA: -1, TOF: 4},
			},
			ammunition: []state.Ammunition{
				{Type: "SHOT", Range: 1, Pen: 1.6, DC: 10, SALM: intToPointer(-13), BasePelletHitChance: nil, PelletHits: nil, PelletRadius: floatToPointer(0)},
				{Type: "SHOT", Range: 2, Pen: 1.6, DC: 4, SALM: intToPointer(-8), BasePelletHitChance: nil, PelletHits: intToPointer(11), PelletRadius: floatToPointer(0)},
				{Type: "SHOT", Range: 4, Pen: 1.6, DC: 3, SALM: intToPointer(-3), BasePelletHitChance: nil, PelletHits: intToPointer(10), PelletRadius: floatToPointer(0)},
				{Type: "SHOT", Range: 6, Pen: 1.5, DC: 3, SALM: intToPointer(0), BasePelletHitChance: nil, PelletHits: intToPointer(8), PelletRadius: floatToPointer(0.1)},
				{Type: "SHOT", Range: 8, Pen: 1.4, DC: 3, SALM: intToPointer(2), BasePelletHitChance: nil, PelletHits: intToPointer(5), PelletRadius: floatToPointer(0.1)},
				{Type: "SHOT", Range: 10, Pen: 1.4, DC: 4, SALM: intToPointer(4), BasePelletHitChance: nil, PelletHits: intToPointer(3), PelletRadius: floatToPointer(0.1)},
				{Type: "SHOT", Range: 15, Pen: 1.3, DC: 3, SALM: intToPointer(7), BasePelletHitChance: nil, PelletHits: intToPointer(1), PelletRadius: floatToPointer(0.2)},
				{Type: "SHOT", Range: 20, Pen: 1.2, DC: 2, SALM: intToPointer(9), BasePelletHitChance: intToPointer(74), PelletHits: nil, PelletRadius: floatToPointer(0.3)},
				{Type: "SHOT", Range: 30, Pen: 1.1, DC: 2, SALM: intToPointer(11), BasePelletHitChance: intToPointer(32), PelletHits: nil, PelletRadius: floatToPointer(0.4)},
				{Type: "SHOT", Range: 40, Pen: 0.9, DC: 2, SALM: intToPointer(14), BasePelletHitChance: intToPointer(17), PelletHits: nil, PelletRadius: floatToPointer(0.6)},
				{Type: "SHOT", Range: 80, Pen: 0.5, DC: 1, SALM: intToPointer(19), BasePelletHitChance: intToPointer(3), PelletHits: nil, PelletRadius: floatToPointer(1.2)},
				{Type: "SLUG", Range: 1, Pen: 3.6, DC: 9},
				{Type: "SLUG", Range: 2, Pen: 3.6, DC: 9},
				{Type: "SLUG", Range: 4, Pen: 3.5, DC: 9},
				{Type: "SLUG", Range: 6, Pen: 3.5, DC: 9},
				{Type: "SLUG", Range: 8, Pen: 3.4, DC: 9},
				{Type: "SLUG", Range: 10, Pen: 3.4, DC: 9},
				{Type: "SLUG", Range: 15, Pen: 3.2, DC: 9},
				{Type: "SLUG", Range: 20, Pen: 3.1, DC: 9},
				{Type: "SLUG", Range: 30, Pen: 2.8, DC: 8},
				{Type: "SLUG", Range: 40, Pen: 2.6, DC: 8},
				{Type: "SLUG", Range: 80, Pen: 1.8, DC: 6},
			},
			magazines: []state.Magazine{},
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
			FirearmDTO:    domain.FirearmDTO{Name: "TEST SHOTGUN"},
			AimTimes:      []int{20, 15, 10, 5},
			RangeBrackets: []int{1, 2, 4, 6, 8, 10, 15, 20, 30, 40, 80},
			Ballistics: domain.Ballistics{
				BallisticAccuracy: []int{61, 52, 42, 36, 32, 29, 23, 19, 13, 9, -1},
				TimeOfFlight:      []int{0, 0, 0, 0, 0, 0, 1, 1, 1, 2, 4},
			},
			Ammunition: []domain.Ammunition{
				{
					Type: "SHOT",
					Pen:  []float32{1.6, 1.6, 1.6, 1.5, 1.4, 1.4, 1.3, 1.2, 1.1, 0.9, 0.5},
					DC:   []int{10, 4, 3, 3, 3, 4, 3, 2, 2, 2, 1},
					Pellet: &domain.Pellet{
						SALM:                []int{-13, -8, -3, 0, 2, 4, 7, 9, 11, 14, 19},
						BasePelletHitChance: []*int{nil, nil, nil, nil, nil, nil, nil, intToPointer(74), intToPointer(32), intToPointer(17), intToPointer(3)},
						BasePelletHits:      []*int{nil, intToPointer(11), intToPointer(10), intToPointer(8), intToPointer(5), intToPointer(3), intToPointer(1), nil, nil, nil, nil},
						PelletRadius:        []float32{0, 0, 0, 0.1, 0.1, 0.1, 0.2, 0.3, 0.4, 0.6, 1.2},
					},
				},
				{
					Type: "SLUG",
					Pen:  []float32{3.6, 3.6, 3.5, 3.5, 3.4, 3.4, 3.2, 3.1, 2.8, 2.6, 1.8},
					DC:   []int{9, 9, 9, 9, 9, 9, 9, 9, 8, 8, 6},
				},
			},
			AvailableMagazines: []domain.Magazine{},
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
