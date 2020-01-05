package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type RolePrivilege struct {
	Ctime time.Time `orm:"column(ctime);type(datetime)" description:"register time" json:"ctime"`
	Id    int       `orm:"column(id);auto" json:"id"`
	Pid   int       `orm:"column(pid)" description:"privilege id" json:"pid"`
	Rid   int       `orm:"column(rid)" description:"role id" json:"rid"`
	Mid   int       `orm:"column(mid)" description:"menu id" json:"mid"`
}

func (t *RolePrivilege) TableName() string {
	return "core_role_privilege"
}

func init() {
	orm.RegisterModel(new(RolePrivilege))
}

func GetAllRolePrivilege() (ml []*RolePrivilege, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RolePrivilege))
	if _, err := qs.All(&ml); err == nil {
		return ml, nil
	}
	return nil, err
}

// Get Privilege Role relation by role ids
func GetRolePrivilegeByRids(rids []int) (ml []*RolePrivilege, err error) {
	if len(rids) == 0 {
		return nil, nil
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(RolePrivilege))
	if _, err := qs.Filter("rid__in", rids).All(&ml); err == nil {
		return ml, nil
	}
	return nil, err
}

func DelRolePrivilegeByMid(mid int) (rows int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RolePrivilege))
	rows, err = qs.Filter("mid", mid).Delete()
	return
}

func GetAllMenusByRids(rids []int) (ml []*RolePrivilege, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RolePrivilege))
	if len(rids) == 0 {
		return nil, nil
	}
	qs = qs.Filter("rid__in", rids).Filter("mid__gt", 0)
	if _, err := qs.All(&ml); err == nil {
		return ml, err
	}
	return nil, err
}

func GetMenuPrivilegeByRid(rid int) (ml []*RolePrivilege, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RolePrivilege))
	qs = qs.Filter("rid", rid)
	if _, err := qs.All(&ml); err == nil {
		return ml, err
	}
	return nil, err
}

// modify user roles
func ModifyRolesPrivileges(rid int, mNames []string, pNames []string) (err error) {
	// filter invalid mids
	menus, err := GetAllMenuByNames(mNames)
	if err != nil {
		return err
	}
	mids := make([]int, 0)
	for _, item := range menus {
		mids = append(mids, item.Id)
	}
	// filter invalid pnames
	privileges, err := GetPrivilegesByNames(pNames)
	if err != nil {
		return err
	}
	pids := make([]int, 0)
	for _, item := range privileges {
		pids = append(pids, item.Id)
	}
	// get all current user roles
	l, err := GetMenuPrivilegeByRid(rid)
	if err != nil {
		return err
	}
	existsMids := make(map[int]int) // for delete none exist roles
	existsPids := make(map[int]int)
	for _, item := range l {
		if item.Mid > 0 {
			existsMids[item.Mid] = item.Id
		}
		if item.Pid > 0 {
			existsPids[item.Pid] = item.Id
		}
	}

	newRelations := make([]RolePrivilege, 0) // for insert
	// deal with menu data
	for _, mid := range mids {
		if _, exist := existsMids[mid]; exist {
			// the rest of existsMids is no longer existing in new submitted roles list, we need to delete them
			delete(existsMids, mid)
		} else {
			// not exist in new rids
			relation := RolePrivilege{Ctime: time.Now(), Rid: rid, Mid: mid, Pid: 0}
			newRelations = append(newRelations, relation)
		}
	}
	// deal with Privilege data
	for _, pid := range pids {
		if _, exist := existsPids[pid]; exist {
			// the rest of existsPids is no longer existing in new submitted roles list, we need to delete them
			delete(existsPids, pid)
		} else {
			// not exist in new rids
			relation := RolePrivilege{Ctime: time.Now(), Rid: rid, Mid: 0, Pid: pid}
			newRelations = append(newRelations, relation)
		}
	}
	o := orm.NewOrm()
	if err = o.Begin(); err != nil {
		return err
	}

	// has delete mids, delete
	delMids := make([]int, 0)
	for delmid := range existsMids {
		delMids = append(delMids, delmid)
	}
	if len(delMids) > 0 {
		if _, err = o.QueryTable(new(RolePrivilege)).Filter("rid", rid).Filter("mid__in", delMids).Delete(); err != nil {
			_ = o.Rollback()
			return err
		}
	}
	// has delete pids, delete
	delPids := make([]int, 0)
	for delpid := range existsPids {
		delPids = append(delPids, delpid)
	}
	if len(delPids) > 0 {
		if _, err = o.QueryTable(new(RolePrivilege)).Filter("rid", rid).Filter("pid__in", delPids).Delete(); err != nil {
			_ = o.Rollback()
			return err
		}
	}
	// has new rids, insert
	if len(newRelations) > 0 {
		if _, err = o.InsertMulti(1, newRelations); err != nil {
			_ = o.Rollback()
			return nil
		}
	}

	//println("============")
	//fmt.Printf("%v", delMids)
	//fmt.Printf("%v", delPids)
	//fmt.Printf("%v", newRelations)
	//println("============")

	_ = o.Commit()
	return nil
}
