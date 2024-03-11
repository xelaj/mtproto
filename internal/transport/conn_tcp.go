package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"syscall"
	"time"
)

type tcpConn struct {
	conn *net.TCPConn
}

func NewTCP(host string) (io.ReadWriteCloser, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return nil, fmt.Errorf("resolving tcp: %w", err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, fmt.Errorf("dialing tcp: %w", err)
	}

	return &tcpConn{conn: conn}, nil
}

func (t *tcpConn) Close() error {
	return t.conn.Close()
}

func (t *tcpConn) Write(b []byte) (int, error) {
	return t.conn.Write(b)
}

func (t *tcpConn) Read(b []byte) (int, error) {
	for {
		if err := t.conn.SetReadDeadline(time.Now().Add(1 * time.Second)); err != nil {
			panic(err)
		}

		n, err := t.conn.Read(b)
		if err != nil {
			if e, ok := err.(*net.OpError); ok {
				if errors.Is(err, syscall.ECONNRESET) {
					return n, io.EOF
				}
				switch errText := e.Err.Error(); errText {
				case "i/o timeout":
					return n, context.DeadlineExceeded

				case "use of closed network connection":
					return n, io.EOF
				}
			}
			switch err {
			case io.EOF:
				return n, err
			default:
				return n, fmt.Errorf("unexpected error: %w", err)
			}
		}
		return n, nil
	}
}
