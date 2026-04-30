package player

import (
	"fmt"
	"travel_survival_game/items"
	"travel_survival_game/resource"
)

type Player struct {
	X         int                 `json:"x"`
	Y         int                 `json:"y"`
	HP        float64             `json:"hp"`
	Hunger    float64             `json:"hunger"`
	Resources []resource.Resource `json:"resources"`
	Items     items.ItemList      `json:"items"`
	Icon      string              `json:"icon"`
}

const (
	DEFAULT_HP                 float64 = 100
	HOURS_CAN_LIVE_WITH_HUNGER float64 = 240
	EXHAUST_K                          = DEFAULT_HP / HOURS_CAN_LIVE_WITH_HUNGER
	FOOD_EFFECT                float64 = 10
)

func New(x, y int) *Player {
	return &Player{x, y, DEFAULT_HP, 0, make([]resource.Resource, 0), make([]items.Item, 0), "😊"}
}

func (this *Player) Live(hours float64) {
	if hours <= 0 {
		hours = 1
	}
	exhaust := hours * EXHAUST_K
	if this.Hunger > 99.99 {
		this.HP -= exhaust
	} else {
		this.Hunger += exhaust
	}
}

func (this *Player) Eat() bool {
	for i, r := range this.Resources {
		if r.Type.Eatable && r.Cnt > 0 {
			this.Resources[i].Cnt--
			if this.HP < 100 {
				this.HP += FOOD_EFFECT
				if this.HP > 100 {
					this.HP = 100
				}
			}
			if this.Hunger > 0 {
				this.Hunger -= FOOD_EFFECT
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

func (this *Player) Craft(itm items.Item) error {
	switch itm.(type) {
	case *items.Boat:
		boat, err := items.NewBoat(this.Resources)
		if err != nil {
			return err
		}
		this.Items = append(this.Items, boat)
		return nil
	}
	return fmt.Errorf("Unknown item: %v", itm.Describe())
}
