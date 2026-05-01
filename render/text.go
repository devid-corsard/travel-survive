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
	t.amap(w, p)
	t.stat(p)
	fmt.Println("Info: ", info)
	fmt.Print("\n> ")
}

func (t *Text) amap(w *world.World, p *player.Player) {
	half := t.lookDistance / 2

	fmt.Println()

	for dy := -half; dy < half; dy++ {
		for dx := -half; dx < half; dx++ {
			x := p.X + dx
			y := p.Y + dy

			if dx == 0 && dy == 0 {
				ic := p.Icon
				if false {
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

func (t *Text) stat(p *player.Player) {
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
