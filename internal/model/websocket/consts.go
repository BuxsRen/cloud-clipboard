package websocket

// Message 消息
type Message struct {
	Code   int         `json:"code"`   // 状态码
	Msg    string      `json:"msg"`    // 信息
	Action string      `json:"action"` // 动作
	Time   int64       `json:"time"`   // 消息时间
	Data   interface{} `json:"data"`   // 消息体
}
