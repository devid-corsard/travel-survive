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
	"travel_survival_game/render"
	"travel_survival_game/world"
)

//go:embed map.png
var mapFile []byte

func main() {
	w, err := world.NewFromBytes(mapFile, 512, 256)
	if err != nil {
		panic(err)
	}

	p, err := game.Restore()
	if err != nil {
		fmt.Printf("fail to restore game, creating new: %v\n", err)
		p = player.New(w.WorldWidth/2, w.WorldHeight/2)
	}
	ren := render.NewText(25)
	var info string

	reader := bufio.NewReader(os.Stdin)

	ctx, cancel := context.WithCancel(context.Background())
	go w.Live(ctx)
	for {
		ren.Display(w, p, info)
		fmt.Print("\n> ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("failed to read command: %v", err.Error())
		}
		exit, inf := game.RunComand(w, p, command, ren)
		info = inf
		if exit {
			break
		}
		if p.HP <= 0 {
			ren.Display(w, p, "Game over, you died")
			break
		}
	}
	cancel()

	if err = game.Save(p); err != nil {
		fmt.Printf("failed to save game: %v\n", err)
	}
}
