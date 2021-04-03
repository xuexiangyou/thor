package pack

import (
	"bufio"
	"errors"
	"io"
	"net"
)

const (
	defaulReaderSize = 8 * 1024
	maxPlayloadLen 	int = 1<<16 - 1 // 计算两个字段的最大值
)

var (
	ErrBadConn       = errors.New("connection was bad")
)

type PacketIO struct {
	rb *bufio.Reader
	wb io.Writer
}

func NewPacketIO(conn net.Conn) *PacketIO {
	p := new(PacketIO)
	p.rb = bufio.NewReaderSize(conn, defaulReaderSize) // 读取缓存
	p.wb = conn
	return p
}

func (p *PacketIO) ReadPacket() ([]byte, error) {
	header := []byte{0, 0} // 定义消息头、目前只存消息内容的大小、其他消息头内容可以自定义扩展添加
	if _, err := io.ReadFull(p.rb, header); err != nil {
		return nil, ErrBadConn
	}

	length := int(uint16(header[0]) | uint16(header[1])<<8)
	data := make([]byte, length)

	if _, err := io.ReadFull(p.rb, data); err != nil {
		return nil, ErrBadConn
	} else {
		if length < maxPlayloadLen {
			return data, nil
		}
		var buf []byte
		buf, err := p.ReadPacket()
		if err != nil {
			return nil, ErrBadConn
		} else {
			return append(data, buf...), nil
		}
	}
}

func (p *PacketIO) WritePacket(data []byte) error {
	length := len(data) - 2

	for length >= maxPlayloadLen {
		data[0] = 0xff
		data[1] = 0xff

		if n, err := p.wb.Write(data[:2+maxPlayloadLen]); err != nil {
			return ErrBadConn
		} else if n != (2 + maxPlayloadLen) {
			return ErrBadConn
		} else {
			length -= maxPlayloadLen
			data = data[maxPlayloadLen:]
		}
	}

	data[0] = byte(length)
	data[1] = byte(length >> 8)

	if n, err := p.wb.Write(data); err != nil {
		return ErrBadConn
	} else if n != len(data) {
		return ErrBadConn
	} else {
		return nil
	}
}

