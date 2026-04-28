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
	for i, r := range this.Resources {
		if r.Type.Eatable && r.Cnt > 0 {
			this.Resources[i].Cnt--
			if this.HP < 100 {
				this.HP += 5
				if this.HP > 100 {
					this.HP = 100
				}
			}
			if this.Hunger > 0 {
				this.Hunger -= 5
				if this.Hunger < 0 {
					this.Hunger = 0
				}
			}
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
