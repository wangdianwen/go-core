package controllers

import (
	"github.com/astaxie/beego"
	"github.com/skip2/go-qrcode"
)

// UserController operations for User
type VersionController struct {
	beego.Controller
}

func (c *VersionController) AndroidAPK() {
	var png []byte
	png, _ = qrcode.Encode("https://storage.googleapis.com/brunton-demo-244223.appspot.com/app-debug.apk", qrcode.Medium, 256)
	c.Ctx.Output.ContentType("png")
	_ = c.Ctx.Output.Body(png)
}
