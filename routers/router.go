package routers

import (
	"github.com/astaxie/beego"
	"github.com/sy/webterm/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/ws", &controllers.WsController{})
}
