// Package models implements the backend operation
package models

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type clientInfor struct {
	Users map[string]WebConn //string here could be a user name , ip address
}

// WebConn define user websocket conn
type WebConn struct {
	Device string
	*websocket.Conn
}

// Clients record uses online
var Clients clientInfor

func init() {
	Clients.Users = make(map[string]WebConn)
}

// Add implements add an client
func (cs *clientInfor) Add(remoteAddr string, conn *websocket.Conn) {
	if _, ok := cs.Users[remoteAddr]; !ok {
		var c WebConn
		c.Conn = conn
		cs.Users[remoteAddr] = c
	}
}

// Delete implements delete and client
func (cs *clientInfor) Delete(remoteAddr string) {
	c, ok := cs.Users[remoteAddr]
	if ok {
		delete(cs.Users, remoteAddr)
	}
	c.Close()
}

// Broadcast implements broadcast information
func (cs *clientInfor) Broadcast(exChangeData ExChange) {
	exChangeJSON, _ := json.Marshal(exChangeData)
	source := exChangeData.Source
	if source != EVERYone && source != "" {
		conn := Clients.Users[source]
		conn.WriteJSON(string(exChangeJSON))
	} else {
		for _, conn := range Clients.Users {
			conn.WriteJSON(string(exChangeJSON))
		}
	}
}
