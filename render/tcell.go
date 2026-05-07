package render

import (
	"fmt"
	"strings"
	"travel_survival_game/player"
	"travel_survival_game/world"

	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Tcell struct {
	lookDistance int
	screen       tcell.Screen
}

func NewTcell(ld int, screen tcell.Screen) Renderer {
	if ld <= 0 {
		ld = 10
	}
	return &Tcell{ld, screen}
}

// Display implements Renderer.
func (t *Tcell) Display(w *world.World, p *player.Player, info string) {
	t.amap(w, p)
	t.stat(p, info)
	t.screen.Show()
}

func (t *Tcell) amap(w *world.World, p *player.Player) {
	ws, hs := t.screen.Size()
	halfX := (ws - 2) / 2
	halfY := (hs - 6) / 2

	for dy := -halfY; dy < halfY; dy++ {
		for dx := -halfX; dx < halfX; dx++ {
			x := p.X + dx
			y := p.Y + dy

			screenX := (dx + halfX) * 2
			screenY := dy + halfY

			if dx == 0 && dy == 0 {
				t.screen.Put(screenX/2, screenY, p.Icon, tcell.StyleDefault.Foreground(color.Red))
				continue
			}

			b := w.GetBiome(x, y)
			style := tcell.StyleDefault.Foreground(
				tcell.NewRGBColor(int32(b.Color.R), int32(b.Color.G), int32(b.Color.B)),
			)

			t.screen.Put(screenX, screenY, string(b.Symbol), style)
		}
	}
	t.drawText(0, hs-6, ws, hs-6, tcell.StyleDefault, fmt.Sprintf("You are in: %-10s Wind: %-15s", w.GetBiome(p.X, p.Y).Name, w.Wind.String()))

}

func (t *Tcell) stat(p *player.Player, info string) {
	ws, hs := t.screen.Size()
	t.drawText(0, hs-5, ws, hs-5, tcell.StyleDefault, fmt.Sprintf("Health: %v Hunger: %v\n", int(p.HP), int(p.Hunger)))
	var sb strings.Builder
	sb.WriteString("Resourses: ")
	for _, res := range p.Resources {
		sb.WriteString(fmt.Sprintf("%s(%v) ", res.Type.Name, res.Cnt))
	}
	t.drawText(0, hs-4, ws, hs-4, tcell.StyleDefault, sb.String())

	sb.Reset()
	sb.WriteString("Items: ")
	for _, itm := range p.Items {
		sb.WriteString(itm.Describe())
	}
	t.drawText(0, hs-3, ws, hs-3, tcell.StyleDefault, sb.String())

	t.drawText(0, hs-2, ws, hs-2, tcell.StyleDefault, info)
}

func (t *Tcell) drawText(x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	var width int
	for text != "" {
		text, width = t.screen.Put(col, row, text, style)
		col += width
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
		if width == 0 {
			// incomplete grapheme at end of string
			break
		}
	}
}
