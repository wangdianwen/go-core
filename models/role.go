package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type Role struct {
	Cid        int          `orm:"column(cid)" description:"company id" json:"cid"`
	Ctime      time.Time    `orm:"column(ctime);type(datetime)" description:"create time" json:"ctime"`
	Id         int          `orm:"column(id);auto" json:"id"`
	Mtime      time.Time    `orm:"column(mtime);type(datetime)" description:"modify time" json:"mtime"`
	Rolename   string       `orm:"column(rolename);size(20)" description:"role name" json:"rolename"`
	Status     int8         `orm:"column(status)" description:"user status 0: normal, 1: forbidden" json:"status"`
	Members    int          `orm:"-" json:"members"`
	Menus      []*Menu      `orm:"-" json:"menus"`
	Privileges []*Privilege `orm:"-" json:"privileges"`
}

func (t *Role) TableName() string {
	return "core_role"
}

func init() {
	orm.RegisterModel(new(Role))
}

// CreateRole insert a new Role into database and returns
// last inserted Id on success.
func CreateRole(m *Role) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetRoleById retrieves Role by Id. Returns error if
// Id doesn't exist
func GetRoleById(id int) (Role, error) {
	o := orm.NewOrm()
	v := &Role{Id: id}
	if err := o.Read(v); err == nil {
		return *v, nil
	} else {
		return *v, err
	}
}

// UpdateRole updates Role by Id and returns error if
// the record to be updated doesn't exist
func UpdateRoleById(m *Role) (err error) {
	o := orm.NewOrm()
	v := Role{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteRole deletes Role by Id and returns error if
// the record to be deleted doesn't exist
func DeleteRole(id int) (err error) {
	o := orm.NewOrm()
	v := Role{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Role{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// search Roles by ids
func GetRolesByIds(ids []int) (l []*Role, err error) {
	if len(ids) == 0 {
		return l, errors.New("empty id lists")
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(Role))
	_, err = qs.Filter("id__in", ids).Filter("status", 0).All(&l)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// Get all clients list
func GetAllRoles(cid int) (l []*Role, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Role))
	if cid > 0 {
		qs = qs.Filter("cid", cid)
	}
	_, err = qs.All(&l)
	if err != nil {
		return nil, err
	} else {
		return l, nil
	}
}

// admin dashboard search user
func SearchRoles(rname string, cid int, status int) (l []Role, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Role))
	// cid search
	if cid > 0 {
		qs = qs.Filter("cid", cid)
	}
	// name search
	if len(rname) > 0 {
		qs = qs.Filter("rolename__icontains", rname)
	}
	// status
	if status >= 0 {
		qs = qs.Filter("status", status)
	}
	_, err = qs.All(&l)
	return l, err
}
