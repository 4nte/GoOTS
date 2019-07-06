package game

import (
	"fmt"
)

// Direction
const (
	North uint8 = iota
	East
	South
	West
	NorthEast
	SouthEast
	SouthWest
	NorthWest
)

// Position is the real in-game position
type Position struct {
	X uint16
	Y uint16
	Z uint8
}

type Offset struct {
	X, Y, Z int8
}

func (pos *Position) Offset(offset Offset) {
	pos.X += (uint16)(offset.X)
	pos.Y += (uint16)(offset.Y)
	pos.Z += (uint8)(offset.Z)
}

// Equals checks if two positions are the same.
func (pos *Position) Equals(other Position) bool {
	return pos.X == other.X && pos.Y == other.Y && pos.Z == other.Z
}

func (pos *Position) String() string {
	return fmt.Sprintf("X:%d, Y:%d, Z:%d", pos.X, pos.Y, pos.Z)
}

// Light has the same structure for both creatures and world
type Light struct {
	Level uint8
	Color uint8
}
