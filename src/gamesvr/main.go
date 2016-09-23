package main

import (
	"gamesvr/cluster"
	"gamesvr/game"
	"gamesvr/gate"

	"framework"
)

func main() {
	framework.Run(
		game.Module,
		gate.Module,
		cluster.Module,
	)
}
