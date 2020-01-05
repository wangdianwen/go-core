package controllers

import (
	"encoding/json"
	"github.com/wangdianwen/go-core.git/models"
	"github.com/wangdianwen/go-core.git/utils"
	"strconv"
	"time"
)

// RoleController operations for Role
type RoleController struct {
	Core
}

// get select options for roles
func (c *RoleController) OptionRoles() {
	// get all roles by current user
	cid, _ := c.GetInt(":cid")
	c.RequireLogin()
	userDetail := c.CurrentUserDetail
	// only super admin can see all the roles
	if !userDetail.IsSystemAdmin {
		cid = userDetail.Info.Cid
	}
	ret, err := models.GetAllRoles(cid)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	} else {
		res := make([]map[string]interface{}, 0)
		cIds := make([]int, 0)
		for _, role := range ret {
			cIds = append(cIds, role.Cid)
		}
		clientRelation, err := models.ClientsRelations(cIds)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Data: res}
			c.ServeJSON()
			return
		}
		for _, role := range ret {
			item := make(map[string]interface{})
			item["id"] = role.Id
			if cid == 0 {
				item["name"] = role.Rolename + "@" + clientRelation[role.Cid].Name
			} else {
				item["name"] = role.Rolename
			}
			res = append(res, item)
		}
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Data: res}
	}
	c.ServeJSON()
}

// admin search roles ...
func (c *RoleController) RoleSearch() {
	cid, _ := c.GetInt("cid")
	userDetail := c.CurrentUserDetail
	// only super admin can see all the roles
	if !userDetail.IsSystemAdmin {
		cid = userDetail.Info.Cid
	}
	status, err := c.GetInt("status")
	if err != nil {
		status = -1
	}
	rName := c.GetString("rname")
	lists, err := models.SearchRoles(rName, cid, status)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "find role error: " + err.Error()}
	}
	lists, err = c.richRoleListInfo(lists)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "Get extra information error" + err.Error()}
		c.ServeJSON()
		return
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: lists}
	c.ServeJSON()
}

// admin add new client ...
func (c *RoleController) SaveRole() {
	userDetail := c.CurrentUserDetail
	dat := make(map[string]interface{})
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &dat)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	v := models.Role{}
	roleName := dat["rolename"].(string)
	if len(roleName) == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "role must has a name"}
		c.ServeJSON()
		return
	}
	inputCid := int(dat["cid"].(float64))
	if inputCid == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "role must has a client id"}
		c.ServeJSON()
		return
	}
	if !userDetail.IsSystemAdmin && inputCid != userDetail.CInfo.Id {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "you are not allowed add a role for this client"}
		c.ServeJSON()
		return
	}
	v.Rolename = roleName
	v.Cid = inputCid
	v.Ctime = time.Now()
	v.Mtime = time.Now()
	v.Status = 0
	if _, err := models.CreateRole(&v); err == nil {
		menuNames := dat["menuNames"].([]interface{})
		privilegeNames := dat["privilegeNames"].([]interface{})
		mNames := make([]string, len(menuNames))
		for _, name := range menuNames {
			mNames = append(mNames, name.(string))
		}
		pNames := make([]string, len(privilegeNames))
		for _, name := range privilegeNames {
			pNames = append(pNames, name.(string))
		}
		err = models.ModifyRolesPrivileges(v.Id, mNames, pNames)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "modify menus and privileges error: " + err.Error()}
			c.ServeJSON()
			return
		}
		v, _ = c.richRoleInfo(v)
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	}
	c.ServeJSON()
}

func (c *RoleController) EditRole() {
	userDetail := c.CurrentUserDetail
	idStr := c.Ctx.Input.Param(":rid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetRoleById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "find role error!" + err.Error()}
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
	inputCid := int(dat["cid"].(float64))
	if inputCid == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "role must has a client id"}
		c.ServeJSON()
		return
	}
	if !userDetail.IsSystemAdmin && inputCid != userDetail.CInfo.Id {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "you are not allowed edit roles for this client"}
		c.ServeJSON()
		return
	}
	rolename := dat["rolename"].(string)
	if len(rolename) == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "role must has a name"}
		c.ServeJSON()
		return
	}
	// you can not change status here!
	v.Mtime = time.Now()
	v.Cid = inputCid
	v.Rolename = rolename
	if err = models.UpdateRoleById(&v); err == nil {
		menuNames := dat["menuNames"].([]interface{})
		privilegeNames := dat["privilegeNames"].([]interface{})
		mNames := make([]string, len(menuNames))
		for _, name := range menuNames {
			mNames = append(mNames, name.(string))
		}
		pNames := make([]string, len(privilegeNames))
		for _, name := range privilegeNames {
			pNames = append(pNames, name.(string))
		}
		err = models.ModifyRolesPrivileges(id, mNames, pNames)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "modify menus and privileges error: " + err.Error()}
			c.ServeJSON()
			return
		}
		v, _ = c.richRoleInfo(v)
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "modify role error: " + err.Error()}
	}
	c.ServeJSON()
}

// admin BanUser ...
// @router /banClient/:uid [PATCH]
func (c *RoleController) BanRole() {
	userDetail := c.CurrentUserDetail
	idStr := c.Ctx.Input.Param(":rid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetRoleById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "find role detail error!" + err.Error()}
		c.ServeJSON()
		return
	}
	if !userDetail.IsSystemAdmin && v.Cid != userDetail.CInfo.Id {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "you are not allowed to ban this role"}
		c.ServeJSON()
		return
	}
	// check status
	if v.Status == 1 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorLogic, Msg: "this role has already been forbidden!"}
		c.ServeJSON()
		return
	}
	v.Status = 1 // ban role
	v.Mtime = time.Now()
	if err = models.UpdateRoleById(&v); err == nil {
		v, _ = c.richRoleInfo(v)
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	}
	c.ServeJSON()
}

func (c *RoleController) ReleaseRole() {
	userDetail := c.CurrentUserDetail
	idStr := c.Ctx.Input.Param(":rid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetRoleById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "find role detail error!" + err.Error()}
		c.ServeJSON()
		return
	}
	// check permission
	if !userDetail.IsSystemAdmin && v.Cid != userDetail.CInfo.Id {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "you are not allowed to ban this role"}
		c.ServeJSON()
		return
	}
	// check status
	if v.Status == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorLogic, Msg: "this role is a normal!"}
		c.ServeJSON()
		return
	}
	v.Status = 0 // release user
	v.Mtime = time.Now()
	if err = models.UpdateRoleById(&v); err == nil {
		v, _ = c.richRoleInfo(v)
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	}
	c.ServeJSON()
}

// rich single role detail
func (c *RoleController) richRoleInfo(role models.Role) (models.Role, error) {
	rids := make([]int, 0)
	rids = append(rids, role.Id)
	relations, err := models.GetRolePrivilegeByRids(rids)
	if err != nil {
		return role, err
	}
	pIds := make([]int, 0)
	mIds := make([]int, 0)
	for _, item := range relations {
		if item.Pid > 0 {
			pIds = append(pIds, item.Pid)
		}
		if item.Mid > 0 {
			mIds = append(mIds, item.Mid)
		}
	}
	var menus []*models.Menu
	if len(mIds) > 0 {
		menus, err = models.GetAllMenuByIds(mIds)
	}
	if err != nil {
		return role, err
	}
	role.Menus = menus
	var privileges []*models.Privilege
	if len(pIds) > 0 {
		privileges, err = models.GetPrivilegesByIds(pIds)
	}
	if err != nil {
		return role, err
	}
	role.Privileges = privileges
	// get role number
	countMap, err := models.CountByRids(rids)
	if err != nil {
		return role, err
	}
	role.Members = countMap[role.Id]
	return role, nil
}

// rich a role list, increase performance
func (c *RoleController) richRoleListInfo(roles []models.Role) ([]models.Role, error) {
	rids := make([]int, 0)
	for _, item := range roles {
		rids = append(rids, item.Id)
	}
	relations, err := models.GetRolePrivilegeByRids(rids)
	if err != nil {
		return roles, err
	}
	pIds := make([]int, 0)
	mIds := make([]int, 0)
	for _, item := range relations {
		if item.Pid > 0 {
			pIds = append(pIds, item.Pid)
		}
		if item.Mid > 0 {
			mIds = append(mIds, item.Mid)
		}
	}

	// combine menu data
	menus, err := models.GetAllMenuByIds(mIds)
	if err != nil {
		return roles, err
	}

	menuRelation := make(map[int]*models.Menu)
	for _, menu := range menus {
		menuRelation[menu.Id] = menu
	}

	// combine privilege data
	privileges := make([]*models.Privilege, 0)
	if len(pIds) > 0 {
		privileges, err = models.GetPrivilegesByIds(pIds)
	}
	if err != nil {

		return roles, err
	}

	privilegeRelation := make(map[int]*models.Privilege)
	for _, privilege := range privileges {
		privilegeRelation[privilege.Id] = privilege
	}
	rolesMenuRelation := make(map[int][]*models.Menu)
	rolesPrivilegeRelation := make(map[int][]*models.Privilege)
	for _, relation := range relations {
		rid := relation.Rid
		if _, exists := rolesMenuRelation[rid]; !exists {
			rolesMenuRelation[rid] = make([]*models.Menu, 0)
		}
		if _, exists := rolesPrivilegeRelation[rid]; !exists {
			rolesPrivilegeRelation[rid] = make([]*models.Privilege, 0)
		}
		if tmp, exists := menuRelation[relation.Mid]; exists {
			rolesMenuRelation[rid] = append(rolesMenuRelation[rid], tmp)
		}
		if tmp, exists := privilegeRelation[relation.Pid]; exists {
			rolesPrivilegeRelation[rid] = append(rolesPrivilegeRelation[rid], tmp)
		}
	}
	// combine role number
	countMap, err := models.CountByRids(rids)
	if err != nil {
		return roles, err
	}
	for index, role := range roles {
		role.Members = countMap[role.Id]
		role.Privileges = rolesPrivilegeRelation[role.Id]
		role.Menus = rolesMenuRelation[role.Id]
		roles[index] = role
	}
	return roles, nil
}
