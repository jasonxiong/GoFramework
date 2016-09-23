package cluster

import (
	"framework/chanrpc"
	"framework/conf"
	"framework/log"
	"framework/network"
	"math"
	"time"
)

//Client 集群通信客户端
type Client struct {
	id        int32             //通信对端的唯一标识
	tcpclient network.TCPClient //TCP Client结构
}

//Cluster 集群通信管理类
type Cluster struct {
	MaxMsgLen    uint32
	AgentChanRPC *chanrpc.Server

	server  *network.TCPServer
	clients []*Client
}

//Run 运行函数
func (cluster *Cluster) Run(closeSig chan bool) {
	if conf.SvrBase.ListenAddr != "" {
		cluster.server = new(network.TCPServer)
		cluster.server.Addr = conf.SvrBase.ListenAddr
		cluster.server.MaxConnNum = int(math.MaxInt32)
		cluster.server.PendingWriteNum = conf.SvrBase.PendingWriteNum
		cluster.server.LenMsgLen = 2
		cluster.server.MaxMsgLen = cluster.MaxMsgLen
		cluster.server.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := new(agent)
			a.conn = conn
			a.cluster = cluster

			if a.cluster.AgentChanRPC != nil {
				//对端ID不确定，传0
				var id int32
				a.cluster.AgentChanRPC.Go("NewClusterAgent", a, id)
			}
			return a
		}

		cluster.server.Start()
	}

	for _, addr := range conf.SvrBase.ConnAddrs {
		client := new(Client)
		client.id = addr.ID
		client.tcpclient.Addr = addr.ConnAddr
		client.tcpclient.ConnNum = 1
		client.tcpclient.ConnectInterval = 3 * time.Second
		client.tcpclient.PendingWriteNum = conf.SvrBase.PendingWriteNum
		client.tcpclient.LenMsgLen = 2
		client.tcpclient.MaxMsgLen = cluster.MaxMsgLen
		client.tcpclient.NewAgent = func(conn *network.TCPConn) network.Agent {
			a := new(agent)
			a.conn = conn
			a.cluster = cluster

			if a.cluster.AgentChanRPC != nil {
				//对端ID确定
				a.cluster.AgentChanRPC.Go("NewClusterAgent", a, client.id)
			}
			return a
		}

		client.tcpclient.Start()
		cluster.clients = append(cluster.clients, client)
	}

	<-closeSig

	//关闭Cluster服务
	if cluster.server != nil {
		cluster.server.Close()
	}

	//关闭Cluster客户端
	for _, client := range cluster.clients {
		client.tcpclient.Close()
	}
}

//OnDestroy 销毁
func (cluster *Cluster) OnDestroy() {}

type agent struct {
	cluster  *Cluster
	conn     *network.TCPConn
	userData interface{}
}

func (a *agent) Run() {
	for {
		data, err := a.conn.ReadMsg()
		if err != nil {
			log.Debug("read message failed, err %v\n", err)
			break
		}

		//路由给其他模块处理
		if a.cluster.AgentChanRPC != nil {
			a.cluster.AgentChanRPC.Go("OnClusterMsg", data, a)
		}
	}
}

func (a *agent) OnClose() {
	if a.cluster.AgentChanRPC != nil {
		err := a.cluster.AgentChanRPC.Call0("CloseClusterAgent", a)
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
