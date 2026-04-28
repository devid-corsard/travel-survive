package main

import (
	"bufio"
	"context"
	"fmt"
	_ "image/png"
	"os"
	"time"
	"travel_survival_game/game"
	"travel_survival_game/player"
	"travel_survival_game/world"
)

func main() {
	w, err := world.New("map.png", 50, 50)
	if err != nil {
		panic(err)
	}

	p, err := game.Restore()
	if err != nil {
		fmt.Printf("fail to restore game, creating new: %v", err)
		p = player.New(w.WorldWidth/2, w.WorldHeight/2)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Commands: look, move <north|south|east|west>, exit")

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-time.Tick(time.Duration(time.Minute * 1)):
				p.Live(1)
			case <-ctx.Done():
				fmt.Println("Exit life loop")
				return
			}
		}
	}()

	for {
		fmt.Print("\n> ")
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("failed to read command: %v", err.Error())
		}
		exit := game.RunComand(w, p, command)
		if exit {
			cancel()
			break
		}
		if p.HP <= 0 {
			fmt.Printf("Game over, you died\n")
			cancel()
			break
		}
	}
	<-ctx.Done()
	if err = game.Save(p); err != nil {
		fmt.Printf("failed to save game: %v", err)
	}
}
