package controllers

import (
	"encoding/json"
	"github.com/wangdianwen/go-core.git/configs/constants"
	"github.com/wangdianwen/go-core.git/models"
	"github.com/wangdianwen/go-core.git/utils"
	"strconv"
	"strings"
	"time"
)

// ClientController operations for Client
type ClientController struct {
	Core
}

// OptionClients ...
// @Title OptionClients
// @Description get all the options of nornal clients
// @Success 201 {int} models.Client
// @Failure 403 body is empty
// @router /optionClients [get]
func (c *ClientController) OptionClients() {
	userDetail := c.CurrentUserDetail
	cid := 0
	if !userDetail.IsSystemAdmin && userDetail.Info != nil {
		cid = userDetail.Info.Cid
	}
	ret, err := models.GetAllClientsById(cid)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
	} else {
		res := make([]map[string]interface{}, 0)
		for _, client := range ret {
			item := make(map[string]interface{})
			item["id"] = strconv.Itoa(client.Id)
			item["name"] = client.Name
			res = append(res, item)
		}
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: res}
	}
	c.ServeJSON()
}

// ClientSuggestion ...
// @router /client_suggestion/:clientName/:limit [get]
func (c *ClientController) ClientSuggestion() {
	userDetail := c.CurrentUserDetail
	// requires super admin Privilege
	if userDetail.IsSystemAdmin == false {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "require role: superadmin"}
		c.ServeJSON()
		return
	}
	clientName := c.GetString(":clientName")
	limit, err := c.GetInt(":limit")
	if err != nil {
		limit = 10
	}
	if clientName == "" {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "empty clinet name"}
		c.ServeJSON()
		return
	}
	clients, err := models.ClientSuggestions(clientName, limit)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	list := make([]map[string]interface{}, 0)
	for _, client := range clients {
		item := make(map[string]interface{})
		item["id"] = strconv.Itoa(client.Id)
		item["name"] = client.Name
		list = append(list, item)
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: list}
	c.ServeJSON()
}

// admin ClientSearch ...
// @router /search/ [get]
func (c *ClientController) ClientSearch() {
	userDetail := c.CurrentUserDetail
	// requires super admin Privilege
	if userDetail.IsSystemAdmin == false {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "require role: superadmin"}
		c.ServeJSON()
		return
	}
	status, err := c.GetInt("status")
	if err != nil {
		status = -1
	}
	cnameStr := c.GetString("cnames")
	cnameStr = strings.TrimSpace(cnameStr)
	cnames := make([]string, 0)
	if len(cnameStr) > 0 {
		cnames = strings.Split(cnameStr, ",")
	}
	lists, err := models.SearchClients(cnames, status)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
	} else {
		cids := make([]int, 0)
		for _, item := range lists {
			cids = append(cids, item.Id)
		}
		countMap, err := models.CountByClient(cids)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "Get current user number error" + err.Error()}
			c.ServeJSON()
			return
		}
		for _, item := range lists {
			if count, exist := countMap[item.Id]; exist {
				item.CurrentUserNumber = count
			}
		}
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: lists}
	}
	c.ServeJSON()
}

// admin add new client ...
// @router /saveClient/ [post]
func (c *ClientController) SaveClient() {
	userDetail := c.CurrentUserDetail
	if !userDetail.IsSystemAdmin {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "only super admin can add new client"}
		c.ServeJSON()
		return
	}
	var v models.Client
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	if len(v.Name) == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "client must has a name"}
		c.ServeJSON()
		return
	}
	if len(v.Description) == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "client must has a description"}
		c.ServeJSON()
		return
	}
	if len(v.Phone) == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "user must has a initial password"}
		c.ServeJSON()
		return
	}

	if v.AllowUserlog != 0 && v.AllowUserlog != 1 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid value for allow user log"}
		c.ServeJSON()
		return
	}
	if v.MaxUsers == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "at least allow one user"}
		c.ServeJSON()
		return
	}
	if v.Status != 0 && v.Status != 1 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid value for status"}
		c.ServeJSON()
		return
	}
	v.Ctime = time.Now()
	v.Mtime = time.Now()
	if _, err := models.CreateClient(&v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: err.Error()}
	}
	c.ServeJSON()
}

func (c *ClientController) EditClient() {
	userDetail := c.CurrentUserDetail
	if !userDetail.IsSystemAdmin {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "only super admin can add new client"}
		c.ServeJSON()
		return
	}
	idStr := c.Ctx.Input.Param(":cid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetClientById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "find client error!" + err.Error()}
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
	if connection, ok := dat["connection"].(string); ok {
		if len(connection) > 0 {
			v.Connection = strings.TrimRight(connection, "/") + "/"
		}
	}
	backEndType := dat["back_end_type"].(float64)
	v.BackEndType = int8(backEndType)
	if v.BackEndType != constants.ClientBackEndTypeMySql && v.BackEndType != constants.ClientBackEndTypeAdvanced && v.BackEndType != constants.ClientBackEndTypeAccredo {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid value for back end type"}
		c.ServeJSON()
		return
	}
	phone := dat["phone"].(string)
	if len(phone) > 0 {
		v.Phone = phone
	}
	email := dat["email"].(string)
	if len(email) > 0 {
		v.Email = email
	}
	expireTimeStr := dat["expire_time"].(string)
	expireTime, err := time.Parse("2006-01-02", expireTimeStr)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid expire time format"}
	}
	v.ExpireTime = expireTime

	if len(phone) > 0 {
		v.Phone = phone
	}
	allowUserlog := dat["allow_userlog"].(float64)
	v.AllowUserlog = int8(allowUserlog)
	if v.AllowUserlog != 0 && v.AllowUserlog != 1 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid value for allow user log"}
		c.ServeJSON()
		return
	}
	if description, exist := dat["description"]; exist {
		v.Description = description.(string)
	}
	if name, exist := dat["name"]; exist {
		v.Name = name.(string)
	}
	if maxUsers, exist := dat["max_users"]; exist {
		v.MaxUsers = int(maxUsers.(float64))
	}
	// you can not change status here!
	v.Mtime = time.Now()
	if err = models.UpdateClientById(v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
	}
	c.ServeJSON()
}

// admin BanUser ...
// @router /banClient/:uid [PATCH]
func (c *ClientController) BanClient() {
	userDetail := c.CurrentUserDetail
	if !userDetail.IsSystemAdmin {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "only super admin can add new client"}
		c.ServeJSON()
		return
	}
	idStr := c.Ctx.Input.Param(":cid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetClientById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "find edit client error!" + err.Error()}
		c.ServeJSON()
		return
	}
	// check status
	if v.Status == constants.StatusForbidden {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "this client has already been forbidden!"}
		c.ServeJSON()
		return
	}
	v.Status = constants.StatusForbidden // ban user
	v.Mtime = time.Now()
	if err = models.UpdateClientById(v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
	}
	c.ServeJSON()
}

func (c *ClientController) ReleaseClient() {
	userDetail := c.CurrentUserDetail
	if !userDetail.IsSystemAdmin {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "only super admin can add new client"}
		c.ServeJSON()
		return
	}
	idStr := c.Ctx.Input.Param(":cid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetClientById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "find edit client error!" + err.Error()}
		c.ServeJSON()
		return
	}
	// check status
	if v.Status == constants.StatusNormal {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorLogic, Msg: "this client is a normal user!"}
		c.ServeJSON()
		return
	}
	v.Status = constants.StatusForbidden // release users
	v.Mtime = time.Now()
	if err = models.UpdateClientById(v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	}
	c.ServeJSON()
}
