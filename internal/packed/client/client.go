package client

import (
	"context"
	"log"
)

// Login Ws连接成功
func (s *Socket) Login(ctx context.Context) {
	log.Println("√ 已连接到服务端")
}

// Logout Ws断开连接
func (s *Socket) Logout(ctx context.Context) {
	log.Println("× 连接断开")

}

// Error 控制器panic错误处理
func (s *Socket) Error(_ context.Context, _ any) {}
