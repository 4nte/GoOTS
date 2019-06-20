package network

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type Message struct {
	Buffer []uint8
	Cursor uint16
}

func NewMessage() *Message {
	var p Message
	p.Buffer = make([]uint8, 2)
	p.Cursor = 2
	return &p
}

func (p *Message) overflow() bool {
	return p.Cursor >= (uint16)(len(p.Buffer))
}

func (p *Message) ReadUint8() uint8 {
	if p.overflow() {
		return 0
	}
	v := p.Buffer[p.Cursor]
	p.Cursor++
	return v
}

func (p *Message) ReadUint16() uint16 {
	if p.overflow() {
		return 0
	}
	v := binary.LittleEndian.Uint16(p.Buffer[p.Cursor : p.Cursor+2])
	p.Cursor += 2
	return v
}

func (p *Message) ReadUint32() uint32 {
	if p.overflow() {
		return 0
	}
	v := binary.LittleEndian.Uint32(p.Buffer[p.Cursor : p.Cursor+4])
	p.Cursor += 4
	return v
}

func (p *Message) ReadString() string {
	if p.overflow() {
		return ""
	}
	var str string
	strlen := p.ReadUint16()
	for i := (uint16)(0); i < strlen; i++ {
		str += (string)(p.ReadUint8())
	}
	return str
}

func (p *Message) WriteUint8(v uint8) {
	p.Buffer = append(p.Buffer, v)
	binary.LittleEndian.PutUint16(p.Buffer[0:2], (uint16)(len(p.Buffer)-2))
	p.Cursor++
}

func (p *Message) WriteUint16(v uint16) {
	bytes := make([]uint8, 2)
	binary.LittleEndian.PutUint16(bytes, v)
	p.WriteUint8(bytes[0])
	p.WriteUint8(bytes[1])
}

func (p *Message) WriteUint32(v uint32) {
	bytes := make([]uint8, 4)
	binary.LittleEndian.PutUint32(bytes, v)
	p.WriteUint8(bytes[0])
	p.WriteUint8(bytes[1])
	p.WriteUint8(bytes[2])
	p.WriteUint8(bytes[3])
}

func (p *Message) WriteString(str string) {
	p.WriteUint16((uint16)(len(str)))
	for i := 0; i < len(str); i++ {
		p.WriteUint8((uint8)(str[i]))
	}
}

func (p *Message) Length() uint16 {
	return binary.LittleEndian.Uint16(p.Buffer[0:2])
}

func (p *Message) SkipBytes(n uint16) {
	p.Cursor += n
}

func (p *Message) HexDump(prefix string) {
	fmt.Printf("\n[%s]\n%s", prefix, hex.Dump(p.Buffer))
}
