package network

import (
	"testing"

	"gitlab.com/wick3dr0se/goots/game"
)

func TestAddMapArea(t *testing.T) {
	msg := NewMessage()
	m := make(game.Map)
	m.InitializeSector(game.SectorPosition{X: 1, Y: 1, Z: 1}, 104)
	var n int
	n = AddMapArea(msg, &m, game.Position{X: 1, Y: 1, Z: 1}, game.Offset{X: 0, Y: 0, Z: 0}, 18, 1)
	if n != 144 {
		t.Errorf("expected 18*1*8=144, got=%d", n)
	}
	n = AddMapArea(msg, &m, game.Position{X: 1, Y: 1, Z: 1}, game.Offset{X: 0, Y: 0, Z: 0}, 1, 14)
	if n != 112 {
		t.Errorf("expected 1*14*8=112, got=%d", n)
	}
	n = AddMapArea(msg, &m, game.Position{X: 1, Y: 1, Z: 1}, game.Offset{X: 0, Y: 0, Z: 0}, 18, 14)
	if n != 2016 {
		t.Errorf("expected 18*14*8=2016, got=%d", n)
	}
}
