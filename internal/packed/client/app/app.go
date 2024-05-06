package app

import (
	ws "cloud-clipboard/internal/model/websocket"
	"cloud-clipboard/internal/service"
	"cloud-clipboard/utility/client"
	"cloud-clipboard/utility/clipboard"
	"context"
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

	clipboard.WriteClipboardText(text)
}
