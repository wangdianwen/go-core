package controllers

import (
	"go-core/middlewares"
	"go-core/utils"
)

type Captcha struct {
	middlewares.Services
}

func (c *Captcha) Test() {
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "this is a test"}
	c.ServeJSON()
}
