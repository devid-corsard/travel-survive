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
	"travel_survival_game/resource"
	"travel_survival_game/world"
)

const LOOK_DISTANCE = 15

var (
	prevCommand string
	isGoing     bool = false
)

func RunComand(w *world.World, p *player.Player, c string) bool {
	input := strings.TrimSpace(c)
	if isGoing {
		isGoing = false
	}

	args := strings.Split(input, " ")

	switch args[0] {
	case "look", "l":
		display(w, p)
		prevCommand = c

	case "move", "m":
		if len(args) < 2 {
			fmt.Println("Usage: move <direction>")
		} else {
			dir, err := world.NewDirection(args[1])
			if err != nil {
				fmt.Println(err.Error())
				break
			}
			movev2(w, p, dir, 1)
			amap(w, p)
			if p.HP <= 20 {
				fmt.Println("Starving, health:", int(p.HP))
			}
			prevCommand = c
		}

	case "exit", "q":
		return true

	case "gather", "g":
		res := gather(w, p)
		if res == nil {
			fmt.Println("No resorces")
		} else {
			fmt.Printf("You found %s\n", res.Name)
		}
		stat(p)
		prevCommand = c

	case "eat", "e":
		ok := p.Eat()
		if !ok {
			fmt.Println("Nothing to eat")
		} else {
			stat(p)
		}
		prevCommand = c

	case "craft", "c":
		if len(args) < 2 {
			fmt.Println("Usage: craft <item_name>")
		} else {
			craft(p, args[1])
		}
	case "go":
		if len(args) < 2 {
			fmt.Println("Usage: go <direction>")
		} else {
			toGo(w, p, args[1])
		}

	default:
		if prevCommand != "" {
			RunComand(w, p, prevCommand)
		} else {
			fmt.Println("Unknown command")
		}
	}
	return false

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

func display(w *world.World, p *player.Player) {
	amap(w, p)
	stat(p)
}

func amap(w *world.World, p *player.Player) {
	half := LOOK_DISTANCE / 2

	fmt.Println()

	for dy := -half; dy < half; dy++ {
		for dx := -half; dx < half; dx++ {
			x := p.X + dx
			y := p.Y + dy

			if dx == 0 && dy == 0 {
				ic := p.Icon
				if isGoing {
					ic = "⛵️"
				}
				fmt.Print(ic) // гравець
				continue
			}

			b := w.GetBiome(x, y)
			fmt.Print(string(b.Symbol))
		}
		fmt.Println()
	}

	fmt.Println("\nYou are in:", w.GetBiome(p.X, p.Y).Name)
	fmt.Println("Wind: ", w.Wind.String())
}

func stat(p *player.Player) {
	fmt.Printf("Health: %v\n", int(p.HP))
	fmt.Printf("Hunger: %v\n", int(p.Hunger))
	fmt.Printf("Resourses: \n")
	for _, res := range p.Resources {
		fmt.Printf("%s:\t%v\n", res.Type.Name, res.Cnt)
	}
	fmt.Printf("Items: \n")
	for _, itm := range p.Items {
		fmt.Println(itm.Describe())
	}

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

func toGo(w *world.World, p *player.Player, dir string) {
	var boat *items.Boat
	for _, itm := range p.Items {
		switch itm.(type) {
		case *items.Boat:
			boat = itm.(*items.Boat)
		default:
		}
	}
	if boat == nil {
		fmt.Println("you dont have transport")
	} else {
		d, err := world.NewDirection(dir)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		isGoing = true
		go func() {
			for isGoing {
				speed, err := boat.Go(w.Wind, d)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					movev2(w, p, d, speed)
					amap(w, p)
				}
				<-time.After(time.Second * 3)
			}
		}()
	}
}

func craft(p *player.Player, itm string) {
	switch itm {
	case "boat":
		err := p.Craft(&items.Boat{})
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		stat(p)
	default:
		fmt.Println("unknown item: ", itm)
	}
}
