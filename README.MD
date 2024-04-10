# 云内容 WebSocket案例

|   语言    |     框架      |   WebSocket   |
|:-------:|:-----------:|:-------------:|
| Go 1.22 | GoFrame 2.5 | gorilla 1.5.1 |

#### 在线demo
> https://fan.jx.cn:82/ct.html
 
#### WebSocket 底层封装
- utility/websocket


#### WebSocket 业务层
- internal/packed/socket

#### 路由与中间件
```go
package server

import (
	wsc "CloudContent/internal/consts/websocket"
	"CloudContent/utility/websocket"
	"context"
)

func (s *Server) Route() {
	route := websocket.Middleware(func(ctx context.Context, client *websocket.Client, msg *websocket.Message) (context.Context, bool) {
		// 中间件实现
		return ctx, true
	})

	route.Routes(wsc.Auth, s.Auth)// 传入 Actoin 与 Controller
}
```

#### 路由组与中间件
```go
package server

import (
	wsc "CloudContent/internal/consts/websocket"
	"CloudContent/internal/packed/socket/app"
	"CloudContent/utility/websocket"
	"context"
)

func (s *Server) Route() {
	cApp := app.App{}

	websocket.Groups(wsc.Identity, func(group *websocket.RouteGroup) {
		group.Middleware(func(ctx context.Context, client *websocket.Client, msg *websocket.Message) (context.Context, bool) {
			// 中间件实现
			return ctx, true
		})
		
		group.Route(wsc.GetContent, cApp.GetContent)
		group.Route(wsc.SetContent, cApp.SetContent)

		// 继续嵌套路由组
		group.Group(func(group *websocket.RouteGroup) {

		})
	})
}
```

#### 房间内广播消息
```go
package server

import (
	wsc "CloudContent/internal/consts/websocket"
	"CloudContent/internal/model/content"
	"CloudContent/internal/service"
	"CloudContent/utility/websocket"
	"context"
	"github.com/gogf/gf/v2/util/gconv"
)

func (s *Server) Test(ctx context.Context, client *websocket.Client, _ *websocket.Message) {
	_ = websocket.GetSocketServer().SendMessageToRoom(
		ctx, 
		client.GetRoomId(), 
		&websocket.Message{Action: wsc.Content, Data: "消息内容"}, 
		client.GetCId(), 
	)
}
```