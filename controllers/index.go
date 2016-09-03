package controllers

import (
	"github.com/astaxie/beego"
)

// MainController implements the main controller
type MainController struct {
	beego.Controller
}

// Get implements method get
func (c *MainController) Get() {
	c.TplName = "index.tpl"
}
