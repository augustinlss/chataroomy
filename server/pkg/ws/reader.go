package ws

import (
	"encoding/binary"
	"fmt"
	"io"
)

// first byte of a frame contains important information about it.
// MSB is FIN, 4 LSB are the opcode. the remaining 3 bits are extension bits and are generally
// 0 unless an extension is negotiated
func (ws *WebSocket) ReadFrame() (*Frame, error) {
	firstByte, err := ws.reader.ReadByte()

	if err != nil {
		return nil, fmt.Errorf("Unable to read first byte of frame: %w", err)
	}

	// use bit masking to determine fin and opcode.
	// so 0x80 => 10000000 so isolates the first bit of the byte
	// similarly 0x0F => 00001111 so isolates the last 4 bits of the byte,
	// this way we are able to get the information we need
	fin := firstByte&0x80 != 0
	opcode := firstByte & 0x0F

	secondByte, err := ws.reader.ReadByte()

	if err != nil {
		return nil, fmt.Errorf("Unable to read second byte of frame: %w", err)
	}

	masked := (secondByte & 0x80) != 0
	payloadLen := uint64(secondByte & 0x7F)

	if payloadLen == 126 {
		// Next 2 bytes are the actual length
		var extLen uint16
		if err := binary.Read(ws.reader, binary.BigEndian, &extLen); err != nil {
			return nil, fmt.Errorf("failed to read 16-bit payload length: %w", err)
		}
		payloadLen = uint64(extLen)
	} else if payloadLen == 127 {
		// Next 8 bytes are the actual length
		if err := binary.Read(ws.reader, binary.BigEndian, &payloadLen); err != nil {
			return nil, fmt.Errorf("failed to read 64-bit payload length: %w", err)
		}
	}

	var maskingKey [4]byte
	if masked {
		if _, err := io.ReadFull(ws.reader, maskingKey[:]); err != nil {
			return nil, fmt.Errorf("failed to read masking key: %w", err)
		}
	}

	payload := make([]byte, payloadLen)
	if payloadLen > 0 {
		if _, err := io.ReadFull(ws.reader, payload); err != nil {
			return nil, fmt.Errorf("failed to read payload: %w", err)
		}

		// Unmask payload if masked
		if masked {
			for i := range payload {
				payload[i] ^= maskingKey[i%4]
			}
		}
	}

	return &Frame{
		Fin:     fin,
		Opcode:  opcode,
		Masked:  masked,
		Payload: payload,
	}, nil

}

func (ws *WebSocket) ReadMessage() (opcode byte, data []byte, err error) {
	var message []byte
	var messageOpcode byte

	for {
		frame, err := ws.ReadFrame()

		if err != nil {
			return 0, nil, err
		}

		if len(message) == 0 {
			messageOpcode = frame.Opcode
		}

		message = append(message, frame.Payload...)

		if frame.Fin {
			break
		}
	}

	return messageOpcode, message, nil
}
