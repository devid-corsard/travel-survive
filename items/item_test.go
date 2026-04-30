package items

import (
	"testing"
	"travel_survival_game/world"
)

func TestBoatGo_ErrorOnSame(t *testing.T) {
	b := &Boat{}
	_, err := b.Go(world.North, world.North)

	if err == nil {
		t.Errorf("err nil want not nil for same wind and cource directions")
	}
}
