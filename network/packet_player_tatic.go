package network

type PlayerTaticPacket struct{}

func (packet PlayerTaticPacket) Parse(msg *Message, code uint8, tc *TibiaConnection) {
	tc.Player.Tactic.FightMode = msg.ReadUint8()
	tc.Player.Tactic.ChaseOpponent = msg.ReadUint8()
	tc.Player.Tactic.AttackPlayers = msg.ReadUint8()
}

func (packet PlayerTaticPacket) Validate(msg *Message, tc *TibiaConnection) bool {
	return true
}
