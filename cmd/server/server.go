package core

import (
	"net"
	"bufio"
	"io"
	"log"
	"com.ityurika/go-redis-clone/internal/command"
	"com.ityurika/go-redis-clone/internal/protocol"
)

type Server struct {
	listener net.Listener
	quit     chan (chan struct{})
}

/**
 * 创建一个服务器
 * @param l net.Listener 监听器
 */
func NewServer(l net.Listener) *Server {
	return &Server{
		listener: l,
		quit: make(chan chan struct{}),
	}
}

/**
 * 启动服务器
 * @param server *Server 服务器
 */
func (s *Server) Start() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case q := <-s.quit:
				close(q)
				return nil
			default:
				return err
			}
		}
		go s.handleConn(conn)
	}
}

/**
 * 处理连接
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
 * @brief 停止服务器
 * @param server *Server 服务器
 */
func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
}
