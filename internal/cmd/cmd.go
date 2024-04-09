package cmd

import (
	"CloudContent/utility/websocket"
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
			s1 := g.Server("ws")
			s1.BindHandler("/cloud", websocket.HandleClient)
			s1.SetPort(5123)
			s1.Run()
			return nil
		},
	}
)
