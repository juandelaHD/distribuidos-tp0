package common

import (
	"fmt"
	"net"
)

// sendAll writes all bytes in buf to the provided connection, retrying on partial writes.
func sendAll(conn net.Conn, buf []byte) error {
	sent := 0
	for sent < len(buf) {
		n, err := conn.Write(buf[sent:])
		if err != nil {
			return fmt.Errorf("sendAll: %w", err)
		}
		sent += n
	}
	return nil
}

// recvAll reads exactly n bytes from the provided connection, retrying on partial reads.
func recvAll(conn net.Conn, n int) ([]byte, error) {
	buf := make([]byte, n)
	received := 0
	for received < n {
		r, err := conn.Read(buf[received:])
		if err != nil {
			return nil, fmt.Errorf("recvAll: %w", err)
		}
		received += r
	}
	return buf, nil
}
