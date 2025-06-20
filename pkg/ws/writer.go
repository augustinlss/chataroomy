package ws

import (
	"encoding/binary"
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

	var header []byte
	firstByte := byte(frame.Opcode & 0x0F)

	if frame.Fin {
		firstByte |= 0x80
	}

	header = append(header, firstByte)

	payloadLen := len(frame.Payload)
	secondByte := byte(0)
	if frame.Masked {
		secondByte |= 0x80
	}

	if payloadLen < 126 {
		secondByte |= byte(payloadLen)
		header = append(header, secondByte)
	} else if payloadLen < 65536 {
		secondByte |= 126
		header = append(header, secondByte)
		extLen := make([]byte, 2)
		binary.BigEndian.PutUint16(extLen, uint16(payloadLen))
		header = append(header, extLen...)
	} else {
		secondByte |= 127
		header = append(header, secondByte)
		extLen := make([]byte, 8)
		binary.BigEndian.PutUint64(extLen, uint64(payloadLen))
		header = append(header, extLen...)
	}

	if frame.Masked {
		header = append(header, frame.MaskKey[:]...)
	}

	if _, err := ws.conn.Write(header); err != nil {
		return fmt.Errorf("failed to write frame header: %w", err)
	}

	if len(frame.Payload) > 0 {
		payload := frame.Payload
		if frame.Masked {
			payload = make([]byte, len(frame.Payload))
			for i, b := range frame.Payload {
				payload[i] = b ^ frame.MaskKey[i%4]
			}
		}

		if _, err := ws.conn.Write(payload); err != nil {
			return fmt.Errorf("failed to write frame payload: %w", err)
		}
	}

	return nil
}

func (ws *WebSocket) WriteMessage(opcode byte, data []byte) error {
	// for now will send all messages in a single frame
	// in the future, i should probably split messages into
	// multiple frames

	masked := ws.isClientConn()
	var maskKey [4]byte
	if masked {
		// Generate a masking key if this is a client connection
		maskKey = [4]byte{0x11, 0x22, 0x33, 0x44} // Placeholder: Replace with actual random
	}

	frame := &Frame{
		Fin:     true,
		Opcode:  opcode,
		Payload: data,
		Masked:  masked,
		MaskKey: maskKey,
	}
	return ws.WriteFrame(frame)
}

func (ws *WebSocket) WriteText() error {

}

func (ws *WebSocket) WriteBinary() error {

}

func (ws *WebSocket) WritePing(data []byte) error {
	if len(data) > 125 {
		return fmt.Errorf("ping payload cannot exceed 125 bytes")
	}

	masked := ws.isClientConn()
	var maskKey [4]byte
	if masked {
		maskKey = [4]byte{0x11, 0x22, 0x33, 0x44} // Placeholder: Replace with actual random
	}

	frame := &Frame{
		Fin:     true,
		Opcode:  OpcodePing,
		Payload: data,
		Masked:  masked,
		MaskKey: maskKey,
	}

	return ws.WriteFrame(frame)

}

func (ws *WebSocket) WritePong(data []byte) error {
	if len(data) > 125 {
		return fmt.Errorf("pong payload cannot exceed 125 bytes")
	}

	masked := ws.isClientConn()
	var maskKey [4]byte
	if masked {
		maskKey = [4]byte{0x11, 0x22, 0x33, 0x44} // Placeholder: Replace with actual random
	}

	frame := &Frame{
		Fin:     true,
		Opcode:  OpcodePong,
		Payload: data,
		Masked:  masked,
		MaskKey: maskKey,
	}

	return ws.WriteFrame(frame)

}
