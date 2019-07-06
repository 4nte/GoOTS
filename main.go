package main

import (
	"log"
	"net"
	"path/filepath"
	"reflect"

	"github.com/rwxsu/goot/game"
	"github.com/rwxsu/goot/network"
)

var knowPackets map[uint8]network.TibiaPacket
var connectionManager network.ConnectionManager

func main() {
	const sectors = "data/map/sectors/*"
	filenames, _ := filepath.Glob(sectors)

	m := make(game.Map)

	for _, fn := range filenames {
		m.LoadSector(fn)
	}

	connectionManager = network.NewConnectionManager()
	knowPackets = network.LoadKnownPackets()

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

connectionLoop:
	for {
		req := network.RecvMessage(tc.Connection)
		if req == nil {
			return
		}

		code := req.ReadUint8()
		if packet, known := knowPackets[code]; known {
			log.Printf(">>> RECV %s", reflect.TypeOf(packet).Name())
			// manual handling: add player and map on login packet
			if code == 0x0a {
				player := network.GetDumpPlayer()
				tc.Map = m
				tc.Player = &player
				connectionManager.Add(tc)
			}

			if !packet.Validate(req, tc) {
				log.Printf(">>> REJECT %s", reflect.TypeOf(packet).Name())
				break connectionLoop
			}

			// automatic handling
			log.Printf(">>> ACCEPT %s", reflect.TypeOf(packet).Name())
			packet.Parse(req, code, tc)

			// manual handling: remove tc on logout packet
			if code == 0x14 {
				connectionManager.Delete(tc)
				break connectionLoop
			}
		}
	}
	if err := c.Close(); err != nil {
		log.Printf("Unable to close connection %v\n", err)
	}
}
