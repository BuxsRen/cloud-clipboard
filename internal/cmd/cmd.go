package cmd

import (
	wsc "CloudContent/internal/consts/websocket"
	"CloudContent/internal/model/content"
	ws "CloudContent/internal/model/websocket"
	wsClient "CloudContent/internal/packed/client"
	"CloudContent/internal/packed/socket"
	"CloudContent/internal/service"
	"CloudContent/utility/client"
	"CloudContent/utility/websocket"
	"context"
	"github.com/atotto/clipboard"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gtimer"
	"time"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			switch parser.GetOpt("mode").String() {
			case "client":
				gtimer.SetInterval(ctx, time.Second, func(ctx context.Context) {
					mutex := client.GetSocketClient().GetMutex()
					mutex.Lock()
					defer mutex.Unlock()

					text, err := clipboard.ReadAll()
					val, _ := service.Cache().GetCache(ctx, "text")
					if err == nil && text != "" && val.String() != text {
						_ = service.Cache().SetCache(ctx, "text", text)
						client.GetSocketClient().Send(ctx, &ws.Message{Action: wsc.SetContent, Data: g.Map{"content": text}})
					}
				})

				room := parser.GetOpt("room").String()
				if room != "" {
					content.RoomId = room
				}
				client.Run(new(wsClient.Socket))
			default:
				websocket.Init(new(socket.Server))

				s1 := g.Server("ws")
				s1.BindHandler("/cloud", websocket.HandleClient)
				s1.SetPort(5123)
				s1.Run()
			}
			return nil
		},
	}
)
