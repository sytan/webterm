// Package models implements the backend operation
package models

import (
	"fmt"
	"os/exec"
	"strings"
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
)

// Status const
const (
	OK  = "ok"
	NOK = "nok"
)

// ExChange define the format of exchange infor over websocket
type ExChange struct {
	Cmd    string
	Msg    interface{}
	Target string
}

// SerInfor record all serial information in a map named by port name
type SerInfor struct {
	Ports map[string]SerConn
}

//SerConn record single serial information
type SerConn struct {
	*serial.Config
	*serial.Port
}

// Add implements add a new serial port
func (s *SerInfor) Add(portName string) error {
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
	s.Ports[portName] = conn

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

// Sers record the information of all the client
var Sers SerInfor

func init() {
	Sers.Ports = make(map[string]SerConn)
}

// Operate record the operation to serial
var Operate ExChange

func getDevice() []string {
	cmd := exec.Command("/bin/sh", "-c", "ls /dev/tty* |grep '/dev/ttyACM\\|/dev/ttyUSB'")
	out, _ := cmd.Output()
	outStr := string(out)
	outStr = strings.TrimSpace(outStr)
	devices := strings.Split(outStr, "\n")
	return devices

}

// RunSerial implements all the serial operate
func RunSerial() {
	var exChangeData ExChange
	op := &Operate
FOR:
	for {
		if op.Cmd == WRITEport {
			fmt.Println("i'm writing!")
		}
		switch op.Cmd {
		case QUITprogram:
			for portName := range Sers.Ports {
				Sers.Delete(portName)
			}
			break FOR
		case GETdevice:
			exChangeData.Cmd = GETdevice
			exChangeData.Msg = getDevice()
			Clients.Broadcast(exChangeData)
			op.Cmd = DEFAULT
		case OPENport:
			var msg = OK
			if _, ok := Sers.Ports[op.Target]; !ok {
				err := Sers.Add(op.Target)
				if err != nil {
					msg = NOK
					fmt.Println("Failed to open port ", op.Target)
				}
			}
			exChangeData.Cmd = OPENport
			exChangeData.Msg = msg
			exChangeData.Target = op.Target
			Clients.Broadcast(exChangeData)
			op.Cmd = READport
		case CLOSEport:
			Sers.Delete(op.Target)
			exChangeData.Cmd = CLOSEport
			exChangeData.Target = op.Target
			Clients.Broadcast(exChangeData)
			op.Cmd = READport
		case WRITEport:
			fmt.Println("i'm writing: ", op.Target)
			port := Sers.Ports[op.Target]
			if msg, ok := op.Msg.(string); ok {
				err := port.WriteStr(msg)
				if err != nil {
					fmt.Println("Failed to write to port ", op.Target)
				}
			}
			op.Cmd = READport
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
						Clients.Broadcast(exChangeData)
						fmt.Println(data)
					}
				}
			}
		default:
			// fmt.Println("i'm default")
			// To do
		}
		time.Sleep(time.Millisecond * 20)
	}
}
