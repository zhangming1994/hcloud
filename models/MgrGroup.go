package models

import (
	"github.com/astaxie/beego/orm"
)

type MgrGroup struct {
	Id         int64  //主键
	Name       string `orm:"unique"` //组别名称
	DevieCount int64  `orm:"-"`
}

func init() {
	orm.RegisterModel(new(MgrGroup))
}

// add MgrGroup
func AddMgrGroup(MgrGroup *MgrGroup) (num int64, err error) {
	o := orm.NewOrm()
	num, err = o.Insert(MgrGroup)
	return num, err
}

// get the group by id result one
func GetGroupByIdOne(id int64) (mgrgroup MgrGroup, err error) {
	err = mgrgroupQuerySeter().Filter("Id", id).Limit(1).One(&mgrgroup)
	return mgrgroup, err
}

func mgrgroupQuerySeter() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(MgrGroup))
}
