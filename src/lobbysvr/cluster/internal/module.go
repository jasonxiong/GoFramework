package internal

import (
	"framework/cluster"
	"lobbysvr/game"
	"lobbysvr/svrconf"
)

type Module struct {
	*cluster.Cluster
}

func (m *Module) OnInit() {
	m.Cluster = &cluster.Cluster{
		MaxMsgLen:    svrconf.MaxMsgLen,
		AgentChanRPC: game.ChanRPC,
	}
}
