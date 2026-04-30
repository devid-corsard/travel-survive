package items

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"travel_survival_game/resource"
	"travel_survival_game/world"
)

type Item interface {
	From() []resource.Resource
	Describe() string
	ResString() string
}

type ItemList []Item

func (l ItemList) MarshalJSON() ([]byte, error) {
	type TypedItem struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}

	var wrappers []TypedItem
	for _, itm := range l {
		var t string
		switch itm.(type) {
		case *Boat:
			t = "boat"
		// add other cases here
		default:
			t = "unknown"
		}

		wrappers = append(wrappers, TypedItem{
			Type: t,
			Data: itm,
		})
	}
	return json.Marshal(wrappers)
}
func (l *ItemList) UnmarshalJSON(data []byte) error {
	type TypedItem struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}

	var wrappers []TypedItem
	if err := json.Unmarshal(data, &wrappers); err != nil {
		return err
	}

	for _, w := range wrappers {
		var itm Item
		switch w.Type {
		case "boat":
			var b Boat
			if err := json.Unmarshal(w.Data, &b); err != nil {
				return err
			}
			itm = &b
		}
		*l = append(*l, itm)
	}
	return nil
}

type Boat struct {
	HP float64
}

func (b *Boat) From() []resource.Resource {
	return []resource.Resource{{Type: resource.Wood, Cnt: 100}}
}

func (b *Boat) ResString() string {
	var str strings.Builder
	for _, rs := range b.From() {
		str.WriteString(fmt.Sprintf("%v: %v, ", rs.Type.Name, rs.Cnt))
	}
	return str.String()
}

func (b *Boat) Go(wind, course world.Direction) (int, error) {
	attackAngle := course.Delta(wind)
	if attackAngle == 0 {
		return 0, errors.New("Move against the wind")
	}
	speedMap := map[int]int{
		45:  1,
		90:  2,
		135: 3,
		180: 2,
	}
	speed, ok := speedMap[attackAngle]
	if !ok {
		panic(fmt.Sprintf("unexpected attack angle: %v", attackAngle))
	}
	return speed, nil
}
func (b *Boat) Describe() string {
	return fmt.Sprintf("Boat, %v%% HP", int(b.HP))
}

func NewBoat(rs []resource.Resource) (Item, error) {
	b := &Boat{100}
	for _, rr := range rs {
		if rr.Type.Name == b.From()[0].Type.Name && rr.Cnt >= b.From()[0].Cnt {
			return b, nil
		}
	}

	return nil, fmt.Errorf("Not enough resouces, need %v", b.ResString())
}
