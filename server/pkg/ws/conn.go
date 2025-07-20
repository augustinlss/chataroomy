package ws

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
)

const (
	OpcodeContinuation = 0x0
	OpcodeText         = 0x1
	OpcodeBinary       = 0x2
	OpcodeClose        = 0x8
	OpcodePing         = 0x9
	OpcodePong         = 0xA
)

type WebSocket struct {
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	writeMut sync.Mutex
	isClient bool
}

type Frame struct {
	Fin     bool
	Opcode  byte
	Masked  bool
	MaskKey [4]byte
	Payload []byte
}

func (ws *WebSocket) isClientConn() bool {
	return ws.isClient
}

func (ws *WebSocket) isServerConn() bool {
	return !ws.isClient
}

func NewWebSocket(conn net.Conn) *WebSocket {
	return &WebSocket{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func (ws *WebSocket) Close() error {
	return ws.conn.Close()
}
