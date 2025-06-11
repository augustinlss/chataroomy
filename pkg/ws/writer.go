package ws

import (
	"fmt"
)

func (ws *WebSocketConn) WriteFrame(frame *Frame) error {
	// TODO implement frame writing logic
	// This is a placeholder for the actual implementation

	ws.writeMut.Lock()

	var data []byte

	return fmt.Errorf("Write method not implemented")
}

func (ws *WebSocketConn) WriteMessage(opcode byte, data []byte) error {
	// for now we will send all messages in a single frame
	// in the future, i should probably split messages into
	// multiple framea
	frame := &Frame{
		Fin:     true,
		Opcode:  opcode,
		Payload: data,
		Masked:  false,
	}
	return ws.WriteFrame(frame)
}

func (ws *WebSocketConn) WriteText() error {

}

func (ws *WebSocketConn) WriteBinary() error {

}

func (ws *WebSocketConn) WritePing() error {

}

func (ws *WebSocketConn) WritePong() error {

}
