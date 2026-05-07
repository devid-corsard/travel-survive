package world

import (
	"bytes"
	"context"
	"image"
	"log"
	"math/rand"
	"time"
)

type Color struct {
	R, G, B uint8
}

type World struct {
	// img         image.Image
	// imgWidth    int
	// imgHeight   int
	WorldWidth  int
	WorldHeight int
	Wind        Direction
	biomMap     [][]*Biome
}

func NewFromBytes(data []byte, worldW, worldH int) (*World, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	biomMap := make([][]*Biome, 0, worldW)
	for x := range worldW {
		row := make([]*Biome, 0, worldH)
		for y := range worldH {
			bm := getBiome(img, worldW, worldH, x, y)
			row = append(row, bm)
		}
		biomMap = append(biomMap, row)
	}
	log.Printf("world W: %v, world H: %v, biom map len: %v,biom map 0 row len: %v", worldW, worldH, len(biomMap), len(biomMap[0]))

	return &World{
		WorldWidth:  worldW,
		WorldHeight: worldH,
		biomMap:     biomMap,
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

func (w *World) GetBiome(x, y int) *Biome {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("failed get biom at x: %v, y: %v, len map: %v,len 0 el: %v\n", x, y, len(w.biomMap), len(w.biomMap[0]))
			panic(r)
		}
	}()
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
	bm := w.biomMap[x][y]
	return bm

	// imgX := x * w.imgWidth / w.WorldWidth
	// imgY := y * w.imgHeight / w.WorldHeight

	// r, g, b, _ := w.img.At(imgX, imgY).RGBA()

	// r8 := uint8(r >> 8)
	// g8 := uint8(g >> 8)
	// b8 := uint8(b >> 8)

	// switch {
	// case colorMatch(Water.Color, r8, g8, b8):
	// 	return Water
	// case colorMatch(Forest.Color, r8, g8, b8):
	// 	return Forest
	// case colorMatch(Steppe.Color, r8, g8, b8):
	// 	return Steppe
	// case colorMatch(Desert.Color, r8, g8, b8):
	// 	return Desert
	// case colorMatch(Mount.Color, r8, g8, b8):
	// 	return Mount
	// case colorMatch(Snow.Color, r8, g8, b8):
	// 	return Snow
	// default:
	// 	return Steppe
	// }
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

func getBiome(img image.Image, worldWidth, worldHeight, x, y int) *Biome {
	// wrap (кругла земля)
	if x < 0 {
		x = (x%worldWidth + worldWidth) % worldWidth
	}
	if y < 0 {
		y = (y%worldHeight + worldHeight) % worldHeight
	}
	if x >= worldWidth {
		x = x % worldWidth
	}
	if y >= worldHeight {
		y = y % worldHeight
	}
	ib := img.Bounds()

	imgX := x * ib.Dx() / worldWidth
	imgY := y * ib.Dy() / worldHeight

	r, g, b, _ := img.At(imgX, imgY).RGBA()

	r8 := uint8(r >> 8)
	g8 := uint8(g >> 8)
	b8 := uint8(b >> 8)

	switch {
	case colorMatch(Water.Color, r8, g8, b8):
		return &Water
	case colorMatch(Forest.Color, r8, g8, b8):
		return &Forest
	case colorMatch(Steppe.Color, r8, g8, b8):
		return &Steppe
	case colorMatch(Desert.Color, r8, g8, b8):
		return &Desert
	case colorMatch(Mount.Color, r8, g8, b8):
		return &Mount
	case colorMatch(Snow.Color, r8, g8, b8):
		return &Snow
	default:
		return &Steppe
	}
}
