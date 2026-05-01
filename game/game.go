package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"travel_survival_game/items"
	"travel_survival_game/player"
	"travel_survival_game/render"
	"travel_survival_game/resource"
	"travel_survival_game/world"
)

const LOOK_DISTANCE = 35

var (
	prevCommand string
	isGoing     bool = false
)

func RunComand(w *world.World, p *player.Player, c string, ren render.Chan) (bool, string) {
	var comInfo string
	input := strings.TrimSpace(c)
	if isGoing {
		isGoing = p.SetSailing(false)
	}

	args := strings.Split(input, " ")

	switch args[0] {
	case "look", "l":
		prevCommand = c

	case "move", "m":
		if len(args) < 2 {
			comInfo = "Usage: move <direction>"
		} else {
			dir, err := world.NewDirection(args[1])
			if err != nil {
				comInfo = err.Error()
				break
			}
			movev2(w, p, dir, 1)
			if p.HP <= 20 {
				comInfo = fmt.Sprint("Starving, health:", int(p.HP))
			}
			prevCommand = c
		}

	case "exit", "q":
		return true, ""

	case "gather", "g":
		res := gather(w, p)
		if res == nil {
			comInfo = "No resorces"
		} else {
			comInfo = fmt.Sprintf("You found %s", res.Name)
		}
		prevCommand = c

	case "eat", "e":
		ok := p.Eat()
		if !ok {
			comInfo = "Nothing to eat"
		}
		prevCommand = c

	case "craft", "c":
		if len(args) < 2 {
			comInfo = "Usage: craft <item_name>"
		} else {
			comInfo = craft(p, args[1])
		}
	case "go":
		if len(args) < 2 {
			comInfo = "Usage: go <direction>"
		} else {
			comInfo = toGo(w, p, args[1], ren)
		}

	default:
		if prevCommand != "" {
			RunComand(w, p, prevCommand, ren)
		} else {
			comInfo = "Unknown command"
		}
	}
	return false, comInfo

}

func Save(p *player.Player) error {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return os.WriteFile("save.json", b, 0644)
}

func Restore() (*player.Player, error) {
	b, err := os.ReadFile("save.json")
	if err != nil {
		return nil, err
	}
	var p player.Player
	err = json.Unmarshal(b, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func movev2(w *world.World, p *player.Player, dir world.Direction, speed int) {
	dx, dy := dir.Vector()
	p.X += dx * speed
	p.Y += dy * speed

	// wrap
	if p.X < 0 {
		p.X = w.WorldWidth - 1
	}
	if p.Y < 0 {
		p.Y = w.WorldHeight - 1
	}
	if p.X >= w.WorldWidth {
		p.X = 0
	}
	if p.Y >= w.WorldHeight {
		p.Y = 0
	}
	p.Live(1)
}

func gather(w *world.World, p *player.Player) *resource.Type {
	p.Live(1)
	b := w.GetBiome(p.X, p.Y)
	var commonProb uint
	if len(b.Resources) == 0 {
		return nil
	}
	for _, p := range b.Resources {
		commonProb += p
	}
	if commonProb < 100 {
		commonProb = 100
	}
	r := uint(rand.Intn(int(commonProb)))

	cumulative := uint(0)
	for res, w := range b.Resources {
		cumulative += w
		if r < cumulative {
			p.Collect(res)
			return &res
		}
	}
	return nil
}

func toGo(w *world.World, p *player.Player, dir string, ren render.Chan) string {
	var boat *items.Boat
	for _, itm := range p.Items {
		switch itm := itm.(type) {
		case *items.Boat:
			boat = itm
		}
	}
	if boat == nil {
		return "you dont have transport"
	} else {
		d, err := world.NewDirection(dir)
		if err != nil {
			return err.Error()
		}

		isGoing = p.SetSailing(true)
		go func() {
			for isGoing {
				speed, err := boat.Go(w.Wind, d)
				if err != nil {
					ren <- err.Error()
				} else {
					movev2(w, p, d, 1)
					ren <- fmt.Sprint("Sailing, speed: ", speed)
				}
				if speed <= 0 {
					speed = 1
				}
				<-time.After(time.Second * 3 / time.Duration(speed))
			}
		}()
	}
	return ""
}

func craft(p *player.Player, itm string) string {
	switch itm {
	case "boat":
		err := p.Craft(&items.Boat{})
		if err != nil {
			return err.Error()
		}
		return ""
	default:
		return fmt.Sprint("unknown item: ", itm)
	}
}
