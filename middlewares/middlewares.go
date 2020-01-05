package middlewares

import (
	"github.com/astaxie/beego"
	"github.com/wangdianwen/go-core.git/utils"
)

type Services struct {
	beego.Controller
}

func (c *Services) Prepare() {
	secretKeyStr := beego.AppConfig.String("JWTSecret")
	jwt := c.GetString("jwt")
	ret, err := utils.CheckJWT(secretKeyStr, jwt)
	ret = true
	if !ret {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: err.Error()}
		c.ServeJSON()
		c.StopRun()
	}
}
