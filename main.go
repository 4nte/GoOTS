package main

import "fmt"
import "net"
import "github.com/rwxsu/goot/netmsg"
import "github.com/rwxsu/goot/packet"

func main() {
	l, err := net.Listen("tcp", ":7171")
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go func(c net.Conn) {
			msg := netmsg.New(&c)
			reqCode := OnRecvHeader(msg)
			switch reqCode {
			case packet.REQUEST_CHARACTER_LOGIN:
				OnRequestCharacterLogin(msg)
				break
			case packet.REQUEST_CHARACTER_LIST:
				OnRequestCharacterList(msg)
				break
			}
		}(conn)
	}
}

func OnRecvHeader(msg *netmsg.NetMsg) uint8 {
	length := msg.ReadUint16()
	fmt.Println("\nlogin.packet.len:", length)

	reqCode := msg.ReadUint8()

	// os := msg.ReadUint16()
	msg.SkipBytes(2)

	protocolVersion := msg.ReadUint16()
	if protocolVersion != 740 {
		packet.SendDisconnect(msg, "Only protocol 7.40 allowed!")
		fmt.Println("login.protocol: invalid")
		return 0
	}
	fmt.Println("login.protocol:", protocolVersion)
	return reqCode
}

func OnRequestCharacterList(msg *netmsg.NetMsg) {
	fmt.Println("[REQUEST_CHARACTER_LIST]")

	msg.SkipBytes(12)
	// msg.ReadUint32() // Tibia.spr version
	// msg.ReadUint32() // Tibia.dat version
	// msg.ReadUint32() // Tibia.pic version

	acc := msg.ReadUint32()
	pwd := msg.ReadString()
	fmt.Println("login.acc:", acc)
	fmt.Println("login.pwd:", pwd)

	// TODO: Authenticate
	// packet.SendCharacterList(msg, characters)
	packet.SendCharacterList(msg)
}

func OnRequestCharacterLogin(msg *netmsg.NetMsg) {
	fmt.Println("[REQUEST_CHARACTER_LOGIN]")
	msg.SkipBytes(1) // ???

	acc := msg.ReadUint32()
	name := msg.ReadString()
	pwd := msg.ReadString()
	fmt.Printf("login.acc: %d\n", acc)
	fmt.Printf("login.pwd: %s\n", pwd)
	fmt.Printf("login.character.name: %s\n", name)

	// TODO: Authenticate
	// packet.SendCharacterLogin(msg, character)
}
