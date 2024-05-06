package main

import (
	_ "cloud-clipboard/internal/logic"
	_ "cloud-clipboard/internal/packed"
	"github.com/gogf/gf/v2/os/gctx"

	"cloud-clipboard/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
