package controllers

import (
	"encoding/json"
	"go-core/configs/constants"
	"go-core/models"
	"go-core/utils"
	"strconv"
	"time"
)

// PrivilegeController operations for Privilege
type PrivilegeController struct {
	Core
}

// admin ClientSearch ...
// @router /search/ [get]
func (c *PrivilegeController) PrivilegeSearch() {
	userDetail := c.CurrentUserDetail
	// requires super admin Privilege
	if userDetail.IsSystemAdmin == false {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "require role: superadmin"}
		c.ServeJSON()
		return
	}
	status, err := c.GetInt8("status")
	if err != nil {
		status = -1
	}
	name := c.GetString("name")
	mid, _ := c.GetInt("mid")
	lists, err := models.SearchPrivileges(name, status, mid)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: lists}
	}
	c.ServeJSON()
}

// admin add new Privilege ...
func (c *PrivilegeController) SavePrivilege() {
	userDetail := c.CurrentUserDetail
	if !userDetail.IsSystemAdmin {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "only super admin can add new privilege"}
		c.ServeJSON()
		return
	}
	var v models.Privilege
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	if len(v.Controller) == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "privilege must has a name"}
		c.ServeJSON()
		return
	}
	if len(v.Action) == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "privilege must has a action"}
		c.ServeJSON()
		return
	}
	if v.Mid <= 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "privilege must belong to a menu"}
		c.ServeJSON()
		return
	}

	v.Ctime = time.Now()
	v.Mtime = time.Now()
	v.Status = 0
	if _, err := models.CreatePrivilege(&v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	}
	c.ServeJSON()
}

func (c *PrivilegeController) EditPrivilege() {
	userDetail := c.CurrentUserDetail
	if !userDetail.IsSystemAdmin {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "only super admin can add new client"}
		c.ServeJSON()
		return
	}
	idStr := c.Ctx.Input.Param(":pid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetPrivilegeById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "find client error!" + err.Error()}
		c.ServeJSON()
		return
	}
	var dat map[string]interface{}
	err = json.Unmarshal([]byte(c.Ctx.Input.RequestBody), &dat)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	name := dat["name"].(string)
	if len(name) > 0 {
		v.Name = name
	}
	controller := dat["controller"].(string)
	if len(controller) > 0 {
		v.Controller = controller
	}
	action := dat["action"].(string)
	v.Action = action
	mid := dat["mid"].(float64)
	v.Mid = int(mid)
	v.Mtime = time.Now()
	// you can not change status here!
	if err = models.UpdatePrivilegeById(v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	}
	c.ServeJSON()
}

// admin BanUser ...
// @router /banClient/:uid [PATCH]
func (c *PrivilegeController) BanPrivilege() {
	userDetail := c.CurrentUserDetail
	if !userDetail.IsSystemAdmin {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "only super admin can add new client"}
		c.ServeJSON()
		return
	}
	idStr := c.Ctx.Input.Param(":pid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetPrivilegeById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "find privilege detail error!" + err.Error()}
		c.ServeJSON()
		return
	}
	// check status
	if v.Status == constants.StatusForbidden {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorLogic, Msg: "this privilege has already been forbidden!"}
		c.ServeJSON()
		return
	}
	v.Status = constants.StatusForbidden // ban user
	v.Mtime = time.Now()
	if err = models.UpdatePrivilegeById(v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	}
	c.ServeJSON()
}

func (c *PrivilegeController) ReleasePrivilege() {
	userDetail := c.CurrentUserDetail
	if !userDetail.IsSystemAdmin {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "only super admin can add new client"}
		c.ServeJSON()
		return
	}
	idStr := c.Ctx.Input.Param(":pid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetPrivilegeById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "find privilege detail error!" + err.Error()}
		c.ServeJSON()
		return
	}
	// check status
	if v.Status == constants.StatusNormal {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorLogic, Msg: "this client is a normal user!"}
		c.ServeJSON()
		return
	}
	v.Status = constants.StatusNormal // release user
	v.Mtime = time.Now()
	if err = models.UpdatePrivilegeById(v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	}
	c.ServeJSON()
}
