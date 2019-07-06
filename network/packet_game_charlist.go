package network

import "github.com/rwxsu/goot/game"

type GameCharListPacket struct{}

func (packet GameCharListPacket) Parse(msg *Message, code uint8, tc *TibiaConnection) {
	characters := make([]game.Creature, 2)
	characters[0].Name = "admin"
	characters[0].World.Name = "test"
	characters[0].World.Port = 7171
	characters[1].Name = "rwxsu"
	characters[1].World.Name = "test"
	characters[1].World.Port = 7171
	output := NewMessage()
	output.WriteUint8(0x14) // MOTD
	output.WriteString("Welcome to GoOT.")
	output.WriteUint8(0x64) // character list
	output.WriteUint8((uint8)(len(characters)))
	for i := 0; i < len(characters); i++ {
		output.WriteString(characters[i].Name)
		output.WriteString(characters[i].World.Name)
		output.WriteUint8(127)
		output.WriteUint8(0)
		output.WriteUint8(0)
		output.WriteUint8(1)
		output.WriteUint16(characters[i].World.Port)
	}
	output.WriteUint16(0) // premium days
	SendMessage(tc.Connection, output)
}

func (packet GameCharListPacket) Validate(msg *Message, tc *TibiaConnection) bool {
	msg.SkipBytes(2) // os := msg.ReadUint16()
	if msg.ReadUint16() != 740 {
		SendInvalidClientVersion(tc.Connection)
		return false
	}

	return true
}
