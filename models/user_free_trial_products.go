package models

import (
	"encoding/json"
	"github.com/astaxie/beego/orm"
	"github.com/wangdianwen/go-core/configs/constants"
	task "github.com/wangdianwen/go-core/middlewares/task/models"
	"github.com/wangdianwen/go-core/utils"
	"math"
	"math/rand"
	"strconv"
	"time"
)

type UserFreeTrialProducts struct {
	Id    int `orm:"column(id);auto" json:"id"`
	Appid int `orm:"column(appid)" description:"application id" json:"appid"`
	Uid   int `orm:"column(uid)" description:"user id" json:"uid"`
}

func (t *UserFreeTrialProducts) TableName() string {
	return "core_user_free_trial_products"
}

func init() {
	orm.RegisterModel(new(UserFreeTrialProducts))
	// orm.Debug = true
}

func GetAllUserFreeTrailProducts() (l []*UserFreeTrialProducts, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserFreeTrialProducts))
	_, err = qs.All(&l)
	return
}

func GetUserFreeTrialProductsByUid(uid int) (l []*UserFreeTrialProducts, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UserFreeTrialProducts))
	_, err = qs.Filter("uid", uid).All(&l)
	return
}

func StartFreeTrial(uid int, appIds []int) (bErr *utils.BError, products []*CoreApps) {
	o := orm.NewOrm()
	if err := o.Begin(); err != nil {
		return &utils.BError{Code: utils.ErrorDB, Message: "start transaction error", Error: err}, nil
	}
	user, err := GetUserById(uid)
	if err != nil && err.Error() == utils.BeegoNoData {
		_ = o.Rollback()
		return &utils.BError{Code: utils.ErrorNodata, Message: "can't find user information", Error: err}, nil
	} else if err != nil {
		_ = o.Rollback()
		return &utils.BError{Code: utils.ErrorNodata, Message: "find user information error", Error: err}, nil
	}
	if user.EmailVerifyStatus == constants.EmailNotVerifiedUser {
		return &utils.BError{Code: utils.ErrorLogic, Message: "please verify your email first.", Error: err}, nil
	}
	if user.Type != constants.TypeEmptyUser && user.Type != constants.TypeFreeTrialUser {
		return &utils.BError{Code: utils.ErrorLogic, Message: "you are not eligible for free trial.", Error: err}, nil
	}
	if len(appIds) == 0 {
		return &utils.BError{Code: utils.ErrorParameter, Message: "please select at least one products.", Error: err}, nil
	}
	user.Type = constants.TypeFreeTrialUser
	user.Mtime = time.Now()
	if user.Cid == -1 {
		user.Cid = -(rand.Intn(100)+10)*int(math.Pow10(len(strconv.Itoa(uid)))) - uid // random client id: random(10 - 99) + uid
		// while user first time to use, reset its ctime to now
		// because we use ctime to judge if user is free trial expired
		user.Ctime = time.Now()
	}

	err = UpdateUserById(user)
	if err != nil {
		_ = o.Rollback()
		return &utils.BError{Code: utils.ErrorDB, Message: "can't find user information", Error: err}, nil
	}

	list, err := GetUserFreeTrialProductsByUid(uid)
	if err != nil && err.Error() == utils.BeegoNoData {
		// do nothing
	} else if err != nil {
		_ = o.Rollback()
		return &utils.BError{Code: utils.ErrorDB, Message: "find user information error", Error: err}, nil
	}
	appRelations := make(map[int]bool)
	allAppIds := make([]int, 0)
	// find existing products
	for _, item := range list {
		appRelations[item.Appid] = true
		allAppIds = append(allAppIds, item.Appid)
	}
	appForInsertIds := make([]int, 0)
	for _, item := range appIds {
		if _, ok := appRelations[item]; ok {
			continue
		}
		appForInsertIds = append(appForInsertIds, item)
		allAppIds = append(allAppIds, item)
	}
	// insert into products lists
	newFreeTrials := make([]*UserFreeTrialProducts, 0)
	for _, item := range appForInsertIds {
		tmp := &UserFreeTrialProducts{Appid: item, Uid: uid}
		newFreeTrials = append(newFreeTrials, tmp)
	}
	if len(newFreeTrials) > 0 {
		if _, err = o.InsertMulti(1, newFreeTrials); err != nil {
			_ = o.Rollback()
			return &utils.BError{Code: utils.ErrorDB, Message: "insert product free trial list error", Error: err}, nil
		}
	}
	tmp, err := GetAppsByIds(allAppIds)
	if err != nil {
		_ = o.Rollback()
		return &utils.BError{Code: utils.ErrorDB, Message: "get products details fail.", Error: err}, nil
	}
	apps := make([]*CoreApps, 0)
	taskData := make([]string, 0)
	for _, item := range tmp {
		item.AppSecret = "***"
		item.AppKey = "***"
		apps = append(apps, item)

		dataTmp := make(map[string]interface{})
		dataTmp["appName"] = item.Name
		dataTmp["uid"] = uid
		dataJson, err := json.Marshal(dataTmp)
		if err != nil {
			_ = o.Rollback()
			return &utils.BError{Code: utils.ErrorParseJson, Message: "parse data to json string error", Error: err}, nil
		}
		taskData = append(taskData, string(dataJson))
	}
	// insert to task
	err = task.MultiInsert(task.TaskTypeInitialApp, taskData, constants.AppName)
	if err != nil {
		_ = o.Rollback()
		return &utils.BError{Code: utils.ErrorDB, Message: "submit task error", Error: err}, nil
	}
	err = o.Commit()
	if err != nil {
		return &utils.BError{Code: utils.ErrorDB, Message: "commit all sql failed.", Error: err}, nil
	}
	return nil, apps
}
