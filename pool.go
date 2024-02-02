package main

import (
	"database/sql"
	"sync"
	"time"
)

type DBConnection struct {
	Db           *sql.DB
	usedLastTime int
	Status       string
}

type DBConnectionPool struct {
	activeConnections    int
	maxActiveConnections int
	userToConnectionMap  map[string][]*DBConnection
	idleTimeout          int
	mu                   *sync.Mutex
}

func NewDBConnection() *DBConnection {
	return &DBConnection{
		Db: new(sql.DB),
	}
}

func NewDBConnectionPool() *DBConnectionPool {
	newConnectionPool := new(DBConnectionPool)
	go func() {
		for {
			newConnectionPool.closeIdleConnections()
			time.Sleep(1 * time.Minute)
		}
	}()
	return newConnectionPool
}

func (db *DBConnectionPool) GetConnection(userId string) *DBConnection {
	if val, ok := db.getIdleConnections(userId); ok {
		return val
	}
	newConnection := new(DBConnection)
	for newConnection == nil {
		db.mu.Lock()
		if db.maxActiveConnections > db.activeConnections {
			db.activeConnections++
			newConnection = NewDBConnection()
			db.userToConnectionMap[userId] = append(db.userToConnectionMap[userId], newConnection)
		}
		db.mu.Unlock()
	}
	return newConnection
}

func (db *DBConnectionPool) getIdleConnections(userId string) (*DBConnection, bool) {
	if connections, ok := db.userToConnectionMap[userId]; ok {
		for _, connection := range connections {
			db.mu.Lock()
			if connection.Status == "Idle" {
				connection.Status = "Allocated"
				db.mu.Unlock()
				return connection, true
			}
			db.mu.Unlock()
		}
	}
	return nil, false
}

func (db *DBConnectionPool) closeIdleConnections() {
	idleConnections := make([]*DBConnection, 0)
	for userId, connections := range db.userToConnectionMap {
		for index, connection := range connections {
			db.mu.Lock()
			if connection.usedLastTime >= db.idleTimeout {
				idleConnections = append(idleConnections, connection)
				connection.Status = "Expired"
				db.removeConnection(userId, index)
			}
			db.mu.Unlock()
		}
	}
	for _, connection := range idleConnections {
		connection.Db.Close()
	}
}

func (db *DBConnectionPool) removeConnection(userId string, index int) {
	copy(db.userToConnectionMap[userId][index:], db.userToConnectionMap[userId][index+1:])
}
