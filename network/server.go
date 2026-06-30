package network

import (
	"fmt"
	"time"
)

type ServerOpts struct {
	Transport []Transport
}

type Server struct {
	ServerOpts

	rpcCh chan RPC

	quitCh chan struct{}
}

func NewServer(opts ServerOpts) *Server {
	return &Server{
		ServerOpts: opts,
		rpcCh:      make(chan RPC),
		quitCh:     make(chan struct{}, 1),
	}
}

func (s *Server) Start() {
	s.initTransport()
	ticker := time.NewTicker(5 * time.Second)
free:
	for {
		select {
		case rpc := <-s.rpcCh:
			fmt.Printf("%+v\n", rpc)
		case <-s.quitCh:
			break free
		case <-ticker.C:
			fmt.Println("do stuff every x seconds")

		}

	}

	fmt.Println("Server shutdown")
}

func (s *Server) initTransport() {
	// 步骤 1：遍历所有的传输通道（比如可能有本地通道、TCP通道、UDP通道等）
	for _, tr := range s.Transport {

		// 步骤 2：用 "go" 关键字，为当前这个通道单独启动一个独立的后台“搬运工”协程
		go func(tr Transport) {

			// 步骤 3：这是一个死循环！
			// tr.Consume() 拿到了我们在 LocalTransport 里看到的那个容量 1024 的收件箱（consumeCh）。
			// 只要这个收件箱里有新消息（rpc），循环就会抓住它。
			for rpc := range tr.Consume() {

				// 步骤 4：搬运工把抓到的消息，顺手扔进服务器的总传送带 `s.rpcCh` 中
				s.rpcCh <- rpc
			}
		}(tr) // 这里的 (tr) 是把当前的通道当成参数传进去，确保搬运工没有认错通道
	}
}
