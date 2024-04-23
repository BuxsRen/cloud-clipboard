package app

import (
	wsc "CloudContent/internal/consts/websocket"
	"CloudContent/internal/model/content"
	ws "CloudContent/internal/model/websocket"
	"CloudContent/internal/service"
	"CloudContent/utility/websocket"
	"context"
	"github.com/gogf/gf/v2/util/gconv"
)

type App struct{}

// GetContent 获取内容
func (a *App) GetContent(ctx context.Context, client *websocket.Client, msg *ws.Message) {
	cache, err := service.Cache().GetCache(ctx, client.GetRoomId())
	if err != nil {
		return
	}
	client.Send(ctx, &ws.Message{Action: wsc.Content, Data: cache.String()})
}

// SetContent 保存内容
func (a *App) SetContent(ctx context.Context, client *websocket.Client, msg *ws.Message) {
	var req content.Content
	_ = gconv.Scan(msg.Data, &req)

	err := service.Cache().SetCache(ctx, client.GetRoomId(), req.Content)
	if err != nil {
		return
	}

	// 内容广播到房间内
	_ = websocket.GetSocketServer().SendMessageToRoom(ctx, client.GetRoomId(), &ws.Message{Action: wsc.Content, Data: req.Content}, client.GetCId())
}
