package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	// "github.com/tarm/serial"
)

// WsController implements the main controller
type WsController struct {
	beego.Controller
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Get implements method Get
func (c *WsController) Get() {
	ports := getSerialPort()
	portsJSON, err := json.Marshal(ports)
	if err != nil {
		log.Fatal(err)
	}
	c.EnableRender = false //it's a must to set Controller.EnableRender to false when there's no TplName
	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(portsJSON)
	conn.WriteJSON(string(portsJSON))
	// conn.WriteJSON(portsJSON)
	fmt.Println(string(portsJSON))
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			return
		}
		if err = conn.WriteMessage(messageType, message); err != nil {
			return
		}

	}
}

func getSerialPort() []string {
	cmd := exec.Command("/bin/sh", "-c", "ls /dev/ttyU*")
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	outStr := string(out)
	outStr = strings.TrimSpace(outStr)
	return strings.Split(outStr, "\n")
}
