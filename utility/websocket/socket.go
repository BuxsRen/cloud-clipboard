package websocket

import (
	"CloudContent/internal/consts/websocket"
	"context"
	"errors"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/guid"
	"sync"
	"time"
)

func Init(si ISocket) {
	if socket != nil {
		return
	}

	socket = &WebSocket{
		iSocket:    si,
		list:       gmap.New(true),
		room:       gmap.New(true),
		routeGroup: make(map[string]*RouteGroup),
		route:      make(map[string]Route),
		group:      make(map[string]map[string]Route),
	}

	socket.iSocket.Route()

	// 绑定路由组路由事件
	for k, v := range socket.routeGroup {
		route := make(map[string]Route)

		socket.handleRouteGroup(&route, v)

		socket.group[k] = route
	}

	socket.PrintRoute()
}

// HandleClient 初始化并处理一个ws链接
func HandleClient(r *ghttp.Request) {
	ws, err := r.WebSocket()
	if err != nil {
		r.Response.Writeln("Protocol Upgrade Failed")
		r.Exit()
	}

	roomId := r.Get("roomId")
	if roomId.IsEmpty() {
		_ = ws.Close()
		return
	}

	client := &Client{
		cId:       guid.S([]byte(r.GetClientIp())),
		socket:    ws,
		httpCtx:   r.GetCtx(),
		msg:       make(chan []byte, 10), // 设置消息缓存池
		msgPing:   make(chan []byte),
		status:    true,
		roomId:    roomId.String(),
		ip:        r.GetClientIp(),
		userAgent: r.UserAgent(),
		time:      gtime.Now().Unix(),
		tk:        time.NewTicker(time.Minute),
		once:      sync.Once{},
		mutex:     sync.Mutex{},
	}

	client.ctx, client.cancel = context.WithCancel(gctx.New())

	socket.connect(client.ctx, client)

}

// GetSocketServer 获取Socket实例
func GetSocketServer() *WebSocket {
	return socket
}

// ForcedClient 下线客户端
func (ws *WebSocket) ForcedClient(ctx context.Context, list *gmap.Map, filter *Client) {
	// 将这个用户下登录的所有app踢下线
	filter.Send(ctx, &Message{Action: websocket.Close, Msg: "您的账号在另一台设备登录了，您被迫下线", Code: websocket.CodeClose})
	go func(c *Client) {
		time.Sleep(time.Second)
		c.Close()
	}(filter)
}

// GetAllRoom 获取所有房间客户端列表
func (ws *WebSocket) GetAllRoom() *gmap.Map {
	return ws.room
}

// GetRoom 获取一个房间内的客户端
func (ws *WebSocket) GetRoom(roomId string) (*gmap.Map, error) {
	room := ws.room.Get(roomId)
	if room == nil {
		return nil, errors.New("没有这个房间")
	}

	return room.(*gmap.Map), nil
}

// GetClient 获取某个房间内的客户端
func (ws *WebSocket) GetClient(roomId, cId string) (*Client, error) {
	room := ws.room.Get(roomId)
	if room == nil {
		return nil, errors.New("没有这个房间")
	}

	client := room.(*gmap.Map).Get(cId)
	if client == nil {
		return nil, errors.New("没有客户端在线")
	}

	return client.(*Client), nil
}

// SendMessageToAll 给所有客户端发送消息
func (ws *WebSocket) SendMessageToAll(ctx context.Context, msg *Message, filter string) {
	ws.list.Iterator(func(roomId interface{}, v interface{}) bool {
		ws.sendToRoom(ctx, roomId.(string), msg, filter)
		return true
	})
}

// SendMessageToRoom 给指定房间内所有客户端发送消息
func (ws *WebSocket) SendMessageToRoom(ctx context.Context, roomId string, msg *Message, filter string) error {
	room := ws.room.Get(roomId)
	if room == nil {
		return errors.New("没有这个房间")
	}

	ws.sendToRoom(ctx, roomId, msg, filter)
	return nil
}

// SendMessageToClientInRoom 给指定房间内指定客户端发送消息
func (ws *WebSocket) SendMessageToClientInRoom(ctx context.Context, roomId, cId string, msg *Message) error {
	room := ws.room.Get(roomId)
	if room == nil {
		return errors.New("没有这个房间")
	}
	client := room.(*gmap.Map).Get(cId)
	if room == nil {
		return errors.New("没有客户端在线")
	}
	client.(*Client).Send(ctx, msg)
	return nil
}
