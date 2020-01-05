package controllers

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/wangdianwen/go-core/configs/constants"
	"github.com/wangdianwen/go-core/configs/structures"
	"github.com/wangdianwen/go-core/models"
	"github.com/wangdianwen/go-core/utils"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego"
)

type Core struct {
	CurrentUserDetail structures.UserDetail
	Cid               int
	beego.Controller
}

func (c *Core) Prepare() {
	// if login get user detail
	if c.GetSession("uname") != nil && c.GetSession("passwd") != nil {
		uname := c.GetSession("uname").(string)
		passwd := c.GetSession("passwd").(string)
		CurrentUser, err := CurrentUserDetails(uname, passwd, "")
		if err == nil {
			c.CurrentUserDetail = CurrentUser
			c.Cid = CurrentUser.Info.Cid
		} else {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "fail to get login user detail, error:" + err.Error()}
			c.DelSession("uname")
			c.DelSession("passwd")
			c.DelSession("uid")
			c.ServeJSON()
			c.StopRun()
		}
	}
	// get current Privilege info
	bErr := c.CheckPrivileges()
	if bErr.Error != nil {
		c.Data["json"] = &utils.JSONStruct{Code: bErr.Code, Msg: bErr.Message}
		if bErr.Code == utils.ErrorUnLogin {
			c.Ctx.Output.SetStatus(utils.ErrorUnLogin)
		}
		c.ServeJSON()
		c.StopRun()
	}
}

func (c *Core) RequireAdmin() {
	c.RequireLogin()
	if !c.CurrentUserDetail.IsSystemAdmin {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "require system admin", Data: nil}
		c.ServeJSON()
		c.StopRun()
		return
	}
}

func (c *Core) RequireLogin() {
	userInfo := c.CurrentUserDetail.Info
	if userInfo == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorUnLogin, Msg: "require login", Data: nil}
		c.ServeJSON()
		c.StopRun()
		return
	}
	loginToken := ""
	if c.GetSession("loginToken") != nil {
		loginToken = c.GetSession("loginToken").(string)
	}
	if loginToken != userInfo.LastLoginToken {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorUnLogin, Msg: "login somewhere else, please login again"}
		c.ServeJSON()
		c.StopRun()
	}
}

// check Privileges
func (c *Core) CheckPrivileges() utils.BError {
	userDetail := c.CurrentUserDetail
	currentController, currentUserAction := c.GetControllerAndAction()
	// println(currentController, currentUserAction)
	privilege, err := models.GetPrivilegeByCA(currentController, currentUserAction)
	if err != nil {
		return utils.BError{Code: utils.ErrorDB, Message: "fail to get privilege due to db error: " + err.Error(), Error: err}
	}
	checked := false
	if privilege != nil {
		c.RequireLogin()
		log := &models.UserLog{
			Ctime: time.Now(),
			Pid:   privilege.Id,
			Uid:   userDetail.Info.Id,
			Data:  FormatLogData(c.Ctx.Request),
		}
		_, _ = models.CreateUserLog(log)
	}
	// super admin
	if userDetail.IsSystemAdmin == true {
		checked = true
		// normal user
	} else {
		// privilege is not null, check user login and privileges
		if privilege != nil {
			pid := privilege.Id
			for _, userPrivilege := range userDetail.Privileges {
				if pid == userPrivilege.Id {
					checked = true
					break
				}
			}
			// privilege is nil, don't need to check anything, for example login page
		} else {
			checked = true
		}
	}
	if !checked {
		return utils.BError{Code: utils.ErrorForbidden, Message: "You don't have the permission to access this api", Error: errors.New("You don't have the permission to access this api")}
	}
	return utils.BError{Code: utils.Success}
}

func (c *Core) Md5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

func ParseEmail(email string) (err error, domain string, username string) {
	components := strings.Split(email, "@")
	if len(components) != 2 {
		err = errors.New("invalid email address")
		return
	}
	username, domain = strings.TrimSpace(components[0]), strings.TrimSpace(components[1])
	return
}

// password must be md5
func CurrentUserDetails(username string, password string, barcode string) (structures.UserDetail, error) {
	userDetail := structures.UserDetail{}
	userDetail.IsSystemAdmin = false
	userDetail.Info = &models.User{}
	userDetail.CInfo = &models.Client{}
	userDetail.Roles = make([]*models.Role, 0)
	userDetail.Privileges = make([]*models.Privilege, 0)
	userDetail.Menus = make([]*models.Menu, 0)
	userDetail.Products = make([]*models.CoreApps, 0)

	var v *models.User
	var err error
	if username != "" && password != "" {
		v, err = models.GetUserLogin(username, password)
	} else if barcode != "" {
		v, err = models.GetBarcodeLogin(&barcode)
	} else {
		return userDetail, errors.New("empty username or password")
	}
	if err != nil {
		return userDetail, errors.New("user is not exist or incorrect password")
	}
	// check if is system admin
	userDetail.Info = v
	// check if current user is system admin
	email := userDetail.Info.Email
	SystemAdmins := beego.AppConfig.Strings("SystemAdmins")
	for _, adminEmail := range SystemAdmins {
		email = strings.ToLower(email)
		adminEmail = strings.ToLower(adminEmail)
		userNameChecked := false
		domainChecked := false
		err, userDomain, userName := ParseEmail(email)
		if err != nil {
			userDetail.IsSystemAdmin = false
			break
		}
		err, adminDomain, adminName := ParseEmail(adminEmail)
		if err != nil {
			userDetail.IsSystemAdmin = false
			break
		}

		if userName == adminName || adminName == "*" {
			userNameChecked = true
		}
		if userDomain == adminDomain {
			domainChecked = true
		}

		if userNameChecked && domainChecked {
			userDetail.IsSystemAdmin = true
			break
		}
	}
	// not administrator and user forbidden
	if !userDetail.IsSystemAdmin && v.Status == constants.StatusForbidden {
		return userDetail, errors.New("current user has been forbidden ")
	}
	// user type | expired user
	if !userDetail.IsSystemAdmin && v.Type == constants.TypeFreeTrialExpireUser {
		return userDetail, errors.New("current user has expired and not activated ")
	}
	// get user client
	client, err := userClients(userDetail)
	if err != nil {
		return userDetail, err
	}
	userDetail.CInfo = client

	// get user roles/products/menus
	roles, err := userRoles(userDetail)
	if err != nil {
		return userDetail, err
	}
	userDetail.Roles = roles
	// get all front-end menus
	privileges, menus, products, err := UserMenusProductsPrivileges(userDetail, roles)
	if err != nil {
		return userDetail, err
	}
	userDetail.Privileges = privileges
	userDetail.Menus = menus
	userDetail.Products = products
	// get clients infos
	return userDetail, nil
}

func userClients(userDetail structures.UserDetail) (*models.Client, error) {
	// 1. check validate user
	uid := userDetail.Info.Id
	if uid <= 0 {
		return nil, errors.New("invalid user information")
	}
	// 2. check user type, only normal user has a client information
	uType := userDetail.Info.Type
	if uType != constants.TypeNormalUser {
		return nil, nil
	}
	// 3. check system admin

	cid := userDetail.Info.Cid
	client, err := models.GetClientById(cid)
	if err != nil && err.Error() == utils.BeegoNoData {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	// jude if client is forbidden
	if !userDetail.IsSystemAdmin {
		if client.Status == 1 {
			return nil, errors.New("this client has been forbidden")
		}
		if client.ExpireTime.Before(time.Now()) {
			return nil, errors.New("this client has expired at " + client.ExpireTime.Format("2006/03/04 15:16:17"))
		}
	}
	return client, nil
}

func userRoles(userDetail structures.UserDetail) ([]*models.Role, error) {
	// 1. check validate user
	uid := userDetail.Info.Id
	if uid <= 0 {
		return nil, errors.New("invalid user information")
	}
	// 2. check user type, only normal user has roles
	uType := userDetail.Info.Type
	if uType != constants.TypeNormalUser {
		return nil, nil
	}
	roleRelations, err := models.GetValidRoles(uid)
	if err != nil && err.Error() == utils.BeegoNoData {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	roleIds := make([]int, 0)
	for _, role := range roleRelations {
		roleIds = append(roleIds, role.Rid)
	}
	if len(roleIds) == 0 {
		return nil, nil
	}
	roles, err := models.GetRolesByIds(roleIds)
	if err != nil && err.Error() == utils.BeegoNoData {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return roles, nil
}

// get user menus
func UserMenusProductsPrivileges(userDetail structures.UserDetail, roles []*models.Role) ([]*models.Privilege, []*models.Menu, []*models.CoreApps, error) {
	// 1. check validate user
	uid := userDetail.Info.Id
	if uid <= 0 {
		return nil, nil, nil, errors.New("invalid user information")
	}
	// 2. check super admin
	superAdmin := userDetail.IsSystemAdmin
	if superAdmin {
		return superAdminUserMenusProductsPrivileges()
	}
	// 3. check user type
	uType := userDetail.Info.Type
	if uType == constants.TypeNormalUser {
		return normalUserMenusProductsPrivileges(roles)
	} else if uType == constants.TypeFreeTrialUser || uType == constants.TypeEmptyUser {
		return freeTrialUserMenusProductsPrivileges(uid)
	}
	return nil, nil, nil, errors.New("user free trial has expired or unknown user type")
}

// get free trial user menu/products/privileges
// get all apps => get all menus => get all Privileges
func freeTrialUserMenusProductsPrivileges(uid int) ([]*models.Privilege, []*models.Menu, []*models.CoreApps, error) {
	list, err := models.GetUserFreeTrialProductsByUid(uid)
	if err != nil {
		return nil, nil, nil, errors.New("free trial user MPP error:" + err.Error())
	}
	appIds := make([]int, 0)
	for _, item := range list {
		appIds = append(appIds, item.Appid)
	}
	appList, err := models.GetAppsByIds(appIds)
	if err != nil {
		return nil, nil, nil, errors.New("free trial user product list error:" + err.Error())
	}
	// filter forbidden apps
	appIds = make([]int, 0)
	for _, item := range appList {
		if item.Status == constants.StatusForbidden {
			continue
		}
		appIds = append(appIds, item.Id)
	}
	menuList, err := models.GetAllMenuByAppids(appIds)
	if err != nil {
		return nil, nil, nil, errors.New("free trial user menu list error:" + err.Error())
	}
	mids := make([]int, 0)
	for _, item := range menuList {
		mids = append(mids, item.Id)
	}
	// get all Privileges by menuid
	privilegeList, err := models.GetAllPrivilegesByMids(mids)
	if err != nil {
		return nil, nil, nil, errors.New("free trial user privilege list error:" + err.Error())
	}
	return privilegeList, menuList, appList, nil
}

// get super admin menu/products/privileges
// all apps | all menus | all privileges
func superAdminUserMenusProductsPrivileges() ([]*models.Privilege, []*models.Menu, []*models.CoreApps, error) {
	appList, err := models.GetAllAppsList(constants.StatusNormal)
	if err != nil {
		return nil, nil, nil, errors.New("super admin user product list error:" + err.Error())
	}
	menuList, err := models.SearchMenu(-1, "") // search all
	if err != nil {
		return nil, nil, nil, errors.New("super admin user menu list error:" + err.Error())
	}
	privilegeList, err := models.SearchPrivileges("", constants.StatusNormal, -1)
	if err != nil {
		return nil, nil, nil, errors.New("super admin user privilege list error:" + err.Error())
	}
	return privilegeList, menuList, appList, nil
}

// get normal user menu/products/privileges
// roles => role-menu relation => memulist | roles => role-privilege relation => privileges list | menulist => applist
func normalUserMenusProductsPrivileges(roles []*models.Role) ([]*models.Privilege, []*models.Menu, []*models.CoreApps, error) {
	rids := make([]int, 0)
	for _, item := range roles {
		rids = append(rids, item.Id)
	}
	// user do not has nay roles
	if len(rids) == 0 {
		return nil, nil, nil, nil
	}
	// find menu ids by role
	menuIds := make([]int, 0)
	roleMenus, err := models.GetAllMenusByRids(rids)
	if err != nil {
		return nil, nil, nil, errors.New("normal user role, menu relation error:" + err.Error())
	}
	for _, menu := range roleMenus {
		menuIds = append(menuIds, menu.Mid)
	}
	// if normal user doesn't has any menu, he/she doesn't has product and privilege either
	if len(menuIds) == 0 {
		return nil, nil, nil, nil
	}
	// get all validate applications
	apps, err := models.GetAllAppsList(0)
	if err != nil {
		return nil, nil, nil, errors.New("normal user role, get normal applications error:" + err.Error())
	}
	appIds := make([]int, 0)
	appRelations := make(map[int]*models.CoreApps)
	for _, item := range apps {
		appIds = append(appIds, item.Id)
		appRelations[item.Id] = item
	}
	// if Application is forbidden, shouldn't show in menus
	menuList, err := models.GetAllMenuByMidsAndAppids(menuIds, appIds)
	if err != nil {
		return nil, nil, nil, errors.New("normal user role, get menu list error:" + err.Error())
	}
	tmp := make(map[int]*models.CoreApps)
	// check menu, find user owned products
	for _, item := range menuList {
		tmp[item.Appid] = appRelations[item.Appid]
	}
	appList := make([]*models.CoreApps, 0)
	for _, item := range tmp {
		appList = append(appList, item)
	}

	// get Privileges
	pids := make([]int, 0)
	privilegeRoles, err := models.GetRolePrivilegeByRids(rids)
	if err != nil {
		return nil, nil, nil, errors.New("normal user role, get privilege role list error:" + err.Error())
	}
	for _, privilegeRole := range privilegeRoles {
		pids = append(pids, privilegeRole.Pid)
	}
	// no available privileges
	if len(pids) == 0 {
		return nil, nil, nil, nil
	}

	privilegeList, err := models.GetPrivilegesByIds(pids)
	if err != nil {
		return nil, nil, nil, errors.New("normal user role, get privilege list error:" + err.Error())
	}
	return privilegeList, menuList, appList, nil
}

func FormatLogData(request *http.Request) string {
	body, _ := ioutil.ReadAll(request.Body)
	return "url:" + request.URL.String() + ", body: " + string(body) + "\n"
}

//type Handle struct {
//	Host string
//	Port string
//}
//
//func (handle *Handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	remote, err := url.Parse("http://" + handle.Host + ":" + handle.Port)
//	if err != nil {
//		panic(err)
//	}
//	proxy := httputil.NewSingleHostReverseProxy(remote)
//	proxy.ServeHTTP(w, r)
//}
