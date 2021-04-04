package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/xuexiangyou/thor/socket/pack"
)

type Server struct {
	addr string

	listener net.Listener

	wg sync.WaitGroup

	running bool

}

func NewServer(addr string) (*Server, error) {
	s := new(Server)
	s.addr = addr

	var err error
	netProto := "tcp"
	s.listener, err = net.Listen(netProto, s.addr)
	if err != nil {
		return nil, err
	}
	s.wg = sync.WaitGroup{}
	return s, nil
}

func (s *Server) Close() {
	s.running = false
	if s.listener != nil {
		s.listener.Close()
	}
}

func (s *Server) Run() error {
	s.running = true
	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("接收链接", err.Error())
			continue
		}

		s.wg.Add(1)

		go s.onConn(conn)
	}
	s.wg.Wait()
	return nil
}

func (s *Server) onConn(c net.Conn) {
	conn := s.newClientConn(c)
	conn.Run()
}

func (s *Server) newClientConn(co net.Conn) *ClientConn{
	c := new(ClientConn)
	c.c = co
	c.pkg = pack.NewPacketIO(co)
	c.closed = false
	c.proxy = s
	return c
}
