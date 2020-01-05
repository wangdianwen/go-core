package models

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type RoleUser struct {
	Ctime time.Time `orm:"column(ctime);type(datetime)" description:"register time"`
	Id    int       `orm:"column(id);auto"`
	Rid   int       `orm:"column(rid)" description:"role id"`
	Uid   int       `orm:"column(uid)" description:"user id"`
}

func (t *RoleUser) TableName() string {
	return "core_role_user"
}

func init() {
	orm.RegisterModel(new(RoleUser))
}

func GetAllRoleUser() (l []*RoleUser, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RoleUser))
	if _, err := qs.All(&l); err == nil {
		return l, nil
	}
	return nil, err
}

func GetRolesByRid(rid int) (l []*RoleUser, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RoleUser))
	if _, err := qs.Filter("rid", rid).All(&l); err == nil {
		return l, nil
	}
	return nil, err
}

// get all role ids by uid
func GetValidRoles(uid int) (l []*RoleUser, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(RoleUser))
	if _, err := qs.Filter("uid", uid).All(&l); err == nil {
		return l, nil
	}
	return nil, err
}

// get role user number by role id
func CountByRids(rids []int) (map[int]int, error) {
	if len(rids) == 0 {
		return nil, nil
	}
	o := orm.NewOrm()
	sql := "SELECT COUNT(*) AS `TOTAL`, `rid` FROM `core_role_user` WHERE `rid` IN ( "
	for range rids {
		sql += "? ,"
	}
	sql = strings.Trim(sql, ",") + " ) GROUP BY `rid`"
	var maps []orm.Params
	num, err := o.Raw(sql, rids).Values(&maps)
	if err != nil || num == 0 {
		return nil, err
	}
	l := make(map[int]int)
	for _, item := range maps {
		itemRid, _ := strconv.Atoi(item["rid"].(string))
		total, _ := strconv.Atoi(item["TOTAL"].(string))
		l[itemRid] = total
	}
	return l, nil
}

// get uid => rids slice
func GetAllRolesByUids(uids []int) (map[int][]int, error) {
	if len(uids) == 0 {
		return nil, nil
	}
	o := orm.NewOrm()
	qs := o.QueryTable(new(RoleUser))
	var l []RoleUser
	_, err := qs.Filter("uid__in", uids).All(&l)
	if err != nil {
		return nil, err
	}
	ret := make(map[int][]int)
	for _, item := range l {
		if rids, exist := ret[item.Uid]; exist {
			rids = append(rids, item.Rid)
			ret[item.Uid] = rids
		} else {
			rids = make([]int, 0)
			rids = append(rids, item.Rid)
			ret[item.Uid] = rids
		}
	}
	return ret, nil
}

// modify user roles
func ModifyUserRoles(uid int, rids []int) (err error) {
	// filter invalid rids
	roles, err := GetRolesByIds(rids)
	if err != nil {
		return err
	}
	rids = nil
	for _, item := range roles {
		// check role status, you can not change to a forbidden role
		if item.Status != 0 {
			continue
		}
		rids = append(rids, item.Id)
	}
	// no validate id， return
	if len(rids) == 0 {
		return errors.New("no validated rid")
	}
	// get all current user roles
	l, err := GetValidRoles(uid)
	if err != nil {
		return err
	}
	existsRids := make(map[int]int) // for delete none exist roles, roleid =>
	for _, item := range l {
		existsRids[item.Rid] = item.Id
	}
	newRelations := make([]RoleUser, 0) // for insert
	for _, rid := range rids {
		if _, exist := existsRids[rid]; exist {
			// the rest of existsRids is no longer existing in new submitted roles list, we need to delete them
			delete(existsRids, rid)
		} else {
			// not exist in new rids
			relation := RoleUser{Ctime: time.Now(), Uid: uid, Rid: rid}
			newRelations = append(newRelations, relation)
		}
	}
	o := orm.NewOrm()
	if err = o.Begin(); err != nil {
		return err
	}

	// fmt.Printf("%v , %v " , newRelations, existsRids)
	// has new rids, insert
	if len(newRelations) > 0 {
		if _, err = o.InsertMulti(1, newRelations); err != nil {
			_ = o.Rollback()
			return err
		}
	}
	// has delete rids, delete
	if len(existsRids) > 0 {
		rids := make([]int, 0)
		for rid := range existsRids {
			rids = append(rids, rid)
		}
		if _, err = o.QueryTable(new(RoleUser)).Filter("uid", uid).Filter("rid__in", rids).Delete(); err != nil {
			_ = o.Rollback()
			return err
		}
	}
	_ = o.Commit()
	return nil
}
