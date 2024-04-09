package socket

import (
	wsc "CloudContent/internal/consts/websocket"
	"CloudContent/internal/packed/socket/app"
	"CloudContent/utility/websocket"
)

func (s *Server) Route() {
	cApp := app.App{}

	websocket.Groups(wsc.Identity, func(group *websocket.RouteGroup) {
		group.Route(wsc.GetContent, cApp.GetContent)
		group.Route(wsc.SetContent, cApp.SetContent)
	})
}
