package server

import (
	"fmt"
	"net"
	"runtime"

	"github.com/xuexiangyou/thor/socket/pack"
)

type ClientConn struct {
	c net.Conn

	pkg *pack.PacketIO

	proxy *Server

	closed bool
}

func (c *ClientConn) readPacket() ([]byte, error) {
	return c.pkg.ReadPacket()
}

func (c *ClientConn) writePacket(data []byte) error {
	return c.pkg.WritePacket(data)
}

func (c *ClientConn) Close() error {
	if c.closed {
		return nil
	}
	c.c.Close()

	c.closed = true
	return nil
}

func (c *ClientConn) Run() {
	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			const size = 4096
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			fmt.Println(err.Error(), buf) // todo 可以通过存日志后面线上出问题调试
		}
		c.Close()
		c.proxy.wg.Done()
	}()

	for  {
		data, err := c.readPacket()
		if err != nil {
			fmt.Println("读取数据失败", err.Error())
			return
		}
		// 写入数据
		err = c.writeMsg(string(data))
		if err != nil {
			fmt.Println("写入数据失败", err.Error())
			return
		}

		if !c.proxy.running { // 如果关闭了直接返回
			return
		}
	}
}

func (c *ClientConn) writeMsg(message string) error {
	length := len(message)
	data := make([]byte, 2, 2 + length)
	data = append(data, message...)
	if err := c.writePacket(data); err != nil {
		return err
	} else {
		return nil
	}
}


