package main

import(
	"github.com/awebbdev/gameswithgo/rpg/game"
	"github.com/awebbdev/gameswithgo/rpg/ui2d"
	"runtime"
)

func main() {
	game := game.NewGame(1, "game/maps/level1.map")

	for i := 0; i < 1; i++ {
		go func(i int) {
			runtime.LockOSThread()
			ui := ui2d.NewUI(game.InputChan, game.LevelChans[i])
			ui.Run()
		}(i)
	}

	game.Run()
}