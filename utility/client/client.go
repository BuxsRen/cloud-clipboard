package client

import (
	wsc "CloudContent/internal/consts/websocket"
	"CloudContent/internal/model/content"
	ws "CloudContent/internal/model/websocket"
	"context"
	"encoding/json"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

func (c *Client) connect() {
	client := gclient.NewWebSocket()
	client.HandshakeTimeout = 10 * time.Second // 设置超时时间

	conn, _, err := client.Dial(content.Host+"/ws/cloud?room_id="+content.RoomId, nil)
	if err != nil {
		log.Println(err)
		c.safeClose()
		return
	}

	c.conn = conn
	c.Status = true

	go c.write()
	go c.read()

	socket.iSocket.Login(c.ctx)
}

// safeClose 安全的关闭发送消息通道
func (c *Client) safeClose() {
	c.once.Do(func() {
		if c.Status {
			c.Status = false
			c.cancel()
			socket.iSocket.Logout(gctx.New())
		}
		if c.conn != nil {
			_ = c.conn.Close()
		}
		time.Sleep(time.Second)
		socket.againStatus = true
		close(c.msgChan)
	})
}

func (c *Client) write() {
	defer c.safeClose()

	for {
		select {
		case msg, ok := <-c.msgChan:
			if !ok || !c.Status {
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				continue
			}
		}
	}
}

func (c *Client) read() {
	defer c.safeClose()

	for {
		_, b, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		go func(b []byte) {
			ctx := gctx.New()

			defer func(ctx context.Context) {
				if e := recover(); e != nil {
					socket.iSocket.Error(ctx, e)
				}
			}(ctx)

			msg := &ws.Message{}
			err = json.Unmarshal(b, msg)
			// 无用消息不处理
			if err != nil {
				return
			}

			// 查询事件路由
			route, ok := socket.route[msg.Action]
			if !ok {
				// 查询事件路由组
				group, ok := socket.group[wsc.Identity]

				// 未绑定事件组
				if !ok {
					return
				}

				route, ok = group[msg.Action]
				// 未绑定路由分组事件组
				if !ok {
					return
				}
			}

			// 执行路由中间件
			if route.middleware != nil {
				var pass = true
				for _, middleware := range route.middleware {
					if middleware != nil && !middleware(ctx, msg) {
						pass = false
						break
					}
				}
				if !pass {
					return
				}
			}

			route.fun(ctx, msg)
		}(b)
	}
}
