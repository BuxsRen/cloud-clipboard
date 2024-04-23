package app

import (
	ws "CloudContent/internal/model/websocket"
	"CloudContent/internal/service"
	"CloudContent/utility/client"
	"context"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/gogf/gf/v2/util/gconv"
)

type App struct{}

// Content 获取内容
func (a *App) Content(ctx context.Context, msg *ws.Message) {
	mutex := client.GetSocketClient().GetMutex()
	mutex.Lock()
	defer mutex.Unlock()
	text := gconv.String(msg.Data)

	_ = service.Cache().SetCache(ctx, "text", text)

	err := clipboard.WriteAll(text)
	if err != nil {
		fmt.Println("Failed to write to clipboard:", err)
		return
	}
}
