package main

import (
	"fmt"

	"github.com/sy/webterm/models"

	"github.com/astaxie/beego"
	_ "github.com/sy/webterm/routers"
)

func main() {
	fmt.Println("@Webterm by sy.Tang 2016-09-08")
	go models.RunSerial()
	beego.Run()
}
