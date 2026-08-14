package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Minimal RFC 6455 client, sufficient for Finnhub text frames.
type WSClient struct {
	conn net.Conn
	r    *bufio.Reader
	w    *bufio.Writer
	mu   sync.Mutex
}

func DialWebSocket(ctx context.Context, rawURL string) (*WSClient, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	host := u.Host
	if !strings.Contains(host, ":") {
		host += ":443"
	}
	dialer := &net.Dialer{Timeout: 15 * time.Second}
	conn, err := tls.DialWithDialer(dialer, "tcp", host, &tls.Config{ServerName: u.Hostname(), MinVersion: tls.VersionTLS12})
	if err != nil {
		return nil, err
	}
	keyBytes := make([]byte, 16)
	_, _ = rand.Read(keyBytes)
	key := base64.StdEncoding.EncodeToString(keyBytes)
	path := u.RequestURI()
	req := fmt.Sprintf("GET %s HTTP/1.1\r\nHost: %s\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Key: %s\r\nSec-WebSocket-Version: 13\r\nUser-Agent: %s/%s\r\n\r\n", path, u.Host, key, appName, appVersion)
	if _, err = conn.Write([]byte(req)); err != nil {
		conn.Close()
		return nil, err
	}
	reader := bufio.NewReader(conn)
	status, err := reader.ReadString('\n')
	if err != nil {
		conn.Close()
		return nil, err
	}
	if !strings.Contains(status, " 101 ") {
		conn.Close()
		return nil, errors.New("websocket upgrade failed: " + strings.TrimSpace(status))
	}
	headers := map[string]string{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			conn.Close()
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			headers[strings.ToLower(strings.TrimSpace(parts[0]))] = strings.TrimSpace(parts[1])
		}
	}
	expectedRaw := sha1.Sum([]byte(key + "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"))
	expected := base64.StdEncoding.EncodeToString(expectedRaw[:])
	if headers["sec-websocket-accept"] != expected {
		conn.Close()
		return nil, errors.New("invalid websocket accept key")
	}
	return &WSClient{conn: conn, r: reader, w: bufio.NewWriter(conn)}, nil
}
func (c *WSClient) WriteText(text string) error { return c.writeFrame(0x1, []byte(text)) }
func (c *WSClient) WritePing() error            { return c.writeFrame(0x9, []byte("depulse")) }
func (c *WSClient) writeFrame(op byte, payload []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return io.ErrClosedPipe
	}
	header := []byte{0x80 | op}
	n := len(payload)
	if n < 126 {
		header = append(header, byte(n)|0x80)
	} else if n <= 65535 {
		header = append(header, 126|0x80, byte(n>>8), byte(n))
	} else {
		header = append(header, 127|0x80)
		buf := make([]byte, 8)
		binary.BigEndian.PutUint64(buf, uint64(n))
		header = append(header, buf...)
	}
	mask := make([]byte, 4)
	_, _ = rand.Read(mask)
	header = append(header, mask...)
	masked := make([]byte, n)
	for i := range payload {
		masked[i] = payload[i] ^ mask[i%4]
	}
	if _, err := c.w.Write(header); err != nil {
		return err
	}
	if _, err := c.w.Write(masked); err != nil {
		return err
	}
	return c.w.Flush()
}
func (c *WSClient) ReadText(ctx context.Context) ([]byte, error) {
	var assembled []byte
	var currentOpcode byte
	for {
		if deadline, ok := ctx.Deadline(); ok {
			_ = c.conn.SetReadDeadline(deadline)
		} else {
			_ = c.conn.SetReadDeadline(time.Now().Add(35 * time.Second))
		}
		b1, err := c.r.ReadByte()
		if err != nil {
			return nil, err
		}
		b2, err := c.r.ReadByte()
		if err != nil {
			return nil, err
		}
		fin := b1&0x80 != 0
		op := b1 & 0x0f
		masked := b2&0x80 != 0
		length := uint64(b2 & 0x7f)
		if length == 126 {
			buf := make([]byte, 2)
			if _, err = io.ReadFull(c.r, buf); err != nil {
				return nil, err
			}
			length = uint64(binary.BigEndian.Uint16(buf))
		} else if length == 127 {
			buf := make([]byte, 8)
			if _, err = io.ReadFull(c.r, buf); err != nil {
				return nil, err
			}
			length = binary.BigEndian.Uint64(buf)
		}
		if length > 16<<20 {
			return nil, errors.New("websocket frame too large")
		}
		var mask []byte
		if masked {
			mask = make([]byte, 4)
			if _, err = io.ReadFull(c.r, mask); err != nil {
				return nil, err
			}
		}
		payload := make([]byte, int(length))
		if _, err = io.ReadFull(c.r, payload); err != nil {
			return nil, err
		}
		if masked {
			for i := range payload {
				payload[i] ^= mask[i%4]
			}
		}
		switch op {
		case 0x8:
			return nil, io.EOF
		case 0x9:
			_ = c.writeFrame(0xA, payload)
			continue
		case 0xA:
			continue
		case 0x1:
			currentOpcode = op
			assembled = append(assembled, payload...)
		case 0x0:
			assembled = append(assembled, payload...)
		default:
			continue
		}
		if fin && currentOpcode == 0x1 {
			return assembled, nil
		}
	}
}
func (c *WSClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return nil
	}
	_ = c.conn.SetDeadline(time.Now().Add(time.Second))
	_ = c.conn.Close()
	c.conn = nil
	return nil
}
