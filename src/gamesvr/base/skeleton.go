package base

import (
	"framework/chanrpc"
	"framework/module"
	"gamesvr/svrconf"
)

func NewSkeleton() *module.Skeleton {
	skeleton := &module.Skeleton{
		GoLen:              svrconf.GoLen,
		TimerDispatcherLen: svrconf.TimerDispatcherLen,
		AsynCallLen:        svrconf.AsynCallLen,
		ChanRPCServer:      chanrpc.NewServer(svrconf.ChanRPCLen),
	}
	skeleton.Init()
	return skeleton
}
