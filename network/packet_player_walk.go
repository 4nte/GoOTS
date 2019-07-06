package network

import "github.com/rwxsu/goot/game"

type PlayerWalkPacket struct {
	Direction uint8
}

func (packet PlayerWalkPacket) Parse(msg *Message, code uint8, tc *TibiaConnection) {
	var offset game.Offset
	var width, height uint16
	from := tc.Player.Position
	to := tc.Player.Position
	switch packet.Direction {
	case game.North:
		offset.X = -8
		offset.Y = -6
		width = 18
		height = 1
		to.Y--
	case game.South:
		offset.X = -8
		offset.Y = 7
		width = 18
		height = 1
		to.Y++
	case game.East:
		offset.X = 9
		offset.Y = -6
		width = 1
		height = 14
		to.X++
	case game.West:
		offset.X = -8
		offset.Y = -6
		width = 1
		height = 14
		to.X--
	}

	if !tc.Map.MoveCreature(tc.Player, to, packet.Direction) {
		SendSnapback(tc)
	}

	output := NewMessage()
	output.WriteUint8(0x6d)
	AddPosition(output, from)
	output.WriteUint8(0x01) // oldStackPos
	AddPosition(output, to)
	output.WriteUint8(code)
	AddMapArea(output, tc.Map, to, offset, width, height)
	SendMessage(tc.Connection, output)
}

func (packet PlayerWalkPacket) Validate(msg *Message, tc *TibiaConnection) bool {
	return true
}
