package main

import (
	"bufio"
	"context"
	_ "embed"
	"fmt"
	_ "image/png"
	"os"
	"travel_survival_game/game"
	"travel_survival_game/player"
	"travel_survival_game/world"
)

//go:embed map.png
var mapFile []byte

func main() {
	w, err := world.NewFromBytes(mapFile, 50, 50)
	if err != nil {
		panic(err)
	}

	p, err := game.Restore()
	if err != nil {
		fmt.Printf("fail to restore game, creating new: %v\n", err)
		p = player.New(w.WorldWidth/2, w.WorldHeight/2)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Commands: look, move <north|south|east|west>, exit")

	ctx, cancel := context.WithCancel(context.Background())
	go w.Live(ctx)
	for {
		fmt.Print("\n> ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("failed to read command: %v", err.Error())
		}
		exit := game.RunComand(w, p, command)
		if exit {
			break
		}
		if p.HP <= 0 {
			fmt.Printf("Game over, you died\n")
			break
		}
	}
	cancel()

	if err = game.Save(p); err != nil {
		fmt.Printf("failed to save game: %v\n", err)
	}
}
