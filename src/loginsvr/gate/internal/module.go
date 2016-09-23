package internal

import (
	"framework/conf"
	"framework/gate"
	"loginsvr/game"
	"loginsvr/svrconf"
)

type Module struct {
	*gate.Gate
}

func (m *Module) OnInit() {
	m.Gate = &gate.Gate{
		MaxConnNum:      svrconf.Server.MaxConnNum,
		PendingWriteNum: conf.SvrBase.PendingWriteNum,
		MaxMsgLen:       svrconf.MaxMsgLen,
		WSAddr:          svrconf.Server.WSAddr,
		HTTPTimeout:     svrconf.HTTPTimeout,
		TCPAddr:         svrconf.Server.TCPAddr,
		LenMsgLen:       svrconf.LenMsgLen,
		LittleEndian:    svrconf.LittleEndian,
		AgentChanRPC:    game.ChanRPC,
	}
}
