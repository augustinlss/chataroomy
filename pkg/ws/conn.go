package ws

import (
	"bufio"
	"net"
	"sync"
)

const (
	OpcodeContinuation = 0x0
	OpcodeText         = 0x1
	OpcodeBinary       = 0x2
	OpcodeClose        = 0x8
	OpcodePing         = 0x9
	OpcodePong         = 0xa
)

type WebSocketConn struct {
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	writeMut sync.Mutex
}

type Frame struct {
	Fin     bool
	Opcode  byte
	Masked  bool
	Payload []byte
}

func NewWebSocketConn(conn net.Conn) *WebSocketConn {
	return &WebSocketConn{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func (ws *WebSocketConn) ReadFrame() (opcode byte, payload []byte, err error) {
	// TODO implemenet frame reading logic
}
