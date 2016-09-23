package cluster

type Agent interface {
	WriteMsg(data interface{})
	Close()
	Destroy()
	UserData() interface{}
	SetUserData(data interface{})
}
