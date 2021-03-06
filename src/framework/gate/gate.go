package gate

import (
	"time"

	"framework/chanrpc"
	"framework/log"
	"framework/network"
)

//Gate中不再解析消息，所有网络消息都发送到gate中进行处理

type Gate struct {
	MaxConnNum      int
	PendingWriteNum int
	MaxMsgLen       uint32
	AgentChanRPC    *chanrpc.Server

	// websocket
	WSAddr      string
	HTTPTimeout time.Duration

	// tcp
	TCPAddr      string
	LenMsgLen    int
	LittleEndian bool
}

func (gate *Gate) Run(closeSig chan bool) {
	var wsServer *network.WSServer
	if gate.WSAddr != "" {
		wsServer = new(network.WSServer)
		wsServer.Addr = gate.WSAddr
		wsServer.MaxConnNum = gate.MaxConnNum
		wsServer.PendingWriteNum = gate.PendingWriteNum
		wsServer.MaxMsgLen = gate.MaxMsgLen
		wsServer.HTTPTimeout = gate.HTTPTimeout
		wsServer.NewAgent = func(conn *network.WSConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go("NewGateAgent", a)
			}
			return a
		}
	}

	var tcpServer *network.TCPServer
	if gate.TCPAddr != "" {
		tcpServer = new(network.TCPServer)
		tcpServer.Addr = gate.TCPAddr
		tcpServer.MaxConnNum = gate.MaxConnNum
		tcpServer.PendingWriteNum = gate.PendingWriteNum
		tcpServer.LenMsgLen = gate.LenMsgLen
		tcpServer.MaxMsgLen = gate.MaxMsgLen
		tcpServer.LittleEndian = gate.LittleEndian
		tcpServer.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := &agent{conn: conn, gate: gate}
			if gate.AgentChanRPC != nil {
				gate.AgentChanRPC.Go("NewGateAgent", a)
			}
			return a
		}
	}

	if wsServer != nil {
		wsServer.Start()
	}
	if tcpServer != nil {
		tcpServer.Start()
	}
	<-closeSig
	if wsServer != nil {
		wsServer.Close()
	}
	if tcpServer != nil {
		tcpServer.Close()
	}
}

func (gate *Gate) OnDestroy() {}

type agent struct {
	conn     network.Conn
	gate     *Gate
	userData interface{}
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message: %v", err)
			break
		}

		//直接路由给game模块处理
		if a.gate.AgentChanRPC != nil {
			a.gate.AgentChanRPC.Go("OnGateMsg", data, a)
		}

		//Gate模块只路由，不解析消息
		/*
			if a.gate.Processor != nil {
				msg, err := a.gate.Processor.Unmarshal(data)
				if err != nil {
					log.Debug("unmarshal message error: %v", err)
					break
				}
				err = a.gate.Processor.Route(msg, a)
				if err != nil {
					log.Debug("route message error: %v", err)
					break
				}
			}
		*/
	}
}

func (a *agent) OnClose() {
	if a.gate.AgentChanRPC != nil {
		err := a.gate.AgentChanRPC.Call0("NewGateAgent", a)
		if err != nil {
			log.Error("chanrpc error: %v", err)
		}
	}
}

func (a *agent) WriteMsg(data interface{}) {
	err := a.conn.WriteMsg([][]byte{data.([]byte)}...)
	if err != nil {
		log.Error("write message to gate error: %v", err)
	}

	//Gate现在只负责发送消息，而不负责消息的打解包
	/*
		if a.gate.Processor != nil {
			data, err := a.gate.Processor.Marshal(msg)
			if err != nil {
				log.Error("marshal message %v error: %v", reflect.TypeOf(msg), err)
				return
			}
			err = a.conn.WriteMsg(data...)
			if err != nil {
				log.Error("write message %v error: %v", reflect.TypeOf(msg), err)
			}
		}
	*/
}

func (a *agent) Close() {
	a.conn.Close()
}

func (a *agent) Destroy() {
	a.conn.Destroy()
}

func (a *agent) UserData() interface{} {
	return a.userData
}

func (a *agent) SetUserData(data interface{}) {
	a.userData = data
}
