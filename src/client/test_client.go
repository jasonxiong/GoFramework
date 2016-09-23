package main

import (
	"encoding/binary"
	"fmt"
	//"framework/cluster"

	"net"
	"proto/gameproto"

	"github.com/golang/protobuf/proto"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:3563")
	if err != nil {
		panic(err)
	}

	fmt.Println("Success to connect to server.")

	//发送loginreq消息
	reqLogin, err := proto.Marshal(&gameproto.LoginSvr_LoginReq{
		Uin:           proto.Uint32(10000),
		StrSessionKey: proto.String("aaaaabbbbbb"),
	})

	reqData, err := proto.Marshal(&gameproto.ProtoMsg{
		Msgid:   gameproto.MsgID_LOGINSVR_LOGIN_REQ.Enum(),
		Uin:     proto.Uint32(10000),
		Msgdata: reqLogin,
	})

	//2bytes len + data
	m := make([]byte, 2+len(reqData))

	//写入2bytes长度
	binary.BigEndian.PutUint16(m, uint16(len(reqData)))

	//写入实际数据
	copy(m[2:], reqData)

	//发送消息
	conn.Write(m)

	fmt.Println("Success to send data, len ", len(m))

	//接收消息
	recvMsg := make([]byte, 1024)
	len, err := conn.Read(recvMsg)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success to read data from server ", len)

	//接收到消息格式也为2bytes len + data
	msgLen := binary.BigEndian.Uint16(recvMsg[:2])

	fmt.Printf("recv msg %v\n", recvMsg[2:])

	respMsg := &gameproto.ProtoMsg{}
	proto.Unmarshal(recvMsg[2:msgLen], respMsg)

	respLogin := &gameproto.LoginSvr_LoginResp{}
	proto.Unmarshal(respMsg.GetMsgdata(), respLogin)

	//cluster.Agent.Run()
	fmt.Printf("recv msg len: %d, msg uin: %u, result %d\n", msgLen, respMsg.GetUin(), respLogin.GetIResult())

	conn.Close()
}
