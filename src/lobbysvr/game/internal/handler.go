package internal

import (
	"framework/cluster"
	"framework/log"
	"proto/gameproto"

	"github.com/golang/protobuf/proto"
)

//GameMsgHandler 消息handler
type GameMsgHandler func([]interface{})

var (
	//Game模块消息处理
	gameMsgHandlers = make(map[gameproto.MsgID]GameMsgHandler)
)

func init() {
	//集群相关消息
	registerHandler(gameproto.MsgID_CLUSTER_UPDATEINFO_NOTIFY, handleClusterInfo)
}

//注册消息handler
func registerHandler(id gameproto.MsgID, gameMsgHandler GameMsgHandler) {
	if _, ok := gameMsgHandlers[id]; ok {
		//已经注册过
		log.Error("Failed to register msg handler %u\n", id)
		return
	}

	gameMsgHandlers[id] = gameMsgHandler

	return
}

//处理集群信息更新
func handleClusterInfo(args []interface{}) {
	a := args[0].(cluster.Agent)
	data := args[1].([]byte)

	//解析消息
	reqMsg := &gameproto.Cluster_UpdateInfoNotify{}
	err := proto.Unmarshal(data, reqMsg)
	if err != nil {
		log.Error("Failed to parse clustinfo msg, err %v\n", err)
		return
	}

	if _, ok := clusterAgent[reqMsg.GetLocalendID()]; ok {
		//agent已经注册
		log.Error("failed to update endpointid %d, already registered.\n", reqMsg.GetLocalendID())
		return
	}

	//更新Cluster对端ID
	clusterAgent[reqMsg.GetLocalendID()] = a
	a.UserData().(*ClusterAgentInfo).endpointID = reqMsg.GetLocalendID()

	log.Debug("update cluster info, endpointid %d\n", reqMsg.GetLocalendID())

	return
}
