package controllers

import (
	"encoding/json"
	"go-core/configs/constants"
	"go-core/models"
	"go-core/utils"
	"strconv"
	"time"
)

// ClientController operations for Client
type AppController struct {
	Core
}

func (c *AppController) OptionApps() {
	// userDetail := c.CurrentUserDetail
	ret, err := models.AppsList("", -1)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
	} else {
		res := make([]map[string]interface{}, 0)
		for _, app := range ret {
			item := make(map[string]interface{})
			item["id"] = app.Id
			item["name"] = app.Name
			res = append(res, item)
		}
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Data: res}
	}
	c.ServeJSON()
}

func (c *AppController) AppsList() {

	appName := c.GetString("name")
	limit, err := c.GetInt(":limit")
	if err != nil {
		limit = 10
	}
	ret, err := models.AppsList(appName, limit)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Data: ret}
	}
	c.ServeJSON()
}

func (c *AppController) SaveApp() {
	userDetail := c.CurrentUserDetail
	if !userDetail.IsSystemAdmin {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "you don't have right to create new application"}
		c.ServeJSON()
		return
	}
	var v models.CoreApps
	var dat map[string]interface{}
	err := json.Unmarshal([]byte(c.Ctx.Input.RequestBody), &dat)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	v.Ctime = time.Now()
	v.Mtime = time.Now()
	name := ""
	_, ok := dat["name"]
	if ok {
		name = dat["name"].(string)
	}
	if len(name) == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid application name"}
		c.ServeJSON()
		return
	}

	appSecret := ""
	_, ok = dat["app_secret"]
	if ok {
		appSecret = dat["app_secret"].(string)
	}
	if len(appSecret) <= 8 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid application secret"}
		c.ServeJSON()
		return
	}

	appKey := ""
	_, ok = dat["app_key"]
	if ok {
		appKey = dat["app_key"].(string)
	}
	if len(appKey) <= 8 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid application key"}
		c.ServeJSON()
		return
	}

	v.Name = name
	v.Status = constants.StatusNormal
	v.AppSecret = appSecret
	v.AppKey = appKey
	if _, err := models.CreateApp(&v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	}
	c.ServeJSON()
}

// admin EditUser ...
func (c *AppController) EditApp() {
	userDetail := c.CurrentUserDetail
	if !userDetail.IsSystemAdmin {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "you don't have right to edit application"}
		c.ServeJSON()
		return
	}
	idStr := c.Ctx.Input.Param(":appid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetAppById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "find application error!" + err.Error()}
		c.ServeJSON()
		return
	}
	var dat map[string]interface{}
	err = json.Unmarshal([]byte(c.Ctx.Input.RequestBody), &dat)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	v.Mtime = time.Now()
	name := ""
	_, ok := dat["name"]
	if ok {
		name = dat["name"].(string)
	}
	if len(name) == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid application name"}
		c.ServeJSON()
		return
	}

	appSecret := ""
	_, ok = dat["app_secret"]
	if ok {
		appSecret = dat["app_secret"].(string)
	}
	if len(appSecret) <= 8 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid application secret"}
		c.ServeJSON()
		return
	}

	appKey := ""
	_, ok = dat["app_key"]
	if ok {
		appKey = dat["app_key"].(string)
	}
	if len(appKey) <= 8 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid application key"}
		c.ServeJSON()
		return
	}

	allowFreeTrial := int8(0)
	_, ok = dat["allow_free_trial"]
	if ok {
		allowFreeTrial = int8(dat["allow_free_trial"].(float64))
	}
	if allowFreeTrial != 0 && allowFreeTrial != 1 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid allow free trial value"}
		c.ServeJSON()
		return
	}

	v.Name = name
	v.AppSecret = appSecret
	v.AppKey = appKey
	v.AllowFreeTrial = allowFreeTrial
	if err := models.UpdateAppById(v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	}
	c.ServeJSON()
}

func (c *AppController) BanApp() {
	userDetail := c.CurrentUserDetail
	if !userDetail.IsSystemAdmin {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "you don't have right to edit application"}
		c.ServeJSON()
		return
	}
	idStr := c.Ctx.Input.Param(":appid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetAppById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "find application error!" + err.Error()}
		c.ServeJSON()
		return
	}
	if v.Name == "Core" {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "you can't forbidden framework core!"}
		c.ServeJSON()
		return
	}

	// check status
	if v.Status == constants.StatusForbidden {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "this application has already been forbidden!"}
		c.ServeJSON()
		return
	}
	v.Status = 1 // ban user
	v.Mtime = time.Now()
	if err = models.UpdateAppById(v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
	}
	c.ServeJSON()
}

func (c *AppController) ReleaseApp() {
	userDetail := c.CurrentUserDetail
	if !userDetail.IsSystemAdmin {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "you don't have right to edit application"}
		c.ServeJSON()
		return
	}
	idStr := c.Ctx.Input.Param(":appid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetAppById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "find application error!" + err.Error()}
		c.ServeJSON()
		return
	}
	// check status
	if v.Status == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "this application has already been released!"}
		c.ServeJSON()
		return
	}
	v.Status = 0 // release user
	v.Mtime = time.Now()
	if err = models.UpdateAppById(v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
	}
	c.ServeJSON()
}

// get all normal status products
func (c *AppController) ValidateProducts() {
	userDetail := c.CurrentUserDetail
	// does not get the user info, clear session of current user
	if userDetail.Info == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorUnLogin, Msg: "require login"}
		c.ServeJSON()
		return
	}
	products, err := models.GetAllValidateAppsList()
	if err != nil && err.Error() == utils.BeegoNoData {
		products = make([]*models.CoreApps, 0)
	} else if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "get validate products failed:" + err.Error()}
		c.ServeJSON()
		return
	}
	// hide Secret key
	for index, item := range products {
		item.AppKey = "***"
		item.AppSecret = "***"
		products[index] = item
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Data: products}
	c.ServeJSON()
}

func (c *AppController) FreeTrialProducts() {
	userDetail := c.CurrentUserDetail
	// get user conditions
	userInfo := userDetail.Info
	if userInfo == nil || userInfo.Id == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorUnLogin, Msg: "require login"}
		c.ServeJSON()
		return
	}
	uid := userInfo.Id
	var dat map[string]interface{}
	err := json.Unmarshal([]byte(c.Ctx.Input.RequestBody), &dat)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	appIds := make([]int, 0)
	_, ok := dat["appids"]
	if ok {
		tmp := dat["appids"].([]interface{})
		for _, item := range tmp {
			appIds = append(appIds, int(item.(float64)))
		}
	}
	bErr, products := models.StartFreeTrial(uid, appIds)
	if bErr != nil {
		c.Data["json"] = &utils.JSONStruct{Code: bErr.Code, Msg: bErr.Message}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Data: products}
	}
	c.ServeJSON()
}
