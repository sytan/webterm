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

type exChangeFormat struct {
	Cmd string
	Msg interface{}
}

var exChangeMsg exChangeFormat
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Get implements method Get
func (c *WsController) Get() {
	ports := getSerialPort()
	exChangeMsg.Cmd = "select"
	exChangeMsg.Msg = ports

	exChangeJSON, err := json.Marshal(exChangeMsg)
	if err != nil {
		log.Fatal(err)
	}
	c.EnableRender = false //it's a must to set Controller.EnableRender to false when there's no TplName
	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.WriteJSON(string(exChangeJSON))
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}
		var msgUnmarshal exChangeFormat
		err = json.Unmarshal(message, &msgUnmarshal)
		if err != nil {
			fmt.Print("i'm failed to unmarshal")
		}
		exChangeMsg.Cmd = msgUnmarshal.Cmd
		exChangeMsg.Msg = msgUnmarshal.Msg
		exChangeJSON, err := json.Marshal(exChangeMsg)
		if err != nil {
			fmt.Println("woring format of json")
		}
		if err = conn.WriteJSON(string(exChangeJSON)); err != nil {
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
