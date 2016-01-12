package goapi

import (
    "bufio"
    "net"
)

type bufferedConn struct {
    r        *bufio.Reader
    net.Conn // So that most methods are embedded
}

func newBufferedConn(c net.Conn) bufferedConn {
    return bufferedConn{bufio.NewReader(c), c}
}
func (bc bufferedConn) Peek(n int) ([]byte, error) {
    return bc.r.Peek(n)
}
func (bc bufferedConn) Read(b []byte) (n int, err error) {
    return bc.r.Read(b)
}
func (bc bufferedConn) ReadString(b byte) (string, error) {
    return bc.r.ReadString(b)
}
func (bc bufferedConn) Buffered() int {
    return bc.r.Buffered()
}
