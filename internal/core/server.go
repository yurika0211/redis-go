package core

import (
	"bufio"
	"io"
	"log"
	"net"
	"sync"

	"com.ityurika/go-redis-clone/internal/command"
	"com.ityurika/go-redis-clone/internal/protocol"
)

type Server struct {
	listener net.Listener
	quit     chan struct{}
	wg       sync.WaitGroup
}

/**
 * 创建一个服务器
 * @param l net.Listener 监听器
 */
func NewServer(l net.Listener) *Server {
	return &Server{
		listener: l,
		quit:     make(chan struct{}),
	}
}

/**
 * 启动服务器
 * @return error 错误
 */
func (s *Server) Start() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return nil
			default:
				return err
			}
		}
		s.wg.Add(1)
		go func(c net.Conn) {
			defer s.wg.Done()
			s.handleConn(c)
		}(conn)
	}
}

/**
 * 处理一个连接
 * @param conn net.Conn 连接
 */
func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		arr, err := protocol.ParseArray(reader)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("parse error: %v", err)
			return
		}
		if len(arr) == 0 {
			continue
		}
		cmd := arr[0]
		args := arr
		command.HandleCommand(conn, cmd, args)
	}
}

/**
 * 停止服务器
 */
func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
}
