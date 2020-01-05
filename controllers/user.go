package controllers

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"go-core/configs/constants"
	MWCaptcha "go-core/middlewares/captcha/models"
	MWEmail "go-core/middlewares/email/models"
	"go-core/models"
	"go-core/utils"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// UserController operations for User
type UserController struct {
	Core
}

// Login ...
// @Title LoginFF
// @Description User Login
// @Param	uname		path 	string	true		"user name"
// @Param	passwd		path 	string	true		"password"
// @Success utils.Success {string} login success!
// @Failure 403 id is empty
// @router /login/:uname/:passwd [get]
func (c *UserController) Login() {
	uname := c.GetString(":uname")
	passwd := c.GetString(":passwd")
	barcode := c.GetString(":barcode")
	if passwd != "" {
		passwd = c.Md5(passwd)
	}
	// in order to use configs func, set session first
	userDetail, err := CurrentUserDetails(uname, passwd, barcode)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	// set session
	c.SetSession("uname", userDetail.Info.Uname)
	c.SetSession("passwd", userDetail.Info.Passwd)
	// update
	v := userDetail.Info
	v.LastLoginIp = c.Ctx.Request.RemoteAddr
	v.LastLoginTime = time.Now()
	v.LastLoginToken = fmt.Sprintf("%x", md5.Sum([]byte(v.LastLoginIp+strconv.FormatInt(int64(v.LastLoginTime.Nanosecond()), 10))))
	c.SetSession("loginToken", v.LastLoginToken)
	// check user is free trial and expired
	if v.Type == constants.TypeFreeTrialUser && v.Ctime.AddDate(0, 0, constants.FreeTrialMaxDays).Before(time.Now()) {
		v.Type = constants.TypeFreeTrialExpireUser
		v.Status = constants.StatusForbidden
	}
	_ = models.UpdateUserById(v)
	c.CruSession.SessionRelease(c.Ctx.ResponseWriter)
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Data: userDetail}
	c.ServeJSON()
}

// user logout
func (c *UserController) LogOut() {
	c.DelSession("uname")
	c.DelSession("passwd")
	c.DelSession("uid")
	c.CruSession.SessionRelease(c.Ctx.ResponseWriter)
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Data: nil}
	c.ServeJSON()
}

// CurrentUser ...
// @Title CurrentUser
// @Description check login user info
// @Success utils.Success {string} login success!
// @Failure 403 user is not login
// @router /current_user/ [get]
func (c *UserController) CurrentUser() {
	userDetail := c.CurrentUserDetail
	// does not get the user info, clear session of current user
	if userDetail.Info == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorUnLogin, Msg: "require login"}
		c.DelSession("uname")
		c.DelSession("passwd")
		c.DelSession("uid")
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: userDetail}
		c.SetSession("uid", userDetail.Info.Id)
	}
	c.ServeJSON()
}

// UserSuggestion ...
// @Title UserSuggestion
// @Description get user suggestions by name
// @Success utils.Success {string} return user name and id lists success!
// @Failure 403 user is not login
// @router /user_suggestion/:userName/:limit [get]
func (c *UserController) UserSuggestion() {
	userDetail := c.CurrentUserDetail
	cid := userDetail.Info.Cid
	userName := c.GetString(":userName")
	limit, err := c.GetInt(":limit")
	if err != nil {
		limit = 10
	}
	if userName == "" {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "invalid user name"}
		c.ServeJSON()
		return
	}
	users, err := models.UserSuggestions(userName, cid, limit)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	list := make([]map[string]interface{}, 0)
	for _, user := range users {
		item := make(map[string]interface{})
		item["id"] = strconv.Itoa(user.Id)
		item["name"] = user.Uname
		list = append(list, item)
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: list}
	c.ServeJSON()
}

// admin SearchUsers ...
func (c *UserController) SearchUsers() {
	userDetail := c.CurrentUserDetail

	utype, err := c.GetInt("type")
	if err != nil {
		utype = -1
	}

	emailVerifyStatus, err := c.GetInt("email_verify_status")
	if err != nil {
		emailVerifyStatus = -1
	}

	cid, err := c.GetInt("cid")
	if err != nil {
		cid = -1
	}
	// not system admin, user could only see his/her own client users
	if !userDetail.IsSystemAdmin {
		cid = userDetail.Info.Cid
	}

	status, err := c.GetInt("status")
	if err != nil {
		status = -1
	}
	unameStr := c.GetString("unames")
	unameStr = strings.TrimSpace(unameStr)
	unames := make([]string, 0)

	if len(unameStr) > 0 {
		unames = strings.Split(unameStr, ",")
	}

	// none super admin, show only their own company users.
	if !userDetail.IsSystemAdmin {
		cid = userDetail.CInfo.Id
	}
	// search uid by role id
	var uids []int
	uids = nil
	if rid, err := c.GetInt("rid"); err == nil && rid > 0 {
		if roleusers, err := models.GetRolesByRid(rid); err == nil {
			uids = make([]int, 0)
			for _, item := range roleusers {
				uids = append(uids, item.Uid)
			}
		}
	}
	lists, err := models.SearchUsers(unames, uids, cid, status, utype, emailVerifyStatus)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
	} else {
		// add role information
		var resUids []int
		for _, item := range lists {
			resUids = append(resUids, item.Id)
		}
		if roleUsers, err := models.GetAllRolesByUids(resUids); err == nil {
			for index, item := range lists {
				if roles, exist := roleUsers[item.Id]; exist {
					item.Roles = roles
				} else {
					item.Roles = make([]int, 0)
				}
				lists[index] = item
			}
		}
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: lists}
	}
	c.ServeJSON()
}

// admin UploadAvatar ...
// @router /uploadAvatar/ [post]
func (c *UserController) UploadAvatar() {
	// must login and have the Privilege
	c.RequireLogin()
	userInfo := c.CurrentUserDetail.Info
	f, h, err := c.GetFile("avatar")
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "upload file failed"}
		c.ServeJSON()
		return
	} else {
		defer func() {
			_ = f.Close()
		}()
		UploadDir := beego.AppConfig.String("UploadDir")
		if _, err := os.Stat(UploadDir); os.IsNotExist(err) {
			err = os.MkdirAll(UploadDir, 0755)
			if err != nil {
				c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "create upload dir error: " + err.Error()}
				c.ServeJSON()
				return
			}
		}
		fileName := strconv.Itoa(userInfo.Id) + "_" + strconv.FormatInt(time.Now().Unix(), 10) + path.Ext(h.Filename)
		err = c.SaveToFile("avatar", UploadDir+fileName)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
			c.ServeJSON()
			return
		}
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: "/Uploads/" + fileName}
	}
	c.ServeJSON()
}

// admin add new user ...
// @router /saveUser/ [post]
func (c *UserController) SaveUser() {
	userDetail := c.CurrentUserDetail
	cid := userDetail.Info.Cid
	var v models.User
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &v); err == nil {
		v.Ctime = time.Now()
		v.Mtime = time.Now()
		v.Passwd = c.Md5(v.Passwd)
		if v.Cid <= 0 {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "user must belong to a client"}
			c.ServeJSON()
			return
		}
		if len(v.Uname) == 0 {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "user must has a user name"}
			c.ServeJSON()
			return
		}
		if len(v.Passwd) == 0 {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "user must has a initial password"}
			c.ServeJSON()
			return
		}
		if len(v.Email) == 0 {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "user must has a email address"}
			c.ServeJSON()
			return
		}
		if len(v.Phone) == 0 {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "user must has a phone number"}
			c.ServeJSON()
			return
		}
		if len(v.Realname) == 0 {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "user must has a real name"}
			c.ServeJSON()
			return
		}

		// not super admin users and they try to create user for client which is not belong to them.
		if !userDetail.IsSystemAdmin && cid != v.Cid {
			c.Data["json"] = &utils.JSONStruct{Code: 403, Msg: "you don't have right to create new user for this client"}
		} else {
			// check if reach the max user number limitation.
			// check client max user.
			clientInfo, err := models.GetClientById(v.Cid)
			if err != nil {
				c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "find client information error:" + err.Error()}
				c.ServeJSON()
				return
			}
			maxUser := clientInfo.MaxUsers
			// GET current normal users number
			totalUser, err := models.NormalUsersCount(v.Cid)
			if err != nil {
				c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "find client's total user error:" + err.Error()}
				c.ServeJSON()
				return
			}
			if int64(maxUser)-totalUser <= 0 {
				c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "this client has reached the max user number:" + string(maxUser)}
				c.ServeJSON()
				return
			}
			v.Source = constants.SourceAdminAdd // from admin create
			v.Status = constants.TypeNormalUser
			v.EmailVerifyStatus = constants.EmailVerifiedUser
			if _, err := models.CreateUser(&v); err == nil {
				// assign roles to user
				if len(v.Roles) != 0 {
					rids := make([]int, 0)
					for _, tmpRid := range v.Roles {
						if tmpRid > 0 {
							rids = append(rids, tmpRid)
						}
					}
					err = models.ModifyUserRoles(v.Id, rids)
					if err != nil {
						c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "edit user role error, " + err.Error()}
						c.ServeJSON()
						return
					}
					v.Roles = rids
				}
				c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
			} else {
				c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
			}
		}
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// admin EditUser ...
// @router /editUser/:uid [PATCH]
func (c *UserController) EditUser() {
	userDetail := c.CurrentUserDetail
	idStr := c.Ctx.Input.Param(":uid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetUserById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "find edit user error!" + err.Error()}
		c.ServeJSON()
		return
	}
	curCid := userDetail.Info.Cid

	var dat models.User
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &dat)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: err.Error()}
		c.ServeJSON()
		return
	}
	if dat.Cid > 0 {
		// not super admin and the use is not belonging to your client
		if !userDetail.IsSystemAdmin && dat.Cid != curCid {
			c.Data["json"] = &utils.JSONStruct{Code: 403, Msg: "you don't have right to edit user from this client"}
			c.ServeJSON()
			return
		}

		// check if reach the max user number limitation.
		// check client max user.
		clientInfo, err := models.GetClientById(dat.Cid)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "find client information error:" + err.Error()}
			c.ServeJSON()
			return
		}
		maxUser := clientInfo.MaxUsers
		// GET current normal users number
		totalUser, err := models.NormalUsersCount(dat.Cid)
		if err != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "find client's total user error:" + err.Error()}
			c.ServeJSON()
			return
		}
		if int64(maxUser)-totalUser <= 0 {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "this client has reached the max user number:" + string(maxUser)}
			c.ServeJSON()
			return
		}
		v.Cid = dat.Cid
	}
	if dat.Passwd != "" {
		v.Passwd = c.Md5(dat.Passwd)
	}

	if dat.Realname != "" {
		v.Realname = dat.Realname
	}
	if dat.Email != "" {
		v.Email = dat.Email
	}
	if dat.Phone != "" {
		v.Phone = dat.Phone
	}
	if dat.Avatar != "" {
		v.Avatar = dat.Avatar
	}
	if dat.Type >= 0 {
		v.Type = dat.Type
	}
	if dat.EmailVerifyStatus >= 0 {
		v.EmailVerifyStatus = dat.EmailVerifyStatus
	}
	if dat.Status >= 0 {
		v.Status = dat.Status
	}
	v.Roles = dat.Roles
	v.Barcode = dat.Barcode
	// you can not change status here!
	v.Mtime = time.Now()
	if err = models.UpdateUserById(v); err == nil {
		// assign roles to user
		if len(v.Roles) != 0 {
			rids := make([]int, 0)
			for _, tmpRid := range v.Roles {
				if tmpRid > 0 {
					rids = append(rids, tmpRid)
				}
			}
			err = models.ModifyUserRoles(v.Id, rids)
			if err != nil {
				c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "edit user role error, " + err.Error()}
				c.ServeJSON()
				return
			}
			v.Roles = rids
		}
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "edit user information error, " + err.Error()}
	}
	c.ServeJSON()
}

// admin BanUser ...
// @router /banUser/:uid [PATCH]
func (c *UserController) BanUser() {
	userDetail := c.CurrentUserDetail
	idStr := c.Ctx.Input.Param(":uid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetUserById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "find edit user error!" + err.Error()}
		c.ServeJSON()
		return
	}
	// check status
	if v.Status == constants.StatusForbidden {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "this user has already been forbidden!"}
		c.ServeJSON()
		return
	}
	cid := userDetail.Info.Cid
	// not super admin and the use is not belonging to your client
	if !userDetail.IsSystemAdmin && cid != v.Cid {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "you don't have right to edit user from this client"}
		c.ServeJSON()
		return
	}
	v.Status = constants.StatusForbidden // ban user
	v.Mtime = time.Now()
	if err = models.UpdateUserById(v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
	}
	c.ServeJSON()
}

// admin releaseUser ...
// @router /releaseUser/:uid [PATCH]
func (c *UserController) ReleaseUser() {
	userDetail := c.CurrentUserDetail
	idStr := c.Ctx.Input.Param(":uid")
	id, _ := strconv.Atoi(idStr)
	v, err := models.GetUserById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "find edit user error!" + err.Error()}
		c.ServeJSON()
		return
	}
	// check status
	if v.Status == constants.StatusNormal {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "this user is a normal user!"}
		c.ServeJSON()
		return
	}
	cid := userDetail.CInfo.Id
	// not super admin and the use is not belonging to your client
	if !userDetail.IsSystemAdmin && cid != v.Cid {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorForbidden, Msg: "you don't have right to edit user from this client"}
		c.ServeJSON()
		return
	}
	// check if reach the max user number limitation.
	// check client max user.
	clientInfo, err := models.GetClientById(v.Cid)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "find client information error:" + err.Error()}
		c.ServeJSON()
		return
	}
	maxUser := clientInfo.MaxUsers
	// GET current normal users number
	totalUser, err := models.NormalUsersCount(v.Cid)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "find client's total user error:" + err.Error()}
		c.ServeJSON()
		return
	}
	if int64(maxUser)-totalUser <= 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "this client has reached the max user number:" + string(maxUser)}
		c.ServeJSON()
		return
	}
	v.Status = constants.StatusNormal // release user
	v.Mtime = time.Now()
	if err = models.UpdateUserById(v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
	}
	c.ServeJSON()
}

// user register from web page
func (c *UserController) RegisterUser() {
	var v models.User

	var dat map[string]interface{}
	err := json.Unmarshal([]byte(c.Ctx.Input.RequestBody), &dat)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: err.Error()}
		c.ServeJSON()
		return
	}

	paraStr, ok := dat["passwd"]
	passwd := ""
	if ok {
		passwd = paraStr.(string)
		passwordMd5 := []byte(passwd)
		passwd = fmt.Sprintf("%x", md5.Sum(passwordMd5))
	}
	if !ok && len(passwd) == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "user must has a password"}
		c.ServeJSON()
		return
	}
	v.Passwd = passwd

	paraStr, ok = dat["uname"]
	uname := ""
	if ok {
		uname = paraStr.(string)
	}
	if !ok && len(uname) == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "user must has a user name"}
		c.ServeJSON()
		return
	}
	v.Uname = uname

	paraStr, ok = dat["email"]
	email := ""
	if ok {
		email = paraStr.(string)
	}
	if !ok && len(email) == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "user must has a email address"}
		c.ServeJSON()
		return
	}
	v.Email = email

	// search user name and email, check if has been taken
	userInDb, _ := models.GetUserByUname(uname)
	if userInDb != nil && userInDb.Id > 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "user name has been taken"}
		c.ServeJSON()
		return
	}
	userInDb, _ = models.GetUserByEmail(v.Email)
	if userInDb != nil && userInDb.Id > 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorParameter, Msg: "email address has been taken"}
		c.ServeJSON()
		return
	}
	v.Cid = constants.UndefinedId // all register user uid
	v.Ctime = time.Now()
	v.Mtime = time.Now()
	v.EmailVerifyStatus = constants.EmailNotVerifiedUser // didn't verify email
	v.Source = constants.SourceRegister                  // from register
	v.Status = constants.StatusNormal                    // normal user
	v.Type = constants.TypeEmptyUser                     // empty user.
	if _, err := models.CreateUser(&v); err == nil {
		bErr := c.sendVerifyEmail(v.Uname, v.Email)
		if bErr != nil {
			c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: bErr.Message}
			c.ServeJSON()
			return
		}
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: err.Error()}
	}
	c.ServeJSON()
}

func (c *UserController) VerifyEmail() {
	email := c.GetString("mark")
	code := c.GetString("code")
	// set user email verify status
	v, err := models.GetUserByEmail(email)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "find user error." + err.Error()}
		c.ServeJSON()
		return
	}
	// email has been verified or user is normal user, do not need to verify
	if v.EmailVerifyStatus == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
		c.ServeJSON()
		return
	}
	ret, bErr := MWCaptcha.Check(constants.AppName, email, code)
	if bErr != nil {
		errStr := "verify url is invalid or expired(" + strconv.Itoa(bErr.Code) + "), please login and then resend verify email."
		c.Data["json"] = &utils.JSONStruct{Code: bErr.Code, Msg: errStr}
		c.ServeJSON()
		return
	}
	if !ret {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorLogic, Msg: "verify failed, please login and then resend verify email"}
		c.ServeJSON()
		return
	}
	v.EmailVerifyStatus = 0
	err = models.UpdateUserById(v)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorDB, Msg: "update user email status error." + err.Error()}
		c.ServeJSON()
		return
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success", Data: v}
	c.ServeJSON()
}

func (c *UserController) ResendVerifyEmail() {
	userDetail := c.CurrentUserDetail
	if userDetail.Info == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorUnLogin, Msg: "require login"}
		c.DelSession("uname")
		c.DelSession("passwd")
		c.DelSession("uid")
	}
	v := userDetail.Info
	if v.EmailVerifyStatus != 1 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorLogic, Msg: "your email address has been verified."}
		c.ServeJSON()
		return
	}

	d, bErr := MWCaptcha.GetLastCaptcha(constants.AppName, v.Email)
	if bErr != nil && bErr.Code != utils.ErrorNodata {
		c.Data["json"] = &utils.JSONStruct{Code: bErr.Code, Msg: bErr.Message}
		c.ServeJSON()
		return
	}
	// cant resend within 60s
	if bErr == nil && d.Ctime.Add(time.Second*time.Duration(constants.VerifyEmailSendIntervalSeconds)).After(time.Now()) {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorFrequency, Msg: "too frequency, please try later."}
		c.ServeJSON()
		return
	}
	bErr = c.sendVerifyEmail(v.Uname, v.Email)
	if bErr != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: bErr.Message}
		c.ServeJSON()
		return
	}
	c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Msg: "success"}
	c.ServeJSON()
	return
}

func (c *UserController) sendVerifyEmail(uname string, email string) (bErr *utils.BError) {
	code, bErr := MWCaptcha.Insert(constants.AppName, email, "mix", 64, constants.VerifyEmailExpireSeconds)
	if bErr != nil {
		return
	}
	// send verification email
	url := beego.AppConfig.String("AccountActiveUrl") + "?mark=" + email + "&code=" + code + "&app=" + constants.AppName
	_, bErr = MWEmail.ActiveAccount(url, email, uname, constants.AppName)
	return
}

// user edit itself profile
func (c *UserController) EditProfile() {
	userDetail := c.CurrentUserDetail
	id := userDetail.Info.Id
	if id == 0 {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorUnLogin, Msg: "You are not login"}
		c.ServeJSON()
		return
	}
	v, err := models.GetUserById(id)
	if err != nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "find edit user error!" + err.Error()}
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
	if passwd, ok := dat["passwd"].(string); ok {
		if len(passwd) > 0 {
			passwordMd5 := []byte(passwd)
			v.Passwd = fmt.Sprintf("%x", md5.Sum(passwordMd5))
		}
	}
	realname := dat["realname"].(string)
	if len(realname) > 0 {
		v.Realname = realname
	}
	phone := dat["phone"].(string)
	if len(phone) > 0 {
		v.Phone = phone
	}
	avatar := dat["avatar"].(string)
	if len(avatar) > 0 {
		v.Avatar = avatar
	}
	// you can not change status here!
	v.Mtime = time.Now()
	if err = models.UpdateUserById(v); err == nil {
		c.Data["json"] = &utils.JSONStruct{Code: utils.Success, Data: v}
	} else {
		c.Data["json"] = &utils.JSONStruct{Code: utils.ErrorService, Msg: "edit user information error, " + err.Error()}
	}
	c.ServeJSON()
}
