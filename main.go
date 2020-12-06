package main

import (
	"log"
	"net"
	"path/filepath"

	"gitlab.com/wick3dr0se/goots/game"
	"gitlab.com/wick3dr0se/goots/network"
)

var connectionManager network.ConnectionManager

func main() {
	const sectors = "data/map/sectors/*"
	filenames, _ := filepath.Glob(sectors)

	m := make(game.Map)

	for _, fn := range filenames {
		m.LoadSector(fn)
	}

	connectionManager = network.NewConnectionManager()

	l, err := net.Listen("tcp", ":7171")
	if err != nil {
		panic(err)
	}
	defer l.Close()
	acceptConnections(l, &m)
}

func acceptConnections(l net.Listener, m *game.Map) {
	for {
		c, err := l.Accept()
		if err != nil {
			log.Println(err)
			return
		}
		go handleConnection(c, m)
	}
}

func handleConnection(c net.Conn, m *game.Map) {
	var tc *network.TibiaConnection
	if tc = connectionManager.ByConnection(c); tc == nil {
		tc = &network.TibiaConnection{
			Connection: c,
		}
	}
	// Placeholder player
	player := game.Creature{
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
connectionLoop:
	for {
		req := network.RecvMessage(tc.Connection)
		if req == nil {
			return
		}
		code := req.ReadUint8()
		switch code {
		case 0x01: // request character list
			req.SkipBytes(2) // os := req.ReadUint16()
			if req.ReadUint16() != 740 {
				network.SendInvalidClientVersion(tc.Connection)
				break connectionLoop
			}
			network.SendCharacterList(tc.Connection)
			break connectionLoop
		case 0x0a: // request character login
			req.SkipBytes(2) // os := req.ReadUint16()
			if req.ReadUint16() != 740 {
				network.SendInvalidClientVersion(tc.Connection)
				break connectionLoop
			}
			tc.Map = m
			tc.Player = &player
			connectionManager.Add(tc)
			network.SendAddCreature(tc)
		case 0x14: // logout
			tc.Map.GetTile(tc.Player.Position).RemoveCreature(tc.Player)
			connectionManager.Delete(tc)
			break connectionLoop
		default:
			network.ParseCommand(tc, req, code)
		}
	}
	if err := c.Close(); err != nil {
		log.Printf("Unable to close connection %v\n", err)
	}
}
