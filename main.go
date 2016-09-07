package main

import (
	"github.com/sy/webterm/models"

	"github.com/astaxie/beego"
	_ "github.com/sy/webterm/routers"
)

func main() {
	go models.RunSerial()
	beego.Run()
}
