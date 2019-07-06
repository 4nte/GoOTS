package network

type GamePingPacket struct {
}

func (packet GamePingPacket) Parse(msg *Message, code uint8, tc *TibiaConnection) {
	output := NewMessage()
	output.WriteUint8(code)
	SendMessage(tc.Connection, output)
}

func (packet GamePingPacket) Validate(msg *Message, tc *TibiaConnection) bool {
	return true
}
