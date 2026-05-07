package main

import (
	"context"
	_ "embed"
	_ "image/png"
	"log"
	"os"
	"time"
	"travel_survival_game/game"
	"travel_survival_game/player"
	"travel_survival_game/render"
	"travel_survival_game/world"

	"github.com/gdamore/tcell/v3"
)

//go:embed map.png
var mapFile []byte

func main() {
	logFile, err := os.OpenFile("game.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	log.SetOutput(logFile)
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC: %v", r)
			panic(r)
		}
	}()

	var info string
	w, err := world.NewFromBytes(mapFile, 1000, 500)
	if err != nil {
		panic(err)
	}

	p, err := game.Restore()
	if err != nil {
		log.Printf("fail to restore game, creating new: %v\n", err)
		p = player.New(w.WorldWidth/2, w.WorldHeight/2)
	}
	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defer s.Fini()
	ren := render.NewTcell(25, s)

	ctx, cancel := context.WithCancel(context.Background())

	renderCh := make(render.Chan)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case inf := <-renderCh:
				ren.Display(w, p, inf)
			}
		}
	}()
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	go func(info *string) {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				renderCh <- *info
			}
		}
	}(&info)
	go w.Live(ctx)
	for {
		// command, err := reader.ReadString('\n')
		// if err != nil {
		// 	fmt.Printf("failed to read command: %v", err.Error())
		// }
		ev := <-s.EventQ()
		// Process event
		var command string
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				command = "q"
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Key() == tcell.KeyLeft || ev.Str() == "a" {
				command = "m w"
			} else if ev.Key() == tcell.KeyRight || ev.Str() == "d" {
				command = "m e"
			} else if ev.Key() == tcell.KeyUp || ev.Str() == "w" {
				command = "m n"
			} else if ev.Key() == tcell.KeyDown || ev.Str() == "s" {
				command = "m s"
			}
		case *tcell.EventMouse:
			// x, y := ev.Position()

			// switch ev.Buttons() {
			// case tcell.Button1, tcell.Button2:
			// 	if ox < 0 {
			// 		ox, oy = x, y // record location when click started
			// 	}

			// case tcell.ButtonNone:
			// 	if ox >= 0 {
			// 		label := fmt.Sprintf("%d,%d to %d,%d", ox, oy, x, y)
			// 		drawBox(s, ox, oy, x, y, boxStyle, label)
			// 		ox, oy = -1, -1
			// 	}
			// }
		}

		exit, inf := game.RunComand(w, p, command, renderCh)
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
		log.Printf("failed to save game: %v\n", err)
	}
}
