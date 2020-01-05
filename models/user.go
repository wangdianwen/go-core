package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Id                int       `orm:"column(id);auto" json:"id"`
	Uname             string    `orm:"column(uname);size(20)" description:"user name for login" json:"uname" valid:"Required;AlphaDash;MinSize(4);MaxSize(128)"`
	Passwd            string    `orm:"column(passwd);size(32)" description:"md5 user password" json:"passwd"`
	Realname          string    `orm:"column(realname);size(20)" description:"user real name" json:"realname"`
	Email             string    `orm:"column(email);size(50)" description:"user email" json:"email" valid:"Required;Email"`
	Phone             string    `orm:"column(phone);size(20)" description:"phone number" json:"phone"`
	Avatar            string    `orm:"column(avatar);size(128)" description:"user avatar" json:"avatar"`
	Cid               int       `orm:"column(cid)" description:"which company the user belong" json:"cid"`
	Barcode           *string   `orm:"column(barcode);null;unique;default(set_null)" description:"user barcode'" json:"barcode"`
	Status            int8      `orm:"column(status)" description:"user status 0: normal, 1: forbidden" json:"status"`
	Type              int8      `orm:"column(type)" description:"user type 0: normal, 1: free trial user, 2: empty user" json:"type"`
	EmailVerifyStatus int8      `orm:"column(email_verify_status)" description:"email verify status: 0 verifyed, 1 not verified" json:"email_verify_status"`
	Source            int8      `orm:"column(source)" description:"where is the user from: 0 from admin create, 1: from user register" json:"source"`
	Ctime             time.Time `orm:"column(ctime);type(datetime)" description:"register time" json:"ctime"`
	Mtime             time.Time `orm:"column(mtime);type(datetime)" description:"modify time" json:"mtime"`
	LastLoginIp       string    `orm:"column(last_login_ip);size(32)" description:"last login ip" json:"last_login_ip" valid:"IP"`
	LastLoginTime     time.Time `orm:"column(last_login_time);type(datetime)" description:"last login time" json:"last_login_time"`
	LastLoginToken    string    `orm:"column(last_login_token);null;size(32);default(null)" description:"last login token" json:"-"`
	Roles             []int     `orm:"-" json:"roles"`
}

func (t *User) TableName() string {
	return "core_user"
}

func init() {
	orm.RegisterModel(new(User))
	// orm.Debug = true
}

// CreateUser insert a new User into database and returns
// last inserted Id on success.
func CreateUser(m *User) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserById retrieves User by Id. Returns error if
// Id doesn't exist
func GetUserById(id int) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetUserById retrieves User by Uname and password. Returns error if
// uname and password don't exist
func GetUserLogin(uname string, password string) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Uname: uname, Passwd: password}
	if err = o.Read(v, "uname", "passwd"); err == nil {
		return v, nil
	}
	return nil, err
}

func GetBarcodeLogin(barcode *string) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Barcode: barcode}
	if err = o.Read(v, "barcode"); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateUser updates User by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserById(m *User) (err error) {
	o := orm.NewOrm()
	v := User{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// get user names by prefix suggestions
func UserSuggestions(userName string, cid int, limit int) (l []*User, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(User))
	qs = qs.Filter("uname__istartswith", userName).Limit(limit)
	if cid > 0 {
		qs = qs.Filter("cid", cid)
	}
	_, err = qs.All(&l)
	return l, err
}

// for only data provider
func GetRawUserByCid(cid int) (l []*User, err error) {
	o := orm.NewOrm()
	_, err = o.Raw("Select * from core_user where cid = ?", cid).QueryRows(&l)
	return
}

// admin dashboard search user
func SearchUsers(unames []string, uids []int, cid int, status int, utype int, emailVerifyStatus int) (l []*User, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(User))
	if uids != nil && len(uids) == 0 {
		return nil, errors.New("no data")
	} else if uids != nil {
		qs = qs.Filter("id__in", uids)
	}
	if len(unames) > 0 {
		qs = qs.Filter("uname__in", unames)
	}
	if cid > 0 {
		qs = qs.Filter("cid", cid)
	}
	if status >= 0 {
		qs = qs.Filter("status", status)
	}
	if utype >= 0 {
		qs = qs.Filter("type", utype)
	}
	if emailVerifyStatus >= 0 {
		qs = qs.Filter("email_verify_status", emailVerifyStatus)
	}
	_, err = qs.All(&l)
	return l, err
}

// admin normal user count
func NormalUsersCount(cid int) (total int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(User))
	qs = qs.Filter("cid", cid).Filter("status", 0)
	total, err = qs.Count()
	return total, err
}

// get normal user count by client
func CountByClient(cids []int) (map[int]int, error) {
	if len(cids) == 0 {
		return nil, nil
	}
	o := orm.NewOrm()
	sql := "SELECT COUNT(*) AS `TOTAL`, `cid` FROM `core_user` WHERE `status` = 0 AND `cid` IN ( "
	for range cids {
		sql += "? ,"
	}
	sql = strings.Trim(sql, ",") + " )  GROUP BY `cid`"
	var maps []orm.Params
	num, err := o.Raw(sql, cids).Values(&maps)
	if err != nil || num == 0 {
		return nil, err
	}
	l := make(map[int]int)
	for _, item := range maps {
		itemCid, _ := strconv.Atoi(item["cid"].(string))
		total, _ := strconv.Atoi(item["TOTAL"].(string))
		l[itemCid] = total
	}
	return l, nil
}

// find user by unique uname
func GetUserByUname(uname string) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Uname: uname}
	if err = o.Read(v, "uname"); err == nil {
		return v, nil
	}
	return nil, err
}

func GetUserByEmail(email string) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Email: email}
	if err = o.Read(v, "email"); err == nil {
		return v, nil
	}
	return nil, err
}
