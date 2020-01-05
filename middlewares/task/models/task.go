package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/orm"
	"time"
)

type ServicesTask struct {
	Id       int       `orm:"column(id);auto" json:"id"`
	App      string    `orm:"column(app)" json:"app"`
	Type     int       `orm:"column(type);" json:"type"`
	Data     string    `orm:"column(data);type(json)" json:"data"`
	Ctime    time.Time `orm:"column(ctime);type(datetime)" json:"ctime"`
	Mtime    time.Time `orm:"column(mtime);type(datetime)" json:"mtime"`
	TryTimes int       `orm:"column(try_times)" json:"try_times"`
	Status   int8      `orm:"column(status)" json:"status"`
	Reason   string    `orm:"column(reason)" json:"reason"`
}

func (t *ServicesTask) TableName() string {
	return "services_task"
}

func init() {
	orm.RegisterModel(new(ServicesTask))
}

const (
	TaskTypeEmail = iota + 1
	TaskTypeInitialApp

	MaxTryTime = 3
)

func Insert(qtype int, data string, app string) (res bool, err error) {
	if qtype != TaskTypeEmail && qtype != TaskTypeInitialApp {
		return false, errors.New("invalid task type")
	}
	if len(app) == 0 {
		return false, errors.New("invalid app name")
	}
	v := &ServicesTask{}
	v.Type = qtype
	v.Data = data
	v.Ctime = time.Now()
	v.Mtime = time.Now()
	v.App = app
	o := orm.NewOrm()
	_, err = o.Insert(v)
	if err != nil {
		return false, err
	}
	return true, err
}

func MultiInsert(qtype int, data []string, app string) (err error) {
	if qtype != TaskTypeEmail && qtype != TaskTypeInitialApp {
		return errors.New("invalid task type")
	}
	if len(app) == 0 {
		return errors.New("invalid app name")
	}
	if len(data) == 0 {
		return errors.New("invalid data length")
	}
	insert := make([]*ServicesTask, 0)
	for _, item := range data {
		tmp := &ServicesTask{}
		tmp.Type = qtype
		tmp.Data = item
		tmp.Ctime = time.Now()
		tmp.Mtime = time.Now()
		tmp.App = app
		insert = append(insert, tmp)
	}
	fmt.Printf("%v, %d", insert, len(insert))
	o := orm.NewOrm()
	if _, err = o.InsertMulti(1, insert); err != nil {
		return err
	}
	return nil
}

func one(id int) (v *ServicesTask, err error) {
	o := orm.NewOrm()
	v = &ServicesTask{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

func TaskSuccess(id int) (res bool, err error) {
	v, err := one(id)
	if err != nil {
		return false, err
	}
	v.Status = 1
	o := orm.NewOrm()
	_, err = o.Update(v)
	if err != nil {
		return false, err
	}
	return true, nil
}

func TaskFail(id int, reason string) (res bool, err error) {
	v, err := one(id)
	if err != nil {
		return false, err
	}
	v.TryTimes = v.TryTimes + 1
	if v.TryTimes > MaxTryTime {
		v.Status = 2
	}
	v.Reason = reason
	o := orm.NewOrm()
	_, err = o.Update(v)
	if err != nil {
		return false, err
	}
	return true, nil
}

func UnDoTasks(qtype int, limit int) (l []*ServicesTask, err error) {
	// println("++++", qtype, limit, "------")
	o := orm.NewOrm()
	qs := o.QueryTable(new(ServicesTask))
	_, err = qs.Filter("status", 0).Filter("type", qtype).OrderBy("-ctime").Limit(limit).All(&l)
	if err != nil {
		return nil, err
	}
	// filter max try time tasks
	for _, item := range l {
		if item.TryTimes >= MaxTryTime {
			_, _ = TaskFail(item.Id, "fail too many times!")
		}
	}
	return l, err
}
