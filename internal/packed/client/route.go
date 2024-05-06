package client

import (
	wsc "cloud-clipboard/internal/consts/websocket"
	"cloud-clipboard/internal/packed/client/app"
	"cloud-clipboard/utility/client"
)

// Route 事件路由，注册Action绑定到Func
func (s *Socket) Route() {
	cApp := app.App{}

	client.Routes(wsc.Content, cApp.Content)
}
