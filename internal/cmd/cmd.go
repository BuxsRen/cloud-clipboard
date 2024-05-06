package cmd

import (
	"cloud-clipboard/internal/model/app"
	wsClient "cloud-clipboard/internal/packed/client"
	"cloud-clipboard/internal/packed/socket"
	"cloud-clipboard/utility/client"
	"cloud-clipboard/utility/clipboard"
	"cloud-clipboard/utility/websocket"
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			switch parser.GetOpt("mode").String() {
			case "client":
				clipboard.GetClipboardToRoom()

				host := parser.GetOpt("host").String()
				if host != "" {
					app.Host = host
				}

				room := parser.GetOpt("room").String()
				if room != "" {
					app.RoomId = room
				}
				client.Run(new(wsClient.Socket))
			default:
				websocket.Init(new(socket.Server))

				var port = 5123
				val := parser.GetOpt("port").Int()
				if val > 0 {
					port = val
				}

				s1 := g.Server("ws")
				s1.BindHandler("/cloud", websocket.HandleClient)
				s1.SetPort(port)
				s1.Run()
			}
			return nil
		},
	}
)
