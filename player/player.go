package player

import "travel_survival_game/resource"

type Player struct {
	X         int                   `json:"x"`
	Y         int                   `json:"y"`
	HP        uint                  `json:"hp"`
	Hunger    int                   `json:"hunger"`
	Resources map[resource.Res]uint `json:"resources"`
}

func New(x, y int) *Player {
	return &Player{x, y, 100, 0, make(map[resource.Res]uint)}
}

func (this *Player) Live(hours uint) {
	if hours == 0 {
		hours = 1
	}
	if this.Hunger > 24 {
		this.HP--
	} else {
		this.Hunger += 1 * int(hours)
	}
}

func (this *Player) Eat() bool {
	for r, c := range this.Resources {
		if r.Eatable && c > 0 {
			this.Resources[r]--
			return true
		}
	}
	return false
}

func (this *Player) Collect(r resource.Res) {
	this.Resources[r] += 1
}
