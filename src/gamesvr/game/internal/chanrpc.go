package internal

import (
	"framework/cluster"
	"framework/conf"
	"framework/gate"
	"framework/log"
	"proto/gameproto"

	"github.com/golang/protobuf/proto"
)

//ClusterAgentInfo 集群Agent信息
type ClusterAgentInfo struct {
	endpointID int32 //对端ID
}

var (
	clusterAgent = make(map[int32]cluster.Agent)
)

func init() {

	//创建和关闭Gate Agent
	skeleton.RegisterChanRPC("NewGateAgent", newGateAgent)
	skeleton.RegisterChanRPC("CloseGateAgent", closeGateAgent)

	//Gate消息的处理
	skeleton.RegisterChanRPC("OnGateMsg", handleGateMsg)

	//创建和关闭Cluster Agent
	skeleton.RegisterChanRPC("NewClusterAgent", newClusterAgent)
	skeleton.RegisterChanRPC("CloseClusterAgent", closeClusterAgent)

	//Cluster消息处理
	skeleton.RegisterChanRPC("OnClusterMsg", handleClusterMsg)
}

func newGateAgent(args []interface{}) {
	//todo jasonxiong
	a := args[0].(gate.Agent)
	_ = a
}

func closeGateAgent(args []interface{}) {
	//todo jasonxiong
	a := args[0].(gate.Agent)
	_ = a
}

//处理Gate的消息
func handleGateMsg(args []interface{}) {
	data := args[0].([]byte)
	a := args[1].(gate.Agent)

	//data解包
	msg := &gameproto.ProtoMsg{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		log.Error("Failed to parse ProtoMsg from gate, err %v\n", err)
		return
	}

	//处理实际的消息
	if handler, ok := gameMsgHandlers[msg.GetMsgid()]; !ok {
		//不存在
		log.Error("Failed to handler game msg %u, handler not registered.\n", msg.GetMsgid())
	} else {
		handler([]interface{}{a, msg.GetMsgdata()})
	}

	return
}

func newClusterAgent(args []interface{}) {
	a := args[0].(cluster.Agent)
	endpointID := args[1].(int32)

	agentInfo := new(ClusterAgentInfo)
	a.SetUserData(agentInfo)

	if endpointID != 0 {
		if _, ok := clusterAgent[endpointID]; ok {
			//对端agent已经注册过
			log.Error("failed to update endpointID %d, already registered.\n", endpointID)
			return
		}

		//对端为Cluster服务端
		agentInfo.endpointID = endpointID
		clusterAgent[endpointID] = a

		log.Debug("new cluster agent, endpointID %d\n", endpointID)

		//通知对端本端的ID
		notifyClusterInfo, err := proto.Marshal(&gameproto.Cluster_UpdateInfoNotify{LocalendID: proto.Int32(conf.SvrBase.ServerID)})
		if err != nil {
			log.Error("Failed to encode notify cluster info msg, err %v\n", err)
			return
		}

		notifyMsg := &gameproto.ProtoMsg{
			Msgid:   gameproto.MsgID_CLUSTER_UPDATEINFO_NOTIFY.Enum(),
			Msgdata: notifyClusterInfo,
		}

		//序列化
		notifyData, err := proto.Marshal(notifyMsg)
		if err != nil {
			log.Error("Failed to marshal notify msg, err %v\n", err)
			return
		}

		a.WriteMsg(notifyData)
	}
}

func closeClusterAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	endpointID := a.UserData().(*ClusterAgentInfo).endpointID

	//先从map中删除
	delete(clusterAgent, endpointID)

	//关闭agent
	a.SetUserData(nil)
	a.Close()

	log.Debug("Success to close cluster agent, endpointID %d\n", endpointID)
}

//处理Gate的消息
func handleClusterMsg(args []interface{}) {
	data := args[0].([]byte)
	a := args[1].(gate.Agent)

	//data解包
	msg := &gameproto.ProtoMsg{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		log.Error("Failed to parse ProtoMsg from gate, err %v\n", err)
		return
	}

	//处理实际的消息
	if handler, ok := gameMsgHandlers[msg.GetMsgid()]; !ok {
		//不存在
		log.Error("Failed to handler game msg %u, handler not registered.\n", msg.GetMsgid())
	} else {
		handler([]interface{}{a, msg.GetMsgdata()})
	}

	return
}
