package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

//个人资源使用统计
type CloudPersonRecord struct {
	Id             int64
	Username       int64     //使用者
	ResourceDate   time.Time `orm:"type(date)"` //统计日期
	TotalResource  int64     //资源总数
	Team           int64     //所属团队
	UsdeResource   int64     //已经使用的资源
	CanUseResource int64     //未使用的资源
	UsePersent     string    //资源使用率
}

func (r *CloudPersonRecord) TableName() string {
	return "cloud_person_record"
}
func init() {
	orm.RegisterModel(new(CloudPersonRecord))
}

//根据日期和所属团队取得数据
func GetDataByTeam(group int64, resdate time.Time, users []int64) (list []orm.Params, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("cloud_person_record").Filter("Team", group).Filter("ResourceDate", resdate).Filter("Username__in", users).Values(&list)
	return list, err
}

//设置个人资源数量的时候修改数据
func ModifyData(id int64, total int64) error {
	o := orm.NewOrm()
	times := time.Now().Format("2006-01-02")
	_, err := o.QueryTable("cloud_person_record").Filter("Username", id).Filter("ResourceDate", times).Update(orm.Params{
		"TotalResource": total,
	})
	return err
}

//根据日期和id取得数据
func GetDataById(id int64) (list CloudPersonRecord) {
	o := orm.NewOrm()
	var timenow = time.Now().Format("2006-01-02")
	times, _ := time.Parse("2006-01-02", timenow)
	err := o.QueryTable("cloud_person_record").Filter("Username", id).Filter("ResourceDate", times).One(&list)
	if err != nil {
		beego.Error("get person date resource infomation is err", err)
	}
	return list
}

//是否存在该资源
func CheckIsPerson(downtime time.Time, username, groupid int64) bool {
	o := orm.NewOrm()
	flag := o.QueryTable("cloud_person_record").Filter("Username", username).Filter("ResourceDate", downtime).Filter("Team", groupid).Exist()
	return flag
}

//查找资源
func GetInfo(downtime time.Time, username, groupid int64) (list CloudPersonRecord) {
	o := orm.NewOrm()
	err := o.QueryTable("cloud_person_record").Filter("Username", username).Filter("ResourceDate", downtime).Filter("Team", groupid).One(&list)
	beego.Info(err)
	return list
}

//增加个人资源统计
func AddResourcePerson(r *CloudPersonRecord) error {
	o := orm.NewOrm()
	_, err := o.Insert(r)
	return err
}

//修改
func ModifyPersonResouce(downtime time.Time, num int64, id int64) error {
	o := orm.NewOrm()
	persondata := GetDataById(id)
	total := persondata.TotalResource
	used := persondata.UsdeResource + num
	canuse := persondata.CanUseResource - num
	persent := fmt.Sprintf("%.1f", float64(used)/float64(total)*100)
	_, err := o.QueryTable("cloud_person_record").Filter("ResourceDate", downtime).Filter("Username", id).Update(orm.Params{
		"UsdeResource":   used,
		"CanUseResource": canuse,
		"UsePersent":     persent,
	})
	return err
}
