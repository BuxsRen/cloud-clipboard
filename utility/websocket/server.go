package websocket

import (
	ws "CloudContent/internal/model/websocket"
	"CloudContent/utility/utils"
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/gmap"
	"reflect"
	"runtime"
	"strings"
)

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

// connect 客户端连接
func (ws *WebSocket) connect(ctx context.Context, client *Client) {
	// 开启读/写协程
	go client.read()
	go client.write()

	// 开启 ping 消息
	go client.ping()

	ws.saveClientToList(client)
	ws.savaClientToRoom(client)

	// 通知到登录事件
	ws.iSocket.Login(ctx, client)
}

// disconnect 客户端断开连接
func (ws *WebSocket) disconnect(ctx context.Context, client *Client) {
	ws.removeClientInList(client)
	ws.removeClientInRoom(client)

	// 通知到退出事件
	ws.iSocket.LogOut(ctx, client)
}

// handleRouteGroup 处理路由组与路由组链
func (ws *WebSocket) handleRouteGroup(route *map[string]Route, group *RouteGroup) {
	socket.handleGroup(route, group.group, []Middleware{})

	if group.next != nil {
		for _, v := range group.next {
			v.group.Middleware = append(group.group.Middleware, v.group.Middleware...)
			ws.handleRouteGroup(route, v)
		}
	}
}

// handleGroup 处理路由组
func (ws *WebSocket) handleGroup(route *map[string]Route, group *Group, middleware []Middleware) {
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
			ws.handleGroup(route, v1, middleware) // 传递上层中间件
		}
	}
}

// setListClient 存储一个ws链接到在线列表
func (ws *WebSocket) saveClientToList(client *Client) {
	ws.list.Set(client.cId, client)
}

// delListClient 从在线列表中删除一个客户端
func (ws *WebSocket) removeClientInList(client *Client) {
	c := ws.list.Get(client.cId)
	if c == nil {
		return
	}
	if c.(*Client).cId == client.cId {
		ws.list.Remove(client.cId)
	}

	c = ws.room.Get(client.cId)
	if c == nil {
		return
	}
	if c.(*Client).cId == client.cId {
		ws.room.Remove(client.cId)
	}
}

// setList 给指定列表中添加一个客户端
func (ws *WebSocket) savaClientToRoom(client *Client) {
	ok := ws.room.Contains(client.roomId)

	if !ok {
		arr := gmap.New(true)
		arr.Set(client.cId, client)
		ws.room.Set(client.roomId, arr)
		return
	}

	ws.room.Get(client.roomId).(*gmap.Map).Set(client.cId, client)
}

// removeList 从指定列表中删除一个客户端
func (ws *WebSocket) removeClientInRoom(client *Client) {
	val := ws.room.Get(client.roomId)

	if val == nil {
		return
	}

	room := val.(*gmap.Map)

	if room.IsEmpty() {
		ws.room.Remove(client.cId) // 回收房间
		return
	}

	c := room.Get(client.cId)

	if c == nil {
		return
	}

	// 从房间中删除这个客户端
	room.Remove(client.cId)

	// 回收这个房间
	if room.IsEmpty() {
		ws.room.Remove(client.roomId)
	}
}

// sendToAll 给指定房间广播消息
func (ws *WebSocket) sendToRoom(ctx context.Context, room *gmap.Map, msg *ws.Message, filter string) {
	room.Iterator(func(cid interface{}, client interface{}) bool {
		if filter != "" {
			if filter != client.(*Client).cId {
				client.(*Client).Send(ctx, msg)
			}
		} else {
			client.(*Client).Send(ctx, msg)
		}
		return true
	})
}

// PrintRoute 打印Socket路由
func (ws *WebSocket) PrintRoute() {
	fmt.Println("|----------------------------------------------------------------- route -----------------------------------------------------------------|")
	ws.printRoute(ws.route)
	for identity, route := range ws.group {
		fmt.Println("|----------------------------------------------------------------- " + identity + " -----------------------------------------------------------------|")
		ws.printRoute(route)
	}
}

func (ws *WebSocket) printRoute(r map[string]Route) {
	fmt.Printf("| %-40v | %-50v | %-40v\n", "Action", "Controller", "Middleware")
	fmt.Println("|-----------------------------------------------------------------------------------------------------------------------------------------|")
	for action, route := range r {
		controller := utils.GetSubstr(runtime.FuncForPC(reflect.ValueOf(route.fun).Pointer()).Name(), "internal/packed/websocket/", "-fm")
		var middleware string
		for _, v2 := range route.middleware {
			middleware += utils.GetSubstr(runtime.FuncForPC(reflect.ValueOf(v2).Pointer()).Name(), "TaiYunServer/internal/packed/websocket/", ".func1") + ", "
		}
		middleware = strings.TrimSuffix(middleware, ", ") + " |"
		fmt.Printf("| %-40v | %-50v | %-40v\n", action, controller, middleware)
	}
	fmt.Println("|-----------------------------------------------------------------------------------------------------------------------------------------|")
	fmt.Println()
}
