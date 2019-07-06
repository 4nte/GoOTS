package network

type PlayerTurnPacket struct {
	Direction uint8
}

func (packet PlayerTurnPacket) Parse(msg *Message, code uint8, tc *TibiaConnection) {
	tc.Player.Direction = packet.Direction
	output := NewMessage()
	output.WriteUint8(0x6b)
	AddPosition(output, tc.Player.Position)
	output.WriteUint8(1)
	output.WriteUint16(0x63)
	output.WriteUint32(tc.Player.ID)
	output.WriteUint8(tc.Player.Direction)
	SendMessage(tc.Connection, output)
}

func (packet PlayerTurnPacket) Validate(msg *Message, tc *TibiaConnection) bool {
	return true
}
