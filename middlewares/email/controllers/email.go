package controllers

import (
	"encoding/json"
	"github.com/wangdianwen/go-core/middlewares"
	"github.com/wangdianwen/go-core/middlewares/email/models"
	"github.com/wangdianwen/go-core/utils"
)

type Email struct {
	middlewares.Services
}

func (c *Email) EmailVerification() {
	dat := make(map[string]interface{})
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &dat)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: "parse json error: " + err.Error()}
		c.ServeJSON()
		return
	}
	username := ""
	tmp, ok := dat["username"]
	if ok {
		username = tmp.(string)
	}
	if len(username) == 0 || !ok {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid username"}
		c.ServeJSON()
		return
	}
	email := ""
	tmp, ok = dat["email"]
	if ok {
		email = tmp.(string)
	}
	if len(email) == 0 || !ok {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid email"}
		c.ServeJSON()
		return
	}
	app := ""
	tmp, ok = dat["app"]
	if ok {
		app = tmp.(string)
	}
	if len(app) == 0 || !ok {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid app name"}
		c.ServeJSON()
		return
	}

	verifyUrl := ""
	tmp, ok = dat["verifyUrl"]
	if ok {
		verifyUrl = tmp.(string)
	}
	if len(verifyUrl) == 0 || !ok {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid verify url"}
		c.ServeJSON()
		return
	}

	res, berr := models.ActiveAccount(verifyUrl, email, username, app)
	if berr != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "create sending task error: " + err.Error()}
		c.ServeJSON()
		return
	}
	if res != true {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "error occurred!"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success!"}
	c.ServeJSON()
}

func (c *Email) TestEmail() {
	//list, err := TaskModel.UnDoTasks(TaskModel.TaskTypeEmail, 10)
	//if err != nil {
	//	c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "create sending task error: " + err.Error()}
	//	c.ServeJSON()
	//	return
	//}
	//host := beego.AppConfig.String("SMTPHost")
	//port, _ := beego.AppConfig.Int("SMTPPort")
	//sender := beego.AppConfig.String("SenderName")
	//pass := beego.AppConfig.String("SenderPass")
	//
	//d := gomail.NewDialer(host, port, sender, pass)
	//s, err := d.Dial()
	//if err != nil {
	//	c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "dial server error: " + err.Error()}
	//	c.ServeJSON()
	//	return
	//}
	//m := gomail.NewMessage()
	//if err != nil {
	//	c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "new message error: " + err.Error()}
	//	c.ServeJSON()
	//	return
	//}
	//for _, r := range list {
	//	var dat map[string]interface{}
	//	err = json.Unmarshal([]byte(r.Data), &dat)
	//	if err != nil {
	//		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "parse json error: " + err.Error()}
	//		c.ServeJSON()
	//		return
	//	}
	//	m.SetHeader("From", sender)
	//	m.SetAddressHeader("To", dat["toAddress"].(string), dat["toName"].(string))
	//
	//	m.SetHeader("Subject", dat["subject"].(string))
	//	// logo pics
	//	m.Embed("assert/images/bizex_logo.png")
	//	// m.Embed("assert/images/brunton_footer.jpg")
	//	m.SetBody("text/html", dat["content"].(string))
	//
	//	if err := gomail.Send(s, m); err != nil {
	//		m.Reset()
	//		errStr := fmt.Sprintf("Could not send email to %s: %s", dat["toAddress"].(string), err.Error())
	//		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: errStr}
	//		c.ServeJSON()
	//		return
	//	}
	//	m.Reset()
	//}
	//c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success!"}
	//c.ServeJSON()

	bErr := models.SysAlert("POS", "this is a test <br />")
	if bErr != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "error occurred!"}
		c.ServeJSON()
		return
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "Success!"}
	c.ServeJSON()
}
