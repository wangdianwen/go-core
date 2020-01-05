package models

import (
	"fmt"
	"go-core/utils"
	"time"

	"github.com/astaxie/beego/orm"
)

type Privilege struct {
	Action     string    `orm:"column(action);size(20)" description:"privilege action name" json:"action"`
	Controller string    `orm:"column(controller);size(20)" description:"privilege controller name" json:"controller"`
	Ctime      time.Time `orm:"column(ctime);type(datetime)" description:"register time" json:"ctime"`
	Id         int       `orm:"column(id);auto" json:"id"`
	Mid        int       `orm:"column(mid);" json:"mid"`
	Mtime      time.Time `orm:"column(mtime);type(datetime)" description:"modify time" json:"mtime"`
	Name       string    `orm:"column(name);size(50)" description:"privilege name" json:"name"`
	Status     int8      `orm:"column(status)" description:"user status 0: normal, 1: forbidden" json:"status"`
}

func (t *Privilege) TableName() string {
	return "core_privilege"
}

func init() {
	orm.RegisterModel(new(Privilege))
}

// CreatePrivilege insert a new Privilege into database and returns
// last inserted Id on success.
func CreatePrivilege(m *Privilege) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetPrivilegeById retrieves Privilege by Id. Returns error if
// Id doesn't exist
func GetPrivilegeById(id int) (v *Privilege, err error) {
	o := orm.NewOrm()
	v = &Privilege{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdatePrivilege updates Privilege by Id and returns error if
// the record to be updated doesn't exist
func UpdatePrivilegeById(m *Privilege) (err error) {
	o := orm.NewOrm()
	v := Privilege{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeletePrivilege deletes Privilege by Id and returns error if
// the record to be deleted doesn't exist
func DeletePrivilege(id int) (err error) {
	o := orm.NewOrm()
	v := Privilege{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Privilege{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// search Roles by ids
func GetPrivilegesByIds(ids []int) (l []*Privilege, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Privilege))
	if len(ids) > 0 {
		qs = qs.Filter("id__in", ids)
	}
	_, err = qs.Filter("status", 0).All(&l)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// search Roles by ids
func GetPrivilegesByNames(names []string) (l []*Privilege, err error) {
	if len(names) == 0 {
		return
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(Privilege))
	_, err = qs.Filter("name__in", names).Filter("status", 0).All(&l)
	if err != nil {
		return nil, err
	}
	return l, nil
}

func GetPrivilegeByCA(controller string, action string) (d *Privilege, err error) {
	o := orm.NewOrm()
	v := &Privilege{Controller: controller, Action: action}
	if err = o.Read(v, "controller", "action"); err == nil {
		return v, nil
	}
	// this is no a real error
	if err.Error() == utils.BeegoNoData {
		return nil, nil
	}
	return nil, err
}

// admin dashboard search Privilege
func SearchPrivileges(name string, status int8, mid int) (l []*Privilege, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Privilege))
	if len(name) > 0 {
		qs = qs.Filter("name__icontains", name)
	}
	// normal
	if status >= 0 {
		qs = qs.Filter("status", status)
	}
	if mid > 0 {
		qs = qs.Filter("mid", mid)
	}
	_, err = qs.All(&l)
	return l, err
}

func CountPrivilegesByMid(mid int) (count int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Privilege))
	if mid > 0 {
		qs = qs.Filter("mid", mid)
	}
	count, err = qs.Count()
	return
}

// get mids => *Privilege slice
func GetAllPrivilegeByMidsGroupByMids(mids []int) (map[int][]*Privilege, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Privilege))
	var l []Privilege
	_, err := qs.Filter("mid__in", mids).All(&l)
	if err != nil {
		return nil, err
	}
	ret := make(map[int][]*Privilege)
	for _, item := range l {
		if lists, exist := ret[item.Mid]; exist {
			lists = append(lists, &item)
			ret[item.Mid] = lists
		} else {
			lists = make([]*Privilege, 0)
			lists = append(lists, &item)
			ret[item.Mid] = lists
		}
	}
	return ret, nil
}

// get mids => *Privilege slice
func GetAllPrivilegesByMids(mids []int) (l []*Privilege, err error) {
	if len(mids) == 0 {
		return nil, nil
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(Privilege))
	_, err = qs.Filter("mid__in", mids).All(&l)
	return
}
