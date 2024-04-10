package websocket

const (
	Ping    = "ping"    // 心跳消息
	Error   = "error"   // 错误消息
	Online  = "online"  // 客户端上线消息
	Offline = "offline" // 客户端下线消息
	Close   = "close"   // WS服务被迫下线
)

const (
	RoomCount  = "room_count"
	GetContent = "get_content"
	SetContent = "set_content"
	Content    = "content"
)
