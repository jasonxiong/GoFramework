package main

//lobbysvr 为内网服务器，所以只有Cluster，没有Gate

import (
	"lobbysvr/cluster"
	"lobbysvr/game"

	"framework"
)

func main() {
	framework.Run(
		game.Module,
		cluster.Module,
	)
}
