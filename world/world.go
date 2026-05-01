package world

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"math/rand"
	"os"
	"time"
	"travel_survival_game/resource"
)

type Color struct {
	R, G, B uint8
}

type Direction int

const (
	North     Direction = iota * 45 // 0 * 45 = 0
	NorthEast                       // 1 * 45 = 45
	East                            // 2 * 45 = 90
	SouthEast                       // 3 * 45 = 135
	South                           // 4 * 45 = 180
	SouthWest                       // 5 * 45 = 225
	West                            // 6 * 45 = 270
	NorthWest                       // 7 * 45 = 315
)

func NewDirection(ds string) (Direction, error) {
	dm := map[string]Direction{
		"north":      North,
		"north-east": NorthEast,
		"east":       East,
		"south-east": SouthEast,
		"south":      South,
		"south-west": SouthWest,
		"west":       West,
		"north-west": NorthWest,
		"n":          North,
		"ne":         NorthEast,
		"e":          East,
		"se":         SouthEast,
		"s":          South,
		"sw":         SouthWest,
		"w":          West,
		"nw":         NorthWest,
	}
	d, ok := dm[ds]
	if !ok {
		return 0, fmt.Errorf("wrong direction: %v", ds)
	}
	return d, nil
}

// Normalize the difference to be within [-180, 180]
func (d Direction) Delta(od Direction) int {
	diff := int(d - od)

	for diff <= -180 {
		diff += 360
	}
	for diff > 180 {
		diff -= 360
	}

	if diff < 0 {
		return -diff
	}
	return diff
}
func (d Direction) Vector() (int, int) {
	var (
		dx, dy int
	)
	switch d {
	case North:
		dx, dy = 0, -1
	case NorthEast:
		dx, dy = 1, -1
	case East:
		dx, dy = 1, 0
	case SouthEast:
		dx, dy = 1, 1
	case South:
		dx, dy = 0, 1
	case SouthWest:
		dx, dy = -1, 1
	case West:
		dx, dy = -1, 0
	case NorthWest:
		dx, dy = -1, -1
	}
	return dx, dy
}

type resCountMap map[resource.Type]uint

type Biome struct {
	Name      string
	Symbol    rune
	Color     Color
	Resources resCountMap
}

var (
	Forest = Biome{"Forest", '🌲', Color{0, 250, 0}, resCountMap{resource.Food: 15, resource.Wood: 90}}
	Water  = Biome{"Water", '🌊', Color{0, 0, 250}, resCountMap{resource.Water: 100}}
	Steppe = Biome{"Steppe", '🍀', Color{150, 255, 150}, resCountMap{resource.Food: 2}}
	Desert = Biome{"Desert", '🟨', Color{255, 255, 0}, resCountMap{}}
	Snow   = Biome{"Snow", '⬜', Color{255, 255, 255}, resCountMap{resource.Water: 10}}
	Mount  = Biome{"Mountain", '🗻', Color{190, 190, 190}, resCountMap{resource.Food: 5, resource.Wood: 50}}
)

type World struct {
	img         image.Image
	imgWidth    int
	imgHeight   int
	WorldWidth  int
	WorldHeight int
	Wind        Direction
}

func NewFromBytes(data []byte, worldW, worldH int) (*World, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	b := img.Bounds()

	return &World{
		img:         img,
		imgWidth:    b.Dx(),
		imgHeight:   b.Dy(),
		WorldWidth:  worldW,
		WorldHeight: worldH,
	}, nil
}

func New(path string, worldW, worldH int) (*World, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	b := img.Bounds()

	return &World{
		img:         img,
		imgWidth:    b.Dx(),
		imgHeight:   b.Dy(),
		WorldWidth:  worldW,
		WorldHeight: worldH,
	}, nil
}

func colorMatch(c Color, r, g, b uint8) bool {
	const tol = 40
	return abs(int(c.R)-int(r)) < tol &&
		abs(int(c.G)-int(g)) < tol &&
		abs(int(c.B)-int(b)) < tol
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (w *World) GetBiome(x, y int) Biome {
	// wrap (кругла земля)
	if x < 0 {
		x = (x%w.WorldWidth + w.WorldWidth) % w.WorldWidth
	}
	if y < 0 {
		y = (y%w.WorldHeight + w.WorldHeight) % w.WorldHeight
	}
	if x >= w.WorldWidth {
		x = x % w.WorldWidth
	}
	if y >= w.WorldHeight {
		y = y % w.WorldHeight
	}

	imgX := x * w.imgWidth / w.WorldWidth
	imgY := y * w.imgHeight / w.WorldHeight

	r, g, b, _ := w.img.At(imgX, imgY).RGBA()

	r8 := uint8(r >> 8)
	g8 := uint8(g >> 8)
	b8 := uint8(b >> 8)

	switch {
	case colorMatch(Water.Color, r8, g8, b8):
		return Water
	case colorMatch(Forest.Color, r8, g8, b8):
		return Forest
	case colorMatch(Steppe.Color, r8, g8, b8):
		return Steppe
	case colorMatch(Desert.Color, r8, g8, b8):
		return Desert
	case colorMatch(Mount.Color, r8, g8, b8):
		return Mount
	case colorMatch(Snow.Color, r8, g8, b8):
		return Snow
	default:
		return Steppe
	}
}

func (d Direction) String() string {
	m := map[Direction]string{
		North:     "north",
		NorthEast: "north-east",
		East:      "east",
		SouthEast: "south-east",
		South:     "south",
		SouthWest: "south-west",
		West:      "west",
		NorthWest: "north-west",
	}
	s := m[d]
	return s
}

func (w *World) Live(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			way := rand.Intn(3) - 1
			w.Wind += Direction(45 * way)
			if w.Wind >= 360 {
				w.Wind -= 360
			}
			if w.Wind < 0 {
				w.Wind += 360
			}
		}
	}
}
