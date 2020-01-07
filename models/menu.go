package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type Menu struct {
	Id         int          `orm:"column(id);auto" json:"id"`
	Appid      int          `orm:"column(appid)" description:"application id" json:"appid"`
	Paid       int          `orm:"column(paid)" description:"parent menu id" json:"paid"`
	Name       string       `orm:"column(name);size(128)" description:"menu name" json:"name"`
	Url        string       `orm:"column(url);size(64)" description:"menu url" json:"url"`
	Visible    int8         `orm:"column(visible)" description:"is visible;0: visile, 1 hide" json:"visible"`
	Ctime      time.Time    `orm:"column(ctime);type(datetime)" description:"create time" json:"ctime"`
	Mtime      time.Time    `orm:"column(mtime);type(datetime)" description:"modify time" json:"mtime"`
	Level      int          `orm:"-" json:"level"`
	Children   []*Menu      `orm:"-" json:"children"`
	Privileges []*Privilege `orm:"-" json:"privileges"`
}

func (t *Menu) TableName() string {
	return "core_menu"
}

func init() {
	orm.RegisterModel(new(Menu))
	// orm.Debug = true
}

// CreateMenu insert a new Menu into database and returns
// last inserted Id on success.
func CreateMenu(m *Menu) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetMenuById retrieves menu by Id. Returns error if
// Id doesn't exist
func GetMenuById(id int) (v *Menu, err error) {
	o := orm.NewOrm()
	v = &Menu{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateMenu updates Menu by Id and returns error if
// the record to be updated doesn't exist
func UpdateMenuById(m *Menu) (err error) {
	o := orm.NewOrm()
	v := Menu{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteMenu deletes Menu by Id and returns error if
// the record to be deleted doesn't exist
func DeleteMenu(id int) (err error) {
	o := orm.NewOrm()
	v := Menu{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Menu{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// search all menu
func SearchMenu(paid int, name string) (l []*Menu, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Menu))
	if len(name) > 0 {
		qs = qs.Filter("name__icontains", name)
	}
	if paid >= 0 {
		qs = qs.Filter("paid", paid)
	}
	_, err = qs.All(&l)
	return l, err
}

func CountMenu(paid int, name string) (total int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Menu))
	if len(name) > 0 {
		qs = qs.Filter("name__icontains", name)
	}
	if paid >= 0 {
		qs = qs.Filter("paid", paid)
	}
	return qs.Count()
}

// get paid => *Menu slice
func GetAllMenusByPaids(paids []int) (map[int][]*Menu, error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Menu))
	var l []Menu
	_, err := qs.Filter("paid__in", paids).All(&l)
	if err != nil {
		return nil, err
	}
	ret := make(map[int][]*Menu)
	for _, item := range l {
		if lists, exist := ret[item.Paid]; exist {
			lists = append(lists, &item)
			ret[item.Paid] = lists
		} else {
			lists = make([]*Menu, 0)
			lists = append(lists, &item)
			ret[item.Paid] = lists
		}
	}
	return ret, nil
}

// get all menus by ids
func GetAllMenuByIds(ids []int) (l []*Menu, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Menu))
	if len(ids) > 0 {
		qs = qs.Filter("id__in", ids).OrderBy("paid")
	}
	_, err = qs.All(&l)
	return
}

// get menu list by app id
func GetAllMenuByAppids(appIds []int) (l []*Menu, err error) {
	if len(appIds) == 0 {
		return nil, nil
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(Menu))
	qs = qs.Filter("appid__in", appIds).OrderBy("paid")
	_, err = qs.All(&l)
	return
}

// get all menus by ids and appids
func GetAllMenuByMidsAndAppids(ids []int, appids []int) (l []*Menu, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Menu))
	if len(ids) > 0 {
		qs = qs.Filter("id__in", ids)
	}
	if len(appids) > 0 {
		qs = qs.Filter("appid__in", appids)
	}
	qs = qs.OrderBy("paid")
	_, err = qs.All(&l)
	return
}

// get all menus by ids
func GetAllMenuByNames(names []string) (l []*Menu, err error) {
	if len(names) == 0 {
		return
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(Menu))
	if len(names) > 0 {
		qs = qs.Filter("name__in", names).OrderBy("paid")
	}
	_, err = qs.All(&l)
	return
}
