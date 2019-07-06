package network

import "github.com/rwxsu/goot/game"

func LoadKnownPackets() map[uint8]TibiaPacket {
	knownPackets := make(map[uint8]TibiaPacket, 10)

	// charlist packet
	knownPackets[0x01] = GameCharListPacket{}

	// login packet
	knownPackets[0x0a] = GameLoginPacket{}

	// logout packet
	knownPackets[0x14] = GameLogoutPacket{}

	// walk packets
	knownPackets[0x65] = PlayerWalkPacket{Direction: game.North}
	knownPackets[0x66] = PlayerWalkPacket{Direction: game.East}
	knownPackets[0x67] = PlayerWalkPacket{Direction: game.South}
	knownPackets[0x68] = PlayerWalkPacket{Direction: game.West}

	// diagonal walk packets
	knownPackets[0x6A] = PlayerWalkPacket{Direction: game.NorthEast}
	knownPackets[0x6B] = PlayerWalkPacket{Direction: game.SouthEast}
	knownPackets[0x6C] = PlayerWalkPacket{Direction: game.SouthWest}
	knownPackets[0x6D] = PlayerWalkPacket{Direction: game.NorthWest}

	// turn packets
	knownPackets[0x6f] = PlayerTurnPacket{Direction: game.North}
	knownPackets[0x70] = PlayerTurnPacket{Direction: game.East}
	knownPackets[0x71] = PlayerTurnPacket{Direction: game.South}
	knownPackets[0x72] = PlayerTurnPacket{Direction: game.West}

	return knownPackets
}

type TibiaPacket interface {
	Parse(msg *Message, code uint8, tc *TibiaConnection)
	Validate(msg *Message, tc *TibiaConnection) bool
}
