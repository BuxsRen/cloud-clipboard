package clipboard

import (
	wsc "cloud-clipboard/internal/consts/websocket"
	ws "cloud-clipboard/internal/model/websocket"
	"cloud-clipboard/internal/service"
	"cloud-clipboard/utility/client"
	"context"
	"github.com/atotto/clipboard"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtimer"
	"time"
)

// GetClipboardToRoom 读取系统剪贴板并同步到房间
func GetClipboardToRoom() {
	gtimer.SetInterval(gctx.New(), time.Second, func(ctx context.Context) {
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
}

// WriteClipboardText 写入系统剪贴板
func WriteClipboardText(text string) {
	err := clipboard.WriteAll(text)
	if err != nil {
		g.Log().Error(gctx.New(), err)
		return
	}
}
