package models

import (
	"github.com/astaxie/beego/orm"
)

type SysList struct {
	Id        int64 //主键
	Uid       int64 `orm:"index"` //云KUid
	Token     string
	Ip        string
	MgrDevice []*MgrDevice `orm:"reverse(many)" json:"-"` // `orm:"reverse(one)" json:"-"`
}

func init() {
	orm.RegisterModel(new(SysList))
}

// add MgrGroup
func AddSysList(syslist *SysList) (num int64, err error) {
	o := orm.NewOrm()
	num, err = o.Insert(syslist)
	return num, err
}

// SysList
func GetSysList(page int64, page_size int64, sort string) (syslist []orm.Params, count int64) {
	o := orm.NewOrm()
	role := new(Role)
	qs := o.QueryTable(role)
	qs.Limit(page_size, page).OrderBy(sort).Values(&syslist)
	count, _ = qs.Count()
	return syslist, count
}

// get the group by id result one
func GetSysListByIdOne(id int64) (syslist SysList, err error) {
	err = syslistQuerySeter().Filter("Id", id).Limit(1).One(&syslist)
	return syslist, err
}

// 根据群控账户查询
func GetSysListByUid(uid int64) (syslist SysList, err error) {
	err = syslistQuerySeter().Filter("Uid", uid).Limit(1).One(&syslist)
	return syslist, err
}

// 根据pn查询
func GetSyslistByUserToken(username, token string) (syslist SysList, userInfo User, err error) {
	err = userQuerySeter().Filter("Username", username).Limit(1).RelatedSel().One(&userInfo)
	if err != nil {
		return
	} else {
		err = syslistQuerySeter().Filter("Uid", userInfo.Id).Filter("Token", token).Limit(1).One(&syslist)
		return
	}
}

func syslistQuerySeter() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(SysList))
}
