package player

import "travel_survival_game/resource"

type Player struct {
	X         int                 `json:"x"`
	Y         int                 `json:"y"`
	HP        uint                `json:"hp"`
	Hunger    int                 `json:"hunger"`
	Resources []resource.Resource `json:"resources"`
}

func New(x, y int) *Player {
	return &Player{x, y, 100, 0, make([]resource.Resource, 0)}
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
	for _, r := range this.Resources {
		if r.Type.Eatable && r.Cnt > 0 {
			r.Cnt--
			return true
		}
	}
	return false
}

func (this *Player) Collect(r resource.Type) {
	for i, rs := range this.Resources {
		if rs.Type.Name == r.Name {
			this.Resources[i].Cnt++
			return
		}
	}
	this.Resources = append(this.Resources, resource.Resource{r, 1})
}
