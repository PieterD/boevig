package sshserver

import (
	"encoding/binary"
	"fmt"
)

type WindowChangePayload struct {
	Columns     uint32
	Rows        uint32
	PixelWidth  uint32
	PixelHeight uint32
}

func parseWindowChangePayload(raw []byte) (p WindowChangePayload, err error) {
	p.Columns = binary.BigEndian.Uint32(raw)
	raw = raw[4:]
	p.Rows = binary.BigEndian.Uint32(raw)
	raw = raw[4:]
	p.PixelWidth = binary.BigEndian.Uint32(raw)
	raw = raw[4:]
	p.PixelHeight = binary.BigEndian.Uint32(raw)
	raw = raw[4:]
	if len(raw) != 0 {
		return WindowChangePayload{}, fmt.Errorf("garbage bytes at the end: %d", len(raw))
	}
	return p, nil
}

type PtyReqPayload struct {
	Term        string
	Columns     uint32
	Rows        uint32
	PixelWidth  uint32
	PixelHeight uint32
	Modes       []TerminalMode
}

type TerminalMode struct {
	OpCode   byte
	Argument uint32
}

func parsePtyReqPayload(raw []byte) (p PtyReqPayload, err error) {
	stringLen := binary.BigEndian.Uint32(raw)
	raw = raw[4:]
	p.Term = string(raw[:stringLen])
	raw = raw[stringLen:]
	p.Columns = binary.BigEndian.Uint32(raw)
	raw = raw[4:]
	p.Rows = binary.BigEndian.Uint32(raw)
	raw = raw[4:]
	p.PixelWidth = binary.BigEndian.Uint32(raw)
	raw = raw[4:]
	p.PixelHeight = binary.BigEndian.Uint32(raw)
	raw = raw[4:]
	modeLength := binary.BigEndian.Uint32(raw)
	raw = raw[4:]
	if modeLength != uint32(len(raw)) {
		return PtyReqPayload{}, fmt.Errorf("mode length mismatch: %d bytes left in buffer, yet mode length is %d", len(raw), modeLength)
	}
	for len(raw) >= 5 {
		opcode := raw[0]
		raw = raw[1:]
		argument := binary.BigEndian.Uint32(raw)
		raw = raw[4:]
		p.Modes = append(p.Modes, TerminalMode{
			OpCode:   opcode,
			Argument: argument,
		})
	}
	if len(raw) == 1 && raw[0] == 0 {
		return p, nil
	}
	return PtyReqPayload{}, fmt.Errorf("invalid tail: expected 1, got %d bytes %#v", len(raw), raw)
}
