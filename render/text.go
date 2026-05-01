package render

import (
	"fmt"
	"travel_survival_game/player"
	"travel_survival_game/world"
)

type Text struct {
	lookDistance int
}

func NewText(ld int) Renderer {
	if ld <= 0 {
		ld = 10
	}
	return &Text{ld}
}

// Display implements Renderer.
func (t *Text) Display(w *world.World, p *player.Player, info string) {
	clearScreen()
	t.amap(w, p)
	t.stat(p)
	fmt.Println("Info: ", info)
	fmt.Print("\n> ")
}
func clearScreen() {
	fmt.Print("\033[2J") // clear screen
	fmt.Print("\033[H")  // move cursor to top-left
}

func (t *Text) amap(w *world.World, p *player.Player) {
	half := t.lookDistance / 2

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

	fmt.Printf("You are in: %v\tWind: %v\n", w.GetBiome(p.X, p.Y).Name, w.Wind.String())
}

func (t *Text) stat(p *player.Player) {
	fmt.Printf("Health: %v\tHunger: %v\n", int(p.HP), int(p.Hunger))
	fmt.Printf("Resourses: ")
	for _, res := range p.Resources {
		fmt.Printf("%s(%v) ", res.Type.Name, res.Cnt)
	}
	fmt.Println()
	fmt.Printf("Items: ")
	for _, itm := range p.Items {
		fmt.Print(itm.Describe())
	}
	fmt.Println()
}
