package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

//请求资源记录表
type CloudDistributRecord struct {
	Id                 int64
	DownNumber         int64     //分发数量
	DistributName      int64     //分配给的人
	DistriResourceType int64     //请求资源类型
	DownTime           time.Time `orm:type(date)` //分发时间(2006-01-02)
}

func (r *CloudDistributRecord) TableName() string {
	return "cloud_distribut_record"
}
func init() {
	orm.RegisterModel(new(CloudDistributRecord))
}

//根据时间和人和资源类型来判断该条记录是否存在
func CheckIsExistRecord(name string, downtime time.Time, restype int64) bool {
	o := orm.NewOrm()
	flag := o.QueryTable("cloud_distribut_record").Filter("DistriUser", name).Filter("DownTime", downtime).Filter("DistriResourceType", restype).Exist()
	return flag
}

//分发记录添加
func AddDistributRecord(m *CloudDistributRecord) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Insert(m)
	return num, err
}

//根据DistributName查找已经分发数量
func GetGroupTotalDown(distriname int64) (num int64, err error) {
	o := orm.NewOrm()
	var list []orm.Params
	times := time.Now().Format("2006-01-02")
	downtime, _ := time.Parse("2006-01-02", times)
	_, err = o.QueryTable("cloud_distribut_record").Filter("DistributName", distriname).Filter("DownTime", downtime).Values(&list)
	for i := 0; i < len(list); i++ {
		onerecord := list[i]
		num = num + onerecord["DownNumber"].(int64)
	}
	return num, err
}

//根据分发给的人取得总的下载量
func GetAllDownNum(id int64) int64 {
	o := orm.NewOrm()
	var downrecord []orm.Params
	var count int64 = 0
	_, err := o.QueryTable("cloud_distribut_record").Filter("DistributName", id).Values(&downrecord)
	if err != nil {
		beego.Error("get data is err", err)
		return 0
	}
	for i := 0; i < len(downrecord); i++ {
		num, _ := downrecord[i]["DownNumber"].(int64)
		count = count + num
	}
	return count
}
