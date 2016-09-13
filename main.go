package main

import (
	"fmt"
	"os"
	"time"

	"os/signal"

	"github.com/astaxie/beego"
	"github.com/sy/webterm/models"
	_ "github.com/sy/webterm/routers"
)

func main() {
	fmt.Println("@Webterm by sy.Tang 2016-09-13")
	go quit()
	go models.RunSerial()
	beego.Run()
}
func quit() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)
	<-sigs
	for _, user := range models.Clients.Users {
		user.Conn.Close()
	}
	models.CloseSerial() //close serial after close websocket
	time.Sleep(time.Millisecond * 200)
	os.Exit(1)
}
