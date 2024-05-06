package websocket

const (
	CodeError       = 500  // 错误
	CodeLimit       = 429  // 消息发送过于频繁
	CodeNotFound    = 404  // 找不到
	CodeNoAuth      = 401  // 未鉴权
	CodeSuccess     = 0    // 成功
	CodeFail        = -1   // 失败
	CodeNull        = -2   // 空数据
	CodeTokenExpire = -99  // 身份过期
	CodeClose       = -410 // 客户端被迫下线
)
