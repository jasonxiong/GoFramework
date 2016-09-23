捕鱼侠游戏服务器

1.捕鱼侠游戏服务器包含： 登录服务器，游戏服务器和中心服务器；
2.基于捕鱼游戏数据特性，采用mongodb代替mysql；
3.中心服务器上的数据可能需要采用缓存，待定；

集群通信组件
1.集群服务器之间通过socket连接进行通信；
2.每个服务器都有自己唯一的通信组件ID，生成规则为: ServerType*1000+ServerID;
    ServerType为服务器类型，区间0-99；
    ServerID为服务器在同类服务中的index，区间为1-999;
3.目前初步使用的ServerType分配为:
    LOGIN_SERVER    10  //登录服务器组
    GAME_SERVER     11  //游戏服务器组
    LOBBY_SERVER    12  //中心服务器组
