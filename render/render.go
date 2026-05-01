package render

import (
	"travel_survival_game/player"
	"travel_survival_game/world"
)

type Renderer interface {
	Display(*world.World, *player.Player, string)
}
