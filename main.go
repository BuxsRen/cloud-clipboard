package main

import (
	_ "CloudContent/internal/logic"
	_ "CloudContent/internal/packed"
	"CloudContent/internal/packed/socket"
	"CloudContent/utility/websocket"

	"github.com/gogf/gf/v2/os/gctx"

	"CloudContent/internal/cmd"
)

func main() {
	websocket.Init(new(socket.Server))

	cmd.Main.Run(gctx.GetInitCtx())
}
