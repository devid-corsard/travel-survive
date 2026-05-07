package world

import "travel_survival_game/resource"

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
