package main

import (
	_ "CloudContent/internal/logic"
	_ "CloudContent/internal/packed"
	"github.com/gogf/gf/v2/os/gctx"

	"CloudContent/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
