package websocket

import (
	ws "cloud-clipboard/internal/model/websocket"
	"context"
	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/net/ghttp"
	"sync"
	"time"
)

var socket *WebSocket

type WebSocket struct {
	iSocket    ISocket                      // 业务接口
	list       *gmap.Map                    // map[string]*Client 在线客户端 {cid: *Client}
	room       *gmap.Map                    // map[string]map[string]*Client 房间 {roomId: {Cid: *Client}}
	routeGroup map[string]*RouteGroup       // 路由组
	route      map[string]*Route            // 事件路由  {"action":controller}
	group      map[string]map[string]*Route // 事件路由组 {"identity": {"action":controller}}
}

// ISocket 必须实现下面方法
type ISocket interface {
	Error(ctx context.Context, client *Client, e any) //  控制器panic错误处理
	Route()                                           // 事件路由
	Login(ctx context.Context, client *Client)        // 客户端登录回调事件
	LogOut(ctx context.Context, client *Client)       // 客户端退出回调事件，退出后通道关闭，所有事件将不可用
}

// RouteGroup 路由组
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

type (
	Middleware func(ctx context.Context, client *Client, msg *ws.Message) (context.Context, bool) // 中间件
	Controller func(ctx context.Context, client *Client, msg *ws.Message)                         // 控制器
)

// Client 客户端
type Client struct {
	cId       string             // 客户端唯一ID
	socket    *ghttp.WebSocket   // ws 链接
	httpCtx   context.Context    // http上下文
	ctx       context.Context    // 客户端的Ctx，客户端退出时会通知到Ctx.Done()，可用于取消当前客户端的异步任务
	cancel    context.CancelFunc // 取消Ctx方法
	time      int64              // WebSocket 接入时间
	msg       chan []byte        // 消息通道
	msgPing   chan []byte        // PING消息通道
	status    bool               // 客户端状态
	ip        string             // 客户端IP
	userAgent string             // 浏览器标识
	roomId    string             // 所在房间
	tk        *time.Ticker       // 定时ping
	once      sync.Once          // 安全锁
	mutex     sync.Mutex         // 用户锁
}
