package main

import(
	"github.com/awebbdev/gameswithgo/rpg/game"
	"github.com/awebbdev/gameswithgo/rpg/ui2d"
)

func main() {
	ui := &ui2d.UI2d{}
	game.Run(ui)
}