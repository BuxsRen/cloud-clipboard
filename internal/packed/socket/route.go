package socket

import (
	wsc "cloud-clipboard/internal/consts/websocket"
	"cloud-clipboard/internal/middleware"
	"cloud-clipboard/internal/packed/socket/app"
	"cloud-clipboard/utility/websocket"
)

func (s *Server) Route() {
	cApp := app.App{}

	// 使用普通路由方式
	websocket.Routes(wsc.GetContent, cApp.GetContent)
	websocket.Routes(wsc.SetContent, cApp.SetContent).Middleware(middleware.Limit) // 该控制器只对该 Action 生效

	// 使用路由组方式
	//websocket.Groups(wsc.Identity, func(group *websocket.RouteGroup) {
	//	group.Route(wsc.GetContent, cApp.GetContent)
	//
	//	group.Group(func(group *websocket.RouteGroup) {
	//		group.Middleware(middleware.Limit) // 使用消息限流中间件 该组中的所有 Action 以及下级所有的 Action 都受这个中间件控制
	//
	//		group.Route(wsc.SetContent, cApp.SetContent)
	//
	//		// 继续嵌套路由组
	//		group.Group(func(group *websocket.RouteGroup) {
	//			// 中间件...  会先执行上级的中间件
	//
	//			// 分组路由...
	//		})
	//	})
	//})
}
