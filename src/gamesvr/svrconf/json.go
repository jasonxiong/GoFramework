package svrconf

import (
	"encoding/json"
	"framework/conf"
	"framework/log"
	"io/ioutil"
)

//Server 服务器配置
var Server struct {
	WSAddr        string
	TCPAddr       string
	MaxConnNum    int
	SvrBaseConfig *conf.BaseConf
}

func init() {
	data, err := ioutil.ReadFile("conf/gamesvr.json")
	if err != nil {
		log.Fatal("%v", err)
	}

	Server.SvrBaseConfig = &conf.SvrBase
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}
}
