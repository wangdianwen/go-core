package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type StatisticsClientsUsage struct {
	Cid    int64     `orm:"column(cid)" description:"client id"`
	Ctime  time.Time `orm:"column(ctime);type(datetime)" description:"create time"`
	Etime  time.Time `orm:"column(etime);type(datetime)" description:"end time"`
	Id     int       `orm:"column(id);auto"`
	Number int64     `orm:"column(number)" description:"period usage"`
	Stime  time.Time `orm:"column(stime);type(datetime)" description:"start time"`
}

func (t *StatisticsClientsUsage) TableName() string {
	return "core_statistics_clients_usage"
}

func init() {
	orm.RegisterModel(new(StatisticsClientsUsage))
}

// AddStatisticsClientsUsage insert a new StatisticsClientsUsage into database and returns
// last inserted Id on success.
func AddStatisticsClientsUsage(m *StatisticsClientsUsage) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetStatisticsClientsUsageById retrieves StatisticsClientsUsage by Id. Returns error if
// Id doesn't exist
func GetStatisticsClientsUsageById(id int) (v *StatisticsClientsUsage, err error) {
	o := orm.NewOrm()
	v = &StatisticsClientsUsage{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllStatisticsClientsUsage retrieves all StatisticsClientsUsage matches certain condition. Returns empty list if
// no records exist
func GetAllStatisticsClientsUsage(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(StatisticsClientsUsage))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []StatisticsClientsUsage
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateStatisticsClientsUsage updates StatisticsClientsUsage by Id and returns error if
// the record to be updated doesn't exist
func UpdateStatisticsClientsUsageById(m *StatisticsClientsUsage) (err error) {
	o := orm.NewOrm()
	v := StatisticsClientsUsage{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteStatisticsClientsUsage deletes StatisticsClientsUsage by Id and returns error if
// the record to be deleted doesn't exist
func DeleteStatisticsClientsUsage(id int) (err error) {
	o := orm.NewOrm()
	v := StatisticsClientsUsage{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&StatisticsClientsUsage{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
