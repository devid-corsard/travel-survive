package world

import "fmt"

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
