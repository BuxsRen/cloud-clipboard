package websocket

import (
	"CloudContent/internal/consts/websocket"
	ws "CloudContent/internal/model/websocket"
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"sync"
	"time"
)

// ping ping消息
func (c *Client) ping() {
	defer c.safeClose()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-c.tk.C:
			if !c.status {
				return
			}
			b, _ := json.Marshal(&ws.Message{Action: websocket.Ping})
			c.msgPing <- b
		}
	}
}

// safeClose 安全的关闭发送消息通道
func (c *Client) safeClose() {
	c.once.Do(func() {
		c.tk.Stop()
		c.cancel()
		c.status = false
		_ = c.socket.Close()
		socket.disconnect(gctx.New(), c)
		time.Sleep(time.Second)
		close(c.msg)
		close(c.msgPing)
	})
}

// write 写入消息协程
func (c *Client) write() {
	defer c.safeClose()

	for {
		select {
		case msg, ok := <-c.msg: // 普通消息
			// 通道被关闭，退出
			if !ok {
				return
			}
			e := c.socket.WriteMessage(ghttp.WsMsgText, msg)
			// 消息写入失败，退出
			if e != nil {
				continue
			}
		case msg, ok := <-c.msgPing: // PING 消息
			if !ok {
				return
			}
			e := c.socket.WriteMessage(ghttp.WsMsgPing, msg)
			if e != nil {
				continue
			}
		}
	}
}

// read 读取消息协程
func (c *Client) read() {
	defer c.safeClose()

	for {
		_, b, err := c.socket.ReadMessage()
		if err != nil {
			break
		}

		// 异步执行控制器
		go func(b []byte) {
			ctx := gctx.New()

			defer func(ctx context.Context, client *Client) {
				if e := recover(); e != nil {
					socket.iSocket.Error(ctx, client, e)
				}
			}(ctx, c)

			msg := &ws.Message{}
			err = json.Unmarshal(b, msg)
			// 无用消息不处理
			if err != nil {
				//c.Send(ctx, &Message{Action: websocket.Error, Msg: "request data format error", Code: websocket.CodeError})
				_ = c.socket.Close() // 发无用消息强制断开客户端连接
				return
			}

			// 查询事件路由
			route, ok := socket.route[msg.Action]
			if !ok {
				// 查询事件路由组
				group, ok := socket.group[websocket.Identity]

				// 未绑定事件组
				if !ok {
					//c.Send(ctx, &Message{Action: websocket.Error, Msg: "user route 404 not found", Code: websocket.CodeNotFound})
					_ = c.socket.Close() // 发无用消息强制断开客户端连接
					return
				}

				route, ok = group[msg.Action]
				// 未绑定路由分组事件组
				if !ok {
					//c.Send(ctx, &Message{Action: websocket.Error, Msg: "404 not found", Code: websocket.CodeNotFound})
					_ = c.socket.Close() // 发无用消息强制断开客户端连接
					return
				}
			}

			// 执行路由中间件
			if route.middleware != nil {
				var pass = true
				for _, middleware := range route.middleware {
					if middleware != nil {
						ctx, ok = middleware(ctx, c, msg)
						if !ok {
							pass = false
							break
						}
					}
				}
				if !pass {
					return
				}
			}

			route.fun(ctx, c, msg)
		}(b)
	}
}

// Send 给当前客户端发送消息
func (c *Client) Send(ctx context.Context, msg *ws.Message) {
	msg.Time = gtime.Now().Unix()
	if msg.Msg == "" {
		msg.Msg = "success"
	}

	b, err := json.Marshal(msg)
	if err != nil {
		c.Send(ctx, &ws.Message{Action: websocket.Error, Msg: err.Error(), Code: websocket.CodeError})
		return
	}

	if c.status {
		c.msg <- b
	}
}

// Close 关闭当前客户端连接
func (c *Client) Close() {
	if c.status {
		c.safeClose()
	}
}

// GetMutex 获得锁
func (c *Client) GetMutex() *sync.Mutex {
	return &c.mutex
}

// GetCtx 获得CTX
func (c *Client) GetCtx() context.Context {
	return c.ctx
}

// GetRoomId 获得RoomId
func (c *Client) GetRoomId() string {
	return c.roomId
}

// GetCId 获得CId
func (c *Client) GetCId() string {
	return c.cId
}

// GetIp 获得IP
func (c *Client) GetIp() string {
	return c.ip
}

// GetUA 获得UA
func (c *Client) GetUA() string {
	return c.userAgent
}
