package controllers

import "github.com/astaxie/beego"

// MainController implements the main controller
type MainController struct {
	beego.Controller
}

// Get implements method get
func (c *MainController) Get() {
	c.TplName = "index.tpl"
}

// // Post implements method post
// func (c *MainController) Post() {
// 	c.EnableRender = false //it's a must to set Controller.EnableRender to false when there's no TplName
// 	fmt.Println("i'm post")
// }
