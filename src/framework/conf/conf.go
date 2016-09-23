package conf

//AddrInfo 集群通信连接对端信息
type AddrInfo struct {
	ID       int32  //通信对端的唯一标识
	ConnAddr string //通信对端的地址信息
}

//BaseConf 基础配置
type BaseConf struct {
	// 日志
	LogLevel string
	LogPath  string

	//服务ID
	ServerID int32

	//控制台
	ConsolePort   int
	ConsolePrompt string
	ProfilePath   string

	//Cluster
	ListenAddr      string
	ConnAddrs       []AddrInfo
	PendingWriteNum int
}

var (
	LenStackBuf = 4096
	SvrBase     BaseConf
)
