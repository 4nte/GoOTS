package network

import (
	"fmt"
	"net"

	"github.com/rwxsu/goot/game"
)

func SendInvalidClientVersion(c net.Conn) {
	msg := NewMessage()
	msg.WriteUint8(0x0a)
	msg.WriteString("Only protocol 7.40 allowed!")
	SendMessage(c, msg)
}

func SendCancelMessage(tc *TibiaConnection, str string) {
	msg := NewMessage()
	AddPlayerMessage(msg, str, game.PlayerMessageTypeCancel)
	SendMessage(tc.Connection, msg)
}

func SendSnapback(tc *TibiaConnection) {
	msg := NewMessage()
	msg.WriteUint8(0xb5)
	msg.WriteUint8(tc.Player.Direction)
	SendMessage(tc.Connection, msg)
	SendCancelMessage(tc, "Sorry, not possible.")
}

func SendAddCreature(tc *TibiaConnection) {
	res := NewMessage()
	res.WriteUint8(0x0a)
	res.WriteUint32(tc.Player.ID) // ID
	res.WriteUint16(0x32)         // ?
	// can report bugs?
	if tc.Player.Access > game.Regular {
		res.WriteUint8(0x01)
	} else {
		res.WriteUint8(0x00)
	}
	if tc.Player.Access >= game.Gamemaster {
		res.WriteUint8(0x0b)
		for i := 0; i < 32; i++ {
			res.WriteUint8(0xff)
		}
	}
	tile := tc.Map.GetTile(tc.Player.Position)
	tile.AddCreature(tc.Player)
	res.WriteUint8(0x64)
	AddPosition(res, tc.Player.Position)
	AddMapArea(res, tc.Map, tc.Player.Position, game.Offset{X: -8, Y: -6, Z: 0}, 18, 14)
	AddMagicEffect(res, tc.Player.Position, 0x0a)
	AddInventory(res, tc.Player)
	AddStats(res, tc.Player)
	AddSkills(res, tc.Player)
	AddWorldLight(res, tc.Player.World)
	AddCreatureLight(res, tc.Player)
	AddPlayerMessage(res, fmt.Sprintf("Welcome, %s.", tc.Player.Name), game.PlayerMessageTypeInfo)
	AddPlayerMessage(res, "TODO: Last Login String 01-01-1970", game.PlayerMessageTypeInfo)
	AddCreatureLight(res, tc.Player)
	AddIcons(res, tc.Player)
	SendMessage(tc.Connection, res)
}

func AddCreatureLight(msg *Message, c *game.Creature) {
	msg.WriteUint8(0x8d)
	msg.WriteUint32(c.ID)
	msg.WriteUint8(c.Light.Level)
	msg.WriteUint8(c.Light.Color)
}

func AddWorldLight(msg *Message, w game.World) {
	msg.WriteUint8(0x82)
	msg.WriteUint8(w.Light.Level)
	msg.WriteUint8(w.Light.Color)
}

func AddIcons(msg *Message, c *game.Creature) {
	msg.WriteUint8(0xa2)
	msg.WriteUint8(c.Icons)
}

func AddSkills(msg *Message, c *game.Creature) {
	msg.WriteUint8(0xa1)
	msg.WriteUint8(c.Fist.Level)
	msg.WriteUint8(c.Fist.Percent)
	msg.WriteUint8(c.Club.Level)
	msg.WriteUint8(c.Club.Percent)
	msg.WriteUint8(c.Sword.Level)
	msg.WriteUint8(c.Sword.Percent)
	msg.WriteUint8(c.Axe.Level)
	msg.WriteUint8(c.Axe.Percent)
	msg.WriteUint8(c.Distance.Level)
	msg.WriteUint8(c.Distance.Percent)
	msg.WriteUint8(c.Shielding.Level)
	msg.WriteUint8(c.Shielding.Percent)
	msg.WriteUint8(c.Fishing.Level)
	msg.WriteUint8(c.Fishing.Percent)
}

func AddStats(msg *Message, c *game.Creature) {
	msg.WriteUint8(0xa0) // send player stats
	msg.WriteUint16(c.HealthNow)
	msg.WriteUint16(c.HealthMax)
	msg.WriteUint16(c.Cap)
	msg.WriteUint32(c.Combat.Experience)
	msg.WriteUint8(c.Combat.Level)
	msg.WriteUint8(c.Combat.Percent)
	msg.WriteUint16(c.ManaNow)
	msg.WriteUint16(c.ManaMax)
	msg.WriteUint8(c.Magic.Level)
	msg.WriteUint8(c.Magic.Percent)
}

func AddInventory(msg *Message, c *game.Creature) {
	msg.WriteUint8(game.SlotEmpty)
	msg.WriteUint8(game.SlotHead)

	msg.WriteUint8(game.SlotEmpty)
	msg.WriteUint8(game.SlotNecklace)

	msg.WriteUint8(game.SlotNotEmpty)
	msg.WriteUint8(game.SlotBackpack)
	msg.WriteUint16(0x7c4) // backpack

	msg.WriteUint8(game.SlotNotEmpty)
	msg.WriteUint8(game.SlotBody)
	msg.WriteUint16(0x9a8) // magic plate armor

	msg.WriteUint8(game.SlotEmpty)
	msg.WriteUint8(game.SlotShield)

	msg.WriteUint8(game.SlotNotEmpty)
	msg.WriteUint8(game.SlotWeapon)
	msg.WriteUint16(0x997) // crossbow

	msg.WriteUint8(game.SlotEmpty)
	msg.WriteUint8(game.SlotLegs)

	msg.WriteUint8(game.SlotEmpty)
	msg.WriteUint8(game.SlotFeet)

	msg.WriteUint8(game.SlotEmpty)
	msg.WriteUint8(game.SlotRing)

	msg.WriteUint8(game.SlotNotEmpty)
	msg.WriteUint8(game.SlotAmmo)
	msg.WriteUint16(0x9ef) // bolts
	msg.WriteUint8(33)     // count
}

// AddMapArea adds the area starting at position+offset until width and height
// is reached. Returns the number of tiles (counting nil)
func AddMapArea(msg *Message, m *game.Map, pos game.Position, offset game.Offset, width, height uint16) int {
	pos.Offset(offset)
	count := 0
	skip := (uint8)(0)
	if pos.Z < 8 {
		for z := (int8)(7); z > -1; z-- {
			for x := (uint16)(0); x < width; x++ {
				for y := (uint16)(0); y < height; y++ {
					tile := m.GetTile(game.Position{X: pos.X + x, Y: pos.Y + y, Z: (uint8)(z)})
					if tile != nil {
						if skip > 0 {
							msg.WriteUint8(skip - 1)
							msg.WriteUint8(0xff)
							skip = 0
						}
						AddTile(msg, tile)
					}
					skip++
					if skip == 0xff {
						msg.WriteUint8(0xff)
						msg.WriteUint8(0xff)
						skip = 0
					}
					count++
				}
			}
		}
	} else { // TODO: underground

	}
	// Remainder
	if skip > 0 {
		msg.WriteUint8(skip - 1)
		msg.WriteUint8(0xff)
		skip = 0
	}
	return count
}

func AddPosition(msg *Message, pos game.Position) {
	msg.WriteUint16(pos.X)
	msg.WriteUint16(pos.Y)
	msg.WriteUint8(pos.Z)
}

func AddMagicEffect(msg *Message, pos game.Position, kind uint8) {
	msg.WriteUint8(0x83)
	AddPosition(msg, pos)
	msg.WriteUint8(kind)
}

func AddCreature(msg *Message, c *game.Creature) {
	msg.WriteUint16(0x61) // unknown creature
	msg.WriteUint32(0x00) // something about caching known creatures
	msg.WriteUint32(c.ID)
	msg.WriteString(c.Name)
	msg.WriteUint8((uint8)(c.HealthNow*100/c.HealthMax) + 1)
	msg.WriteUint8(c.Direction)
	msg.WriteUint8(c.Outfit.Type)
	msg.WriteUint8(c.Outfit.Head)
	msg.WriteUint8(c.Outfit.Body)
	msg.WriteUint8(c.Outfit.Legs)
	msg.WriteUint8(c.Outfit.Feet)
	msg.WriteUint8(c.Light.Level)
	msg.WriteUint8(c.Light.Color)
	msg.WriteUint16(c.Speed)
	msg.WriteUint8(c.Skull)
	msg.WriteUint8(c.Party)
}

// AddTile adds all the tile items and creatures WITHOUT the end of tile
// delimeter (0xSKIP-0xff)
func AddTile(msg *Message, tile *game.Tile) {
	for _, i := range tile.Items {
		msg.WriteUint16(i.ID)
	}
	for _, c := range tile.Creatures {
		AddCreature(msg, c)
	}
}

func AddPlayerMessage(msg *Message, str string, kind uint8) {
	msg.WriteUint8(0xb4)
	msg.WriteUint8(kind)
	msg.WriteString(str)
}

// Placeholder player
func GetDumpPlayer() game.Creature {
	return game.Creature{
		ID:        0x04030201,
		Access:    game.Tutor,
		Name:      "rwxsu",
		Cap:       50,
		Combat:    game.Skill{Level: 8, Percent: 0, Experience: 4200},
		HealthNow: 100,
		HealthMax: 200,
		ManaNow:   50,
		ManaMax:   100,
		Magic:     game.Skill{Level: 10, Percent: 50},
		Fist:      game.Skill{Level: 10, Percent: 50},
		Club:      game.Skill{Level: 10, Percent: 50},
		Sword:     game.Skill{Level: 10, Percent: 50},
		Axe:       game.Skill{Level: 10, Percent: 50},
		Distance:  game.Skill{Level: 10, Percent: 50},
		Shielding: game.Skill{Level: 10, Percent: 50},
		Fishing:   game.Skill{Level: 10, Percent: 50},
		Direction: game.South,
		Position:  game.Position{X: 32000, Y: 32000, Z: 7},
		Outfit: game.Outfit{
			Type: 0x80,
			Head: 0x50,
			Body: 0x50,
			Legs: 0x50,
			Feet: 0x50,
		},
		Skull: 3,
		Icons: 1,
		Light: game.Light{Level: 0x7, Color: 0xd7},
		World: game.World{Light: game.Light{Level: 0x00, Color: 0xd7}},
		Speed: 60000,
	}
}
