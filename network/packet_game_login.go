package network

type GameLoginPacket struct {
}

func (packet GameLoginPacket) Parse(msg *Message, code uint8, tc *TibiaConnection) {
	SendAddCreature(tc)
}

func (packet GameLoginPacket) Validate(msg *Message, tc *TibiaConnection) bool {
	msg.SkipBytes(2) // os := msg.ReadUint16()
	if msg.ReadUint16() != 740 {
		SendInvalidClientVersion(tc.Connection)
		return false
	}

	return true
}
