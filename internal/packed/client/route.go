package client

import (
	wsc "CloudContent/internal/consts/websocket"
	"CloudContent/internal/packed/client/app"
	"CloudContent/utility/client"
)

// Route 事件路由，注册Action绑定到Func
func (s *Socket) Route() {
	cApp := app.App{}

	client.Routes(wsc.Content, cApp.Content)
}
