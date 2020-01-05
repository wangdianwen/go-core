package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type CoreApps struct {
	Id             int       `orm:"column(id);auto" json:"id"`
	Name           string    `orm:"column(name);size(128)" description:"application name" json:"name"`
	AppKey         string    `orm:"column(app_key);size(128)" description:"application key" json:"app_key"`
	AppSecret      string    `orm:"column(app_secret);size(128)" description:"application secret" json:"app_secret"`
	Ctime          time.Time `orm:"column(ctime);type(datetime)" description:"register time" json:"ctime"`
	Mtime          time.Time `orm:"column(mtime);type(datetime)" description:"modify time" json:"mtime"`
	Status         int8      `orm:"column(status)" description:"application status 0: normal, 1: forbidden" json:"status"`
	AllowFreeTrial int8      `orm:"column(allow_free_trial)" description:"application allow free trial, 0: yes, 1: no" json:"allow_free_trial"`
}

func (t *CoreApps) TableName() string {
	return "core_apps"
}

func init() {
	orm.RegisterModel(new(CoreApps))
}

func GetAllAppsList(status int8) (l []*CoreApps, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CoreApps))
	if status >= 0 {
		qs = qs.Filter("status", status)
	}
	_, err = qs.All(&l)
	return l, err
}

func GetAllValidateAppsList() (l []*CoreApps, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CoreApps)).Filter("status", 0).Filter("allow_free_trial", 0)
	_, err = qs.All(&l)
	return l, err
}

// get user names by prefix suggestions
func AppsList(appName string, limit int) (l []*CoreApps, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(CoreApps))
	if len(appName) > 0 {
		qs = qs.Filter("name__istartswith", appName)
	}
	if limit > 0 {
		qs = qs.Limit(limit)
	}
	_, err = qs.All(&l)
	return l, err
}

// get user by id
func GetAppById(id int) (v *CoreApps, err error) {
	o := orm.NewOrm()
	v = &CoreApps{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func GetAppsByIds(ids []int) (l []*CoreApps, err error) {
	if len(ids) == 0 {
		return nil, nil
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(CoreApps))
	qs = qs.Filter("id__in", ids)
	_, err = qs.All(&l)
	return
}

// add apps
func CreateApp(m *CoreApps) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// update apps
func UpdateAppById(m *CoreApps) (err error) {
	o := orm.NewOrm()
	v := CoreApps{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}
