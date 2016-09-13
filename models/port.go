// Package models implements the backend operation
package models

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/tarm/serial"
)

// Operate const
const (
	QUITprogram = "quit"
	OPENport    = "open"
	CLOSEport   = "close"
	WRITEport   = "write"
	READport    = "read"
	GETdevice   = "device"
	DEFAULT     = "default"
	READonly    = "readonly"
	EVERYone    = "everyone"
)

// Status const
const (
	OK    = "ok"
	NOK   = "nok"
	OWNER = "owner"
)

// ExChange define the format of exchange infor over websocket
type ExChange struct {
	Cmd    string
	Msg    interface{}
	Target string
	Source string
}

// SerInfor record all serial information in a map named by port name
type SerInfor struct {
	Ports map[string]*SerConn
}

//SerConn record single serial information
type SerConn struct {
	Owner string
	*serial.Config
	*serial.Port
}

// Add implements add a new serial port
func (s *SerInfor) Add(portName string, owner string) error {
	config := new(serial.Config)
	config.Name = portName
	config.Baud = 115200
	config.ReadTimeout = time.Millisecond * 20
	port, err := serial.OpenPort(config)
	if err != nil {
		return err
	}
	var conn SerConn
	conn.Config = config
	conn.Port = port
	conn.Owner = owner
	s.Ports[portName] = &conn

	return err
}

// Delete implements delete a serial port
func (s *SerInfor) Delete(portName string) {
	if p, ok := s.Ports[portName]; ok {
		p.Close()
		delete(s.Ports, portName)
	}
}

// WriteStr implements write to serial port
func (s *SerConn) WriteStr(cmd string) error {
	_, err := s.Write([]byte(cmd))
	return err
}

// ReadStr implements read from serial port
func (s *SerConn) ReadStr() (string, error) {
	buf := make([]byte, 1024)
	n, err := s.Read(buf)
	if err == nil {
		return string(buf[:n]), err
	}
	return "", err
}

func (s *SerConn) setOwner(owner string) {
	s.Owner = owner
}

// Sers record the information of all the client
var Sers SerInfor

// Operate record the operation to serial
var Operate ExChange

// Lock is to make written of global var in order
var Lock sync.Mutex

func init() {
	Sers.Ports = make(map[string]*SerConn)
}

func getDevice() []string {
	cmd := exec.Command("/bin/sh", "-c", "ls /dev/tty* |grep '/dev/ttyACM\\|/dev/ttyUSB'")
	out, _ := cmd.Output()
	outStr := string(out)
	outStr = strings.TrimSpace(outStr)
	devices := strings.Split(outStr, "\n")
	return devices

}

// CloseSerial implements quit port operation loop
func CloseSerial() {
	Operate.Cmd = QUITprogram //Operate may be overwirte by websocket onmessage
}

// RunSerial implements all the serial operate
func RunSerial() {
	var exChangeData ExChange
	op := &Operate
FOR:
	for {
		switch op.Cmd {
		case QUITprogram:
			for portName := range Sers.Ports {
				Sers.Delete(portName)
			}
			break FOR
		case GETdevice:
			exChangeData.Cmd = GETdevice
			exChangeData.Msg = getDevice()
			exChangeData.Target = EVERYone
			exChangeData.Source = EVERYone
			Clients.Broadcast(exChangeData)
			op.Cmd = DEFAULT
		case OPENport:
			Lock.Lock()
			var msg = OK
			if port, ok := Sers.Ports[op.Target]; !ok {
				err := Sers.Add(op.Target, op.Source)
				if err != nil {
					msg = NOK
					fmt.Println("Failed to open port ", op.Target)
				}
			} else {
				if _, ok := Clients.Users[port.Owner]; !ok {
					port.Owner = op.Source
					msg = OWNER
				}
			}
			exChangeData.Cmd = OPENport
			exChangeData.Msg = msg
			exChangeData.Target = op.Target
			exChangeData.Source = op.Source
			Clients.Broadcast(exChangeData)
			op.Cmd = READport
			Lock.Unlock()
		case CLOSEport:
			Lock.Lock()
			Sers.Delete(op.Target)
			exChangeData.Cmd = CLOSEport
			exChangeData.Target = op.Target
			Clients.Broadcast(exChangeData)
			op.Cmd = READport
			Lock.Unlock()
		case WRITEport:
			Lock.Lock()
			var m string
			port := Sers.Ports[op.Target]
			if op.Source != port.Owner {
				m = READonly
			} else {
				if msg, ok := op.Msg.(string); ok {
					err := port.WriteStr(msg)
					m = OK
					if err != nil {
						m = NOK
						fmt.Println("Failed to write to port ", op.Target)
					}
				}
			}
			exChangeData.Cmd = WRITEport
			exChangeData.Msg = m
			exChangeData.Target = op.Target
			exChangeData.Source = op.Source

			Clients.Broadcast(exChangeData)
			op.Cmd = READport
			Lock.Unlock()
		case READport:
			for _, port := range Sers.Ports {
				data, err := port.ReadStr()
				if err != nil {
					// fmt.Println("read nothing from port ")
				} else {
					if data != "" {
						exChangeData.Cmd = READport
						exChangeData.Msg = data
						exChangeData.Target = port.Config.Name
						exChangeData.Source = EVERYone
						Clients.Broadcast(exChangeData)
					}
				}
			}
		default:
			// To do
		}
		time.Sleep(time.Millisecond * 20)
	}
}
