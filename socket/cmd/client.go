package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/xuexiangyou/thor/socket/pack"
)

type ClientConn struct {
	pkg *pack.PacketIO
	c net.Conn
}

func main() {
	startClient("127.0.0.1:8089")
}

func startClient(ip string) {
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		fmt.Println("Error dialing", err.Error())
		return // 终止程序
	}

	c := newClientConn(conn)

	for  {
		// 随机生产
		r := rand.New(rand.NewSource(time.Now().Unix()))
		n := r.Intn(100000) + 60000
		message := GetRandomString(n)
		err := c.ClientWrite(message)
		if err != nil {
			if strings.Contains(err.Error(), "connection was bad") {
				break
			}
			fmt.Println(err)
		}
		fmt.Println("发送消息", len(message))
		time.Sleep(2 * time.Second)
		// 读取数据
		result, err := c.pkg.ReadPacket()
		if err != nil {
			fmt.Println("接收消息错误", err.Error())
		}
		fmt.Println("读取消息", len(string(result)))
	}
}

func newClientConn(co net.Conn) *ClientConn {
	c := new(ClientConn)
	c.c = co
	c.pkg = pack.NewPacketIO(co)
	return c
}

func (c *ClientConn) ClientWrite(message string) error {
	length := len(message)
	data := make([]byte, 2, 2 + length)

	data = append(data, message...)

	if err := c.pkg.WritePacket(data); err != nil {
		return err
	} else {
		return nil
	}
}

func GetRandomString(n int) string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	var result []byte
	for i := 0; i < n; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
