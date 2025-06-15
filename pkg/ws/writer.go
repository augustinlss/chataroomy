package ws

import (
	"fmt"
)

func (ws *WebSocket) WriteFrame(frame *Frame) error {
	// TODO implement frame writing logic
	// This is a placeholder for the actual implementation

	ws.writeMut.Lock()
	defer ws.writeMut.Unlock()

	if frame == nil {
		return fmt.Errorf("frame cannot be nil")
	}

	if frame.Opcode < 0 || frame.Opcode > 15 {
		return fmt.Errorf("invalid opcode: %d", frame.Opcode)
	}

	firstByte := byte(frame.Opcode & 0x0F)

	if frame.Fin {
		firstByte |= 0x80
	}

	return fmt.Errorf("Write method not implemented")
}

func (ws *WebSocket) WriteMessage(opcode byte, data []byte) error {
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

func (ws *WebSocket) WriteText() error {

}

func (ws *WebSocket) WriteBinary() error {

}

func (ws *WebSocket) WritePing() error {

}

func (ws *WebSocket) WritePong() error {

}
