package middleware

import (
	wsc "cloud-clipboard/internal/consts/websocket"
	ws "cloud-clipboard/internal/model/websocket"
	"cloud-clipboard/internal/service"
	"cloud-clipboard/utility/websocket"
	"context"
	"time"
)

// Limit 消息限流
func Limit(ctx context.Context, client *websocket.Client, msg *ws.Message) (context.Context, bool) {
	var (
		key    = "socket_msg_limit_" + client.GetCId()
		limit  = 60              // 最多可发送 N 条消息
		expire = 5 * time.Second // N 秒内
	)

	client.GetMutex().Lock()
	defer client.GetMutex().Unlock()

	v, _ := service.Cache().GetCache(ctx, key)
	if v.IsEmpty() {
		_ = service.Cache().SetxCache(ctx, key, 1, expire)
		return ctx, true
	}

	count := v.Int() + 1
	if count > limit {
		client.Send(ctx, &ws.Message{Action: wsc.Error, Code: wsc.CodeLimit, Msg: "消息发送过于频繁"})
		_ = service.Cache().SetxCache(ctx, key, count, expire)
		return ctx, false
	}

	_, _, _ = service.Cache().Update(ctx, key, count)
	return ctx, true
}
