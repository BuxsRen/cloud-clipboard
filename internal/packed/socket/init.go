package socket

import "cloud-clipboard/utility/websocket"

type Server struct{}

var _ websocket.ISocket = (*Server)(nil)
