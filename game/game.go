package game

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"travel_survival_game/player"
	"travel_survival_game/world"
)

const LOOK_DISTANCE = 10

var prevCommand string

func RunComand(w *world.World, p *player.Player, c string) bool {
	input := strings.TrimSpace(c)

	parts := strings.Split(input, " ")

	switch parts[0] {
	case "look", "l":
		look(w, p)
		prevCommand = c

	case "move", "m":
		if len(parts) < 2 {
			fmt.Println("Usage: move <direction>")
		} else {
			move(w, p, parts[1])
			look(w, p)
			prevCommand = c
		}

	case "exit", "q":
		return true

	case "gather", "g":
		gather(w, p)
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
	return os.WriteFile("save.json", b, os.ModeAppend)
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

func look(w *world.World, p *player.Player) {
	half := LOOK_DISTANCE / 2

	fmt.Println()

	for dy := -half; dy < half; dy++ {
		for dx := -half; dx < half; dx++ {
			x := p.X + dx
			y := p.Y + dy

			if dx == 0 && dy == 0 {
				fmt.Print("@") // гравець
				continue
			}

			b := w.GetBiome(x, y)
			fmt.Print(string(b.Symbol))
		}
		fmt.Println()
	}

	fmt.Println("\nYou are in:", w.GetBiome(p.X, p.Y).Name)
	fmt.Printf("Health: %v\n", p.HP)
	fmt.Printf("Hunger: %v\n", p.Hunger)
	fmt.Printf("Resourses: \n")
	for res, cou := range p.Resources {
		fmt.Printf("%s:\t%v\n", res.Name, cou)
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

func gather(w *world.World, p *player.Player) {
	b := w.GetBiome(p.X, p.Y)
	var commonProb uint
	if len(b.Resources) == 0 {
		fmt.Println("No resorces")
		return
	}
	for _, p := range b.Resources {
		commonProb += p
	}
	r := uint(rand.Intn(int(commonProb)))

	cumulative := uint(0)
	for res, w := range b.Resources {
		cumulative += w
		if r < cumulative {
			p.Resources[res] += 1
			fmt.Printf("You found %s\n", res.Name)
		}
	}
}
