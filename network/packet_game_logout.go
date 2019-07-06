package network

type GameLogoutPacket struct {
}

func (packet GameLogoutPacket) Parse(msg *Message, code uint8, tc *TibiaConnection) {
	tc.Map.GetTile(tc.Player.Position).RemoveCreature(tc.Player)
}

func (packet GameLogoutPacket) Validate(msg *Message, tc *TibiaConnection) bool {
	return true
}
