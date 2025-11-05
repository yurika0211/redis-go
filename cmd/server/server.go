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
 * func NewServer(l net.Listener) *Server
 * func (s *Server) Start() error
 * func (s *Server) Stop()
 * func (s *Server) handleConn(conn net.Conn)
 */
func NewServer(l net.Listener) *Server {
	return &Server{
		listener: l,
		quit: make(chan chan struct{}),
	}
}

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
 * @brief 
 */
func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
}
