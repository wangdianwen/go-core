package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

type Client struct {
	AllowUserlog      int8      `orm:"column(allow_userlog)" description:"is allow user log, 0: allow, 1 deny" json:"allow_userlog"`
	BackEndType       int8      `orm:"column(back_end_type);" description:"backendtype: 0 mysql, 1 accredo, 2 advanced" json:"back_end_type"`
	Connection        string    `orm:"column(connection);size(128)" description:"accredo domian of the client" json:"connection"`
	Ctime             time.Time `orm:"column(ctime);type(datetime)" description:"register time" json:"ctime"`
	Description       string    `orm:"column(description)" description:"description of the company" json:"description"`
	Email             string    `orm:"column(email);size(20)" description:"the email address of the client" json:"email"`
	ExpireTime        time.Time `orm:"column(expire_time);type(datetime)" description:"expired day" json:"expire_time"`
	Id                int       `orm:"column(id);auto" json:"id"`
	MaxUsers          int       `orm:"column(max_users)" description:"every client allow max user number" json:"max_users"`
	Mtime             time.Time `orm:"column(mtime);type(datetime)" description:"modify time" json:"mtime"`
	Name              string    `orm:"column(name);size(50)" description:"client name" json:"name"`
	Phone             string    `orm:"column(phone);size(20)" description:"the phone number of the client" json:"phone"`
	Fax               string    `orm:"column(fax)" description:"client fax" json:"fax"`
	AddressLine1      string    `orm:"column(address_line_1)" description:"the address line 1 of the client" json:"address_line_1"`
	AddressLine2      string    `orm:"column(address_line_2)" description:"the address line 2 of the client" json:"address_line_2"`
	City              string    `orm:"column(city)" description:"the city of the client" json:"city"`
	Country           string    `orm:"column(country)" description:"the country of the client" json:"country"`
	PostalCode        string    `orm:"column(postal_code)" description:"the postal code of the client" json:"postal_code"`
	State             string    `orm:"column(state)" description:"the state of the client" json:"state"`
	Web               string    `orm:"column(web)" description:"the web address of the client" json:"web"`
	Gst               string    `orm:"column(gst)" description:"the gst of the client" json:"gst"`
	BankAccount       string    `orm:"column(bank_account)" description:"the bank account of the client" json:"bank_account"`
	Status            int8      `orm:"column(status)" description:"user status 0: normal, 1: forbidden" json:"status"`
	CurrentUserNumber int       `orm:"-" json:"current_user_number"`
}

func (t *Client) TableName() string {
	return "core_client"
}

func init() {
	orm.RegisterModel(new(Client))
}

// CreateClient insert a new client into database and returns
// last inserted Id on success.
func CreateClient(m *Client) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetClientById retrieves Client by Id. Returns error if
// Id doesn't exist
func GetClientById(id int) (v *Client, err error) {
	o := orm.NewOrm()
	v = &Client{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// UpdateClient updates Client by Id and returns error if
// the record to be updated doesn't exist
func UpdateClientById(m *Client) (err error) {
	o := orm.NewOrm()
	v := Client{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteBruntonClient deletes Client by Id and returns error if
// the record to be deleted doesn't exist
func DeleteClientById(id int) (err error) {
	o := orm.NewOrm()
	v := Client{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Client{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// Get all clients list
func GetAllClientsById(cid int) (l []*Client, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Client))
	if cid > 0 {
		qs = qs.Filter("id", cid)
	}
	_, err = qs.All(&l)
	if err != nil {
		return nil, err
	} else {
		return l, nil
	}
}

// get user names by prefix suggestions
func ClientSuggestions(clientName string, limit int) (l []*Client, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Client))
	qs = qs.Filter("name__istartswith", clientName).Limit(limit)
	_, err = qs.All(&l)
	return l, err
}

// admin dashboard search user
func SearchClients(unames []string, status int) (l []*Client, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Client))
	if len(unames) > 0 {
		qs = qs.Filter("name__in", unames)
	}
	// normal
	if status == 0 {
		qs = qs.Filter("status", 0)
		qs = qs.Filter("expire_time__gte", time.Now())
		// forbidden
	} else if status == 1 {
		qs = qs.Filter("status", 1)
	} else if status == 2 {
		qs = qs.Filter("expire_time__lte", time.Now())
	}
	_, err = qs.All(&l)
	return l, err
}

func ClientsRelations(cids []int) (ret map[int]*Client, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Client))
	if len(cids) > 0 {
		qs = qs.Filter("id__in", cids)
	}
	l := make([]*Client, 0)
	_, err = qs.All(&l)
	if err != nil {
		return nil, err
	}
	ret = make(map[int]*Client)
	for _, item := range l {
		ret[item.Id] = item
	}
	return ret, nil
}
