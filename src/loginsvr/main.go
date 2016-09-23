package main

import (
	"loginsvr/cluster"
	"loginsvr/game"
	"loginsvr/gate"

	"framework"
)

func main() {
	framework.Run(
		game.Module,
		gate.Module,
		cluster.Module,
	)
}
