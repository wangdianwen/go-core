package models

import (
	"github.com/astaxie/beego/orm"
	"go-core/utils"
	"math/rand"
	"time"
)

type ServicesCaptcha struct {
	Id       int       `orm:"column(id);auto" json:"id"`
	Code     string    `orm:"column(code);" json:"code"`
	Ctime    time.Time `orm:"column(ctime);type(datetime)" json:"ctime"`
	Mtime    time.Time `orm:"column(mtime);type(datetime)" json:"mtime"`
	Expire   time.Time `orm:"column(expire);type(datetime)" json:"expire"`
	Verify   time.Time `orm:"column(verify);type(datetime)" json:"verify"`
	TryTimes int       `orm:"column(try_times)" json:"try_times"`
	Status   int8      `orm:"column(status)" json:"status"`
	Mark     string    `orm:"column(mark)" json:"mark"`
	App      string    `orm:"column(app)" json:"app"`
}

func (t *ServicesCaptcha) TableName() string {
	return "services_captcha"
}

func init() {
	orm.RegisterModel(new(ServicesCaptcha))
}

// admin normal user count
func Insert(app string, mark string, codeType string, size int, expire int) (code string, bErr *utils.BError) {
	o := orm.NewOrm()
	if len(mark) == 0 {
		return "", &utils.BError{Message: "invalid mark", Code: utils.ErrorParameter}
	}
	if size <= 0 {
		return "", &utils.BError{Message: "invalid size", Code: utils.ErrorParameter}
	}
	if expire <= 0 {
		expire = 600
	}
	qs := o.QueryTable(new(ServicesCaptcha))
	if err := o.Begin(); err != nil {
		return "", &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorService}
	}

	// delete all other exist code
	update := make(map[string]interface{})
	update["status"] = 1
	_, err := qs.Filter("app", app).Filter("mark", mark).Filter("status", 0).Update(update)
	if err != nil {
		_ = o.Rollback()
		return "", &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorDB}
	}
	code = generateCode(codeType, size)
	if len(code) == 0 {
		return "", &utils.BError{Message: "generate code failed", Code: utils.ErrorService}
	}
	v := &ServicesCaptcha{}
	v.Code = code
	v.Ctime = time.Now()
	v.Mtime = time.Now()
	v.Expire = time.Now().Local().Add(time.Second * time.Duration(expire))
	v.Mark = mark
	v.App = app
	_, err = o.Insert(v)
	if err != nil {
		_ = o.Rollback()
		return "", &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorDB}
	}
	_ = o.Commit()
	return code, nil
}

// check captcha
func Check(app string, mark string, code string) (res bool, bErr *utils.BError) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ServicesCaptcha))
	v := ServicesCaptcha{}
	err := qs.Filter("app", app).Filter("mark", mark).Filter("code", code).Filter("status", 0).OrderBy("-ctime").Limit(1).One(&v)
	if err != nil {
		if err.Error() == utils.BeegoNoData {
			return false, &utils.BError{Message: "nodata", Code: utils.ErrorNodata} // non exist code
		}
		return false, &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorDB}
	}
	if v.Code != code {
		v.TryTimes = v.TryTimes + 1
		// checked failed more than 10 times, automatically close this captcha
		if v.TryTimes > 10 {
			v.Status = 2
		}
		_, err = o.Update(v)
		if err != nil {
			return false, &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorDB}
		}
		return false, nil
	} else {
		v.TryTimes = v.TryTimes + 1
		v.Status = 1
		_, err = o.Update(&v)
		if err != nil {
			return false, &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorDB}
		}
		if v.Expire.Before(time.Now()) {
			return false, &utils.BError{Message: "code has expired", Code: utils.ErrorExpire} // expire code
		}
	}
	return true, nil
}

// get last captcha
func GetLastCaptcha(app string, mark string) (res *ServicesCaptcha, bErr *utils.BError) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ServicesCaptcha))
	v := ServicesCaptcha{}
	err := qs.Filter("app", app).Filter("mark", mark).OrderBy("-ctime").Limit(1).One(&v)
	if err != nil {
		if err.Error() == utils.BeegoNoData {
			return nil, &utils.BError{Message: "nodata", Code: utils.ErrorNodata} // non exist code
		}
		return nil, &utils.BError{Error: err, Message: err.Error(), Code: utils.ErrorDB}
	}
	return &v, nil
}

//
func generateCode(codeType string, size int) (ret string) {
	characters := make(map[string][]rune)
	characters["number"] = []rune("0123456789")
	characters["letter"] = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	characters["mix"] = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	if randomSample, ok := characters[codeType]; ok {
		b := make([]rune, size)
		for i := range b {
			b[i] = randomSample[rand.Intn(len(randomSample))]
		}
		ret = string(b)
	} else {
		ret = ""
	}
	return ret
}
