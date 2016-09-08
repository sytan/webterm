package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/gorilla/websocket"
	"github.com/sy/webterm/models"
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
	c.EnableRender = false //it's a must to set Controller.EnableRender to false when there's no TplName
	conn, err := upgrader.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	remoteAddr := conn.RemoteAddr().String()
	models.Clients.Add(remoteAddr, conn)

	var dataUnmarshal models.ExChange
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return
		}
		json.Unmarshal(message, &dataUnmarshal)
		fmt.Println("There is a message ", dataUnmarshal.Cmd)
		models.Operate = dataUnmarshal
	}
}
