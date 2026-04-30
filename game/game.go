package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"travel_survival_game/player"
	"travel_survival_game/resource"
	"travel_survival_game/world"
)

const LOOK_DISTANCE = 15

var prevCommand string

func RunComand(w *world.World, p *player.Player, c string) bool {
	input := strings.TrimSpace(c)

	parts := strings.Split(input, " ")

	switch parts[0] {
	case "look", "l":
		display(w, p)
		prevCommand = c

	case "move", "m":
		if len(parts) < 2 {
			fmt.Println("Usage: move <direction>")
		} else {
			move(w, p, parts[1])
			amap(w, p)
			if p.HP <= 20 {
				fmt.Println("Starving, health:", p.HP)
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
				fmt.Print(p.Icon) // гравець
				continue
			}

			b := w.GetBiome(x, y)
			fmt.Print(string(b.Symbol))
		}
		fmt.Println()
	}

	fmt.Println("\nYou are in:", w.GetBiome(p.X, p.Y).Name)

}

func stat(p *player.Player) {
	fmt.Printf("Health: %v\n", int(p.HP))
	fmt.Printf("Hunger: %v\n", int(p.Hunger))
	fmt.Printf("Resourses: \n")
	for _, res := range p.Resources {
		fmt.Printf("%s:\t%v\n", res.Type.Name, res.Cnt)
	}

}

func move(w *world.World, p *player.Player, dir string) {
	switch dir {
	case "north":
		p.Y--
	case "south":
		p.Y++
	case "west":
		p.X--
	case "east":
		p.X++
	default:
		fmt.Println("Unknown direction")
		return
	}

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
