package world

import (
	"image"
	"os"
	"travel_survival_game/resource"
)

type Color struct {
	R, G, B uint8
}

type resCountMap map[resource.Type]uint

type Biome struct {
	Name      string
	Symbol    rune
	Color     Color
	Resources resCountMap
}

var (
	Forest  = Biome{"Forest", 't', Color{0, 250, 0}, resCountMap{resource.Food: 5, resource.Wood: 100}}
	Water   = Biome{"Water", '~', Color{0, 0, 250}, resCountMap{resource.Water: 100}}
	Steppe  = Biome{"Steppe", '"', Color{150, 255, 150}, resCountMap{resource.Food: 2}}
	Desert  = Biome{"Desert", '.', Color{255, 255, 0}, resCountMap{}}
	Mount   = Biome{"Mountain", '^', Color{190, 190, 190}, resCountMap{resource.Food: 5, resource.Wood: 50}}
	Unknown = Biome{"Unknown", '?', Color{0, 0, 0}, resCountMap{}}
)

type World struct {
	img         image.Image
	imgWidth    int
	imgHeight   int
	WorldWidth  int
	WorldHeight int
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
	default:
		return Unknown
	}
}
