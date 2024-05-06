package socket

import (
	wsc "cloud-clipboard/internal/consts/websocket"
	ws "cloud-clipboard/internal/model/websocket"
	"cloud-clipboard/utility/websocket"
	"context"
	"github.com/gogf/gf/v2/util/gconv"
)

// Login 客户端接入
func (s *Server) Login(ctx context.Context, client *websocket.Client) {
	room, err := websocket.GetSocketServer().GetRoom(client.GetRoomId())
	if err != nil {
		return
	}
	_ = websocket.GetSocketServer().SendMessageToRoom(ctx, client.GetRoomId(), &ws.Message{Action: wsc.RoomCount, Data: room.Size()}, "")
}

// LogOut 客户端断开
func (s *Server) LogOut(ctx context.Context, client *websocket.Client) {
	room, err := websocket.GetSocketServer().GetRoom(client.GetRoomId())
	if err != nil {
		return
	}
	_ = websocket.GetSocketServer().SendMessageToRoom(ctx, client.GetRoomId(), &ws.Message{Action: wsc.RoomCount, Data: room.Size()}, "")
}

// Error 控制器panic错误处理
func (s *Server) Error(ctx context.Context, client *websocket.Client, e any) {
	client.Send(ctx, &ws.Message{Action: wsc.Error, Code: wsc.CodeError, Msg: gconv.String(e)})
}
