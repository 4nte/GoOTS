package network

import (
	"net"
	"sync"

	"gitlab.com/wick3dr0se/goots/game"
)

// ConnectionManager is used to manage TibiaConnections
type ConnectionManager struct {
	connections map[net.Conn]*TibiaConnection
	lock        *sync.Mutex
}

// NewConnectionManager ...
func NewConnectionManager() ConnectionManager {
	return ConnectionManager{
		connections: make(map[net.Conn]*TibiaConnection),
		lock:        &sync.Mutex{},
	}
}

// Add a new TibiaConnection to ConnectionManager
func (cm *ConnectionManager) Add(tc *TibiaConnection) {
	cm.lock.Lock()
	cm.connections[tc.Connection] = tc
	cm.lock.Unlock()
}

// ByConnection retrieves TibiaConnection by net.Conn
func (cm *ConnectionManager) ByConnection(conn net.Conn) *TibiaConnection {
	return cm.connections[conn]
}

// Delete a TibiaConnection from ConnectionManager
func (cm *ConnectionManager) Delete(tc *TibiaConnection) {
	cm.lock.Lock()
	delete(cm.connections, tc.Connection)
	cm.lock.Unlock()
}

// TibiaConnection ...
type TibiaConnection struct {
	Connection net.Conn
	Player     *game.Creature
	Map        *game.Map
}
