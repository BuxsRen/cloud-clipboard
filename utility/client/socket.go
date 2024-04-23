package client

import (
	ws "CloudContent/internal/model/websocket"
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"sync"
	"time"
)

func GetSocketClient() *Socket {
	return socket
}

// Routes 添加通用路由
func Routes(action string, controller Controller) Route {
	if _, ok := socket.route[action]; ok {
		panic("socket route duplicate registration: " + action)
	}
	socket.route[action] = Route{
		action:     action,
		middleware: []Middleware{},
		fun:        controller,
	}
	return socket.route[action]
}

// Middleware 通用路由-中间件
func (r *Route) Middleware(middleware Middleware) *Route {
	r.middleware = append(r.middleware, middleware)
	return r
}

// Groups 添加分组路由
func Groups(identity string, fun func(group *RouteGroup)) {
	if _, ok := socket.routeGroup[identity]; ok {
		panic("socket group duplicate registration: " + identity)
	}
	group := &RouteGroup{
		group: &Group{},
	}
	socket.routeGroup[identity] = group
	fun(group)
}

// Group 分组路由-组
func (group *RouteGroup) Group(fun func(group *RouteGroup)) {
	next := &RouteGroup{
		group: &Group{},
	}
	group.next = append(group.next, next)
	fun(next)
}

// Middleware 分组路由-中间件
func (group *RouteGroup) Middleware(middleware Middleware) {
	group.group.Middleware = append(group.group.Middleware, middleware)
}

// Route 分组路由-路由
func (group *RouteGroup) Route(action string, controller Controller) {
	group.group.Route = append(group.group.Route, &Route{action: action, fun: controller})
}

// handleRouteGroup 处理路由组与路由组链
func (s *Socket) handleRouteGroup(route *map[string]Route, group *RouteGroup) {
	socket.handleGroup(route, group.group, []Middleware{})

	if group.next != nil {
		for _, v := range group.next {
			v.group.Middleware = append(group.group.Middleware, v.group.Middleware...)
			s.handleRouteGroup(route, v)
		}
	}
}

// handleGroup 处理路由组
func (s *Socket) handleGroup(route *map[string]Route, group *Group, middleware []Middleware) {
	if group.Middleware != nil { // 合并上层中间件与本层中间件
		middleware = append(middleware, group.Middleware...)
	}

	for _, v2 := range group.Route { // 遍历路由
		if _, ok := (*route)[v2.action]; ok {
			panic("socket route group duplicate registration: " + v2.action)
		}
		v2.middleware = append(middleware, v2.middleware...) // 合并中间件到路由
		(*route)[v2.action] = *v2
	}

	if group.Group != nil {
		for _, v1 := range group.Group { // 遍历路由组
			s.handleGroup(route, v1, middleware) // 传递上层中间件
		}
	}
}

// Again 重连
func (s *Socket) again() {
	tk := time.NewTicker(time.Second)

	for {
		select {
		case <-tk.C:
			if s.againStatus {
				if s.Client != nil && !s.Client.Status {
					s.againStatus = false
					NewClient(gctx.New())
				}
			}
		}
	}
}

// SetAgainStatus 设置是否需要重连
func (s *Socket) SetAgainStatus(status bool) {
	s.againStatus = status
}

// Close 关闭ws客户端
func (s *Socket) Close() {
	if s.Client != nil && s.Client.conn != nil && s.Client.Status {
		s.Client.safeClose()
	}
}

// Send 给网关发送消息
func (s *Socket) Send(ctx context.Context, msg *ws.Message) {
	if msg.Msg == "" {
		msg.Msg = "success"
	}
	if s.Client != nil && s.Client.Status {
		s.Client.msgChan <- s.BuildMsg(msg)
	}
}

func (s *Socket) BuildMsg(msg *ws.Message) []byte {
	msg.Time = gtime.Now().Unix()
	b, _ := json.Marshal(msg)
	return b
}

func (s *Socket) GetMutex() *sync.Mutex {
	return &s.mutex
}
