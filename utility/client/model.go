package client

import (
	ws "CloudContent/internal/model/websocket"
	"context"
	"github.com/gorilla/websocket"
	"sync"
)

var socket *Socket

type Socket struct {
	iSocket     SocketInterface
	routeGroup  map[string]*RouteGroup
	route       map[string]Route            // 事件路由  {"action":controller}
	group       map[string]map[string]Route // 事件路由组 {"identity": {"action":controller}}
	Client      *Client
	againStatus bool // 重连状态 true 表示需要重新连接
	mutex       sync.Mutex
}

// SocketInterface 必须实现下面方法
type SocketInterface interface {
	Error(ctx context.Context, e any) //  控制器panic错误处理
	Route()                           // 通用事件路由
	Login(ctx context.Context)        // ws连接成功回调事件
	Logout(ctx context.Context)       // ws断开回调事件，断开后通道关闭，所有事件将不可用
}

type RouteGroup struct {
	group *Group
	next  []*RouteGroup
}

// Route 事件路由 优先级 Route > Group
type Route struct {
	action     string       // 路由名称或动作名称
	fun        Controller   // 对应处理的方法
	middleware []Middleware // 消息中间件
}

// Group 事件路由组 优先级 Group.Route > Group.Group
type Group struct {
	Route      []*Route     // 事件路由  同时使用Route和Group时，优先访问Route
	Middleware []Middleware // 路由组消息中间件
	Group      []*Group     // 子路由组
}

type Middleware func(ctx context.Context, msg *ws.Message) bool // 中间件
type Controller func(ctx context.Context, msg *ws.Message)      // 控制器

type Client struct {
	ctx     context.Context
	conn    *websocket.Conn
	msgChan chan []byte // 消息通道
	Status  bool
	cancel  context.CancelFunc // 取消Ctx方法
	once    sync.Once          // 安全锁
	mutex   sync.Mutex         // 用户锁
}
