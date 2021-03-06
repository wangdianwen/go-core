package controllers

import (
	"encoding/json"
	"github.com/wangdianwen/go-core/models"
	"github.com/wangdianwen/go-core/utils"
	"strconv"
	"time"
)

// MenuController operations for Menu
type MenuController struct {
	Core
}

// get select options for menus
func (c *MenuController) OptionMenus() {
	userDetail := c.CurrentUserDetail
	// requires super admin Privilege
	if userDetail.IsSystemAdmin == false {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "require role: superadmin"}
		c.ServeJSON()
		return
	}
	ret, err := models.SearchMenu(-1, "")
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	} else {
		res := make([]map[string]interface{}, 0)
		for _, menu := range ret {
			item := make(map[string]interface{})
			item["id"] = menu.Id
			item["name"] = menu.Name
			res = append(res, item)
		}
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: res}
	}
	c.ServeJSON()
}

// admin ClientSearch ...
func (c *MenuController) MenuSearch() {
	userDetail := c.CurrentUserDetail
	// requires super admin Privilege
	if userDetail.IsSystemAdmin == false {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "require role: superadmin"}
		c.ServeJSON()
		return
	}
	name := c.GetString("name")
	paid, err := c.GetInt("paid")
	if err != nil {
		paid = -1
	}
	lists, err := models.SearchMenu(paid, name)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	} else {
		paids := make([]int, 0)
		mids := make([]int, 0)
		for _, item := range lists {
			paids = append(paids, item.Id)
			mids = append(mids, item.Id)
		}
		if len(paids) > 0 {
			relation, err := models.GetAllMenusByPaids(paids)
			if err != nil {
				c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "find children error " + err.Error()}
				c.ServeJSON()
				return
			}
			for _, item := range lists {
				if children, exists := relation[item.Id]; exists {
					item.Children = children
				}
			}
		}
		if len(mids) > 0 {
			relation, err := models.GetAllPrivilegeByMidsGroupByMids(mids)
			if err != nil {
				c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "find sub privileges error " + err.Error()}
				c.ServeJSON()
				return
			}
			for _, item := range lists {
				if privileges, exists := relation[item.Id]; exists {
					item.Privileges = privileges
				}
			}
		}
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: lists}
	}
	c.ServeJSON()
}

// admin add new Privilege ...
func (c *MenuController) SaveMenu() {
	userDetail := c.CurrentUserDetail
	if !userDetail.IsSystemAdmin {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "only super admin has menu operations"}
		c.ServeJSON()
		return
	}
	var v models.Menu
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParseJson, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	if len(v.Name) == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "menu must has a name"}
		c.ServeJSON()
		return
	}
	if len(v.Url) == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "menu must has a path"}
		c.ServeJSON()
		return
	}
	if v.Appid <= 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "menu must belong to an application"}
		c.ServeJSON()
		return
	}
	v.Ctime = time.Now()
	v.Mtime = time.Now()
	if _, err := models.CreateMenu(&v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	}
	c.ServeJSON()
}

func (c *MenuController) EditMenu() {
	userDetail := c.CurrentUserDetail
	if !userDetail.IsSystemAdmin {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "only super admin has menu operations"}
		c.ServeJSON()
		return
	}
	idStr := c.Ctx.Input.Param(":mid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetMenuById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "find client error!" + err.Error()}
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
	name := dat["name"].(string)
	if len(name) > 0 {
		v.Name = name
	}
	url := dat["url"].(string)
	if len(url) > 0 {
		v.Url = url
	}
	visible := dat["visible"].(float64)
	if visible == 0 || visible == 1 {
		v.Visible = int8(visible)
	}
	appid := dat["appid"].(float64)
	if appid > 0 {
		v.Appid = int(appid)
	}

	v.Mtime = time.Now()
	// you can not change status here!
	if err = models.UpdateMenuById(v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	}
	c.ServeJSON()
}

// admin BanUser ...
func (c *MenuController) DelMenu() {
	userDetail := c.CurrentUserDetail
	if !userDetail.IsSystemAdmin {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "only super admin has menu operations"}
		c.ServeJSON()
		return
	}
	idStr := c.Ctx.Input.Param(":mid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetMenuById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "find menu detail error!" + err.Error()}
		c.ServeJSON()
		return
	}
	// check if there is child menu
	total, err := models.CountMenu(v.Id, "")
	if total > 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorLogic, Msg: "this menu contains submenus, please delete submenus first"}
		c.ServeJSON()
		return
	}
	// check if there are privileges bind to this menu
	total, err = models.CountPrivilegesByMid(v.Id)
	if total > 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorLogic, Msg: "this menu contains privileges, please delete submenus first"}
		c.ServeJSON()
		return
	}
	v.Mtime = time.Now()
	if err = models.DeleteMenu(v.Id); err == nil {
		// delete all menu Privilege for users
		if _, err := models.DelRolePrivilegeByMid(v.Id); err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "delete privilege error:" + err.Error()}
			c.ServeJSON()
			return
		}
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: err.Error()}
	}
	c.ServeJSON()
}

// get current menu
func (c *MenuController) CurrentMenu() {
	userDetail := c.CurrentUserDetail
	if userDetail.Info == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorUnLogin, Msg: "require login", Data: nil}
		c.ServeJSON()
		return
	}
	menus := userDetail.Menus
	menus = generateMenuTree(menus, 0, 0)
	mList := formatMenuData(menus)
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: mList}
	c.ServeJSON()
}


// get current menu
func (c *MenuController) CurrentMenuNames() {
	userDetail := c.CurrentUserDetail
	if userDetail.Info == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorUnLogin, Msg: "require login", Data: nil}
		c.ServeJSON()
		return
	}
	menus := userDetail.Menus
	menus = generateMenuTree(menus, 0, 0)
	var mNames []string
	formatMenuNames(menus, "", &mNames)
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: mNames}
	c.ServeJSON()
}

func generateMenuTree(menus []*models.Menu, pid int, level int) []*models.Menu {
	ret := make([]*models.Menu, 0)
	level += 1
	for _, menu := range menus {
		if menu.Paid == pid {
			menu.Level = level
			menu.Children = generateMenuTree(menus, menu.Id, level)
			ret = append(ret, menu)
		}
	}
	return ret
}

func formatMenuData(menus []*models.Menu) []map[string]interface{} {
	var mList []map[string]interface{}
	for _, menu := range menus {
		tmp := make(map[string]interface{})
		tmp["name"] = menu.Name
		tmp["path"] = menu.Url
		tmp["authority"] = menu.Name
		if menu.Visible == 1 {
			tmp["hideInMenu"] = true
		}
		if len(menu.Children) != 0 {
			tmp["children"] = formatMenuData(menu.Children)
		}
		mList = append(mList, tmp)
	}
	return mList
}

func formatMenuNames(menus []*models.Menu, preName string, ret *[]string) {
	for _, menu := range menus {
		if menu.Paid == 0 {
			preName = menu.Name
			*ret = append(*ret, preName)
		}
		if len(menu.Children) > 0 {
			if menu.Paid > 0 {
				preName = preName + "." + menu.Name
				*ret = append(*ret, preName)
			}
			formatMenuNames(menu.Children, preName, ret)
		} else {
			menuName := preName + "." + menu.Name
			*ret = append(*ret, menuName)
		}
	}
}





