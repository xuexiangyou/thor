package main

import (
	"fmt"
	"log"
	"net"

	"github.com/xuexiangyou/thor/socket/pack"
)

func main() {
	StartServer("127.0.0.1:8089")
}

type ClientConn struct {
	pkg *pack.PacketIO
	c net.Conn
}

func StartServer(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("Error listening", err.Error())
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("error accepting", err.Error())
			return
		}

		go onConn(conn)
	}
}

func onConn(c net.Conn) {
	conn := newServerConn(c) // 新建一个conn

	// 链接逻辑处理
	conn.Run()
}

func (c *ClientConn) Run() {
	for  {
		data, err := c.readPacket()
		if err != nil {
			return
		}
		// 写入数据
		err = c.writeMsg(string(data))
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (c *ClientConn) writeMsg(message string) error {
	length := len(message)
	data := make([]byte, 2, 2 + length)

	data = append(data, message...)

	if err := c.pkg.WritePacket(data); err != nil {
		return err
	} else {
		return nil
	}
}

func newServerConn(co net.Conn) *ClientConn {
	c := new(ClientConn)
	c.c = co

	c.pkg = pack.NewPacketIO(co)

	return c
}

func (c *ClientConn) readPacket() ([]byte, error) {
	return c.pkg.ReadPacket()
}

