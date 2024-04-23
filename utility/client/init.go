package client

import (
	"context"
	"sync"
)

func Run(si SocketInterface) {
	if socket == nil {

		socket = &Socket{
			iSocket:     si,
			routeGroup:  make(map[string]*RouteGroup),
			route:       make(map[string]Route),
			group:       make(map[string]map[string]Route),
			againStatus: true,
			mutex:       sync.Mutex{},
		}

		socket.iSocket.Route()

		// 绑定路由组路由事件
		for k, v := range socket.routeGroup {
			route := make(map[string]Route)

			socket.handleRouteGroup(&route, v)

			socket.group[k] = route
		}

		socket.Client = &Client{
			conn:    nil,
			msgChan: make(chan []byte, 10),
			Status:  false,
			once:    sync.Once{},
			mutex:   sync.Mutex{},
		}

		socket.again()
	}
}

// NewClient 初始化一个ws客户端
func NewClient(ctx context.Context) {
	socket.Client = &Client{
		conn:    nil,
		msgChan: make(chan []byte, 10),
		Status:  false,
		once:    sync.Once{},
		mutex:   sync.Mutex{},
	}

	socket.Client.ctx, socket.Client.cancel = context.WithCancel(ctx)

	go socket.Client.connect()
}
