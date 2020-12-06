package network

import (
	"net"

	"gitlab.com/wick3dr0se/goots/game"
)

func SendInvalidClientVersion(c net.Conn) {
	msg := NewMessage()
	msg.WriteUint8(0x0a)
	msg.WriteString("Only protocol 7.40 allowed!")
	SendMessage(c, msg)
}

func SendCharacterList(c net.Conn) {
	characters := make([]game.Creature, 2)
	characters[0].Name = "admin"
	characters[0].World.Name = "test"
	characters[0].World.Port = 7171
	characters[1].Name = "rwxsu"
	characters[1].World.Name = "test"
	characters[1].World.Port = 7171
	res := NewMessage()
	res.WriteUint8(0x14) // MOTD
	res.WriteString("Welcome to GoOT.")
	res.WriteUint8(0x64) // character list
	res.WriteUint8((uint8)(len(characters)))
	for i := 0; i < len(characters); i++ {
		res.WriteString(characters[i].Name)
		res.WriteString(characters[i].World.Name)
		res.WriteUint8(127)
		res.WriteUint8(0)
		res.WriteUint8(0)
		res.WriteUint8(1)
		res.WriteUint16(characters[i].World.Port)
	}
	res.WriteUint16(0) // premium days
	SendMessage(c, res)
}
