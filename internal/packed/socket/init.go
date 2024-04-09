package socket

import "CloudContent/utility/websocket"

type Server struct{}

var _ websocket.ISocket = (*Server)(nil)
