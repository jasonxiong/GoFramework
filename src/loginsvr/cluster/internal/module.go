package internal

import (
	"framework/cluster"
	"loginsvr/game"
	"loginsvr/svrconf"
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
