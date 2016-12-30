package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

//资源列表页面首页统计
type CloudTotalResource struct {
	Id             int64
	ResourceDate   time.Time `orm:"type(date)"` //统计日期
	Group          int64     //团队
	ManageUser     int64     //管理者
	TotalResource  int64     //资源总数
	UsdeResource   int64     //已经使用的资源
	CanUseResource int64     //未使用的资源
	UsePersent     string    //资源使用率
}

func (r *CloudTotalResource) TableName() string {
	return "cloud_total_resource"
}
func init() {
	orm.RegisterModel(new(CloudTotalResource))
}

//取得该日期下面所有的团队数据
func CheckIsDates(tiemdate time.Time) int64 {
	o := orm.NewOrm()
	num, _ := o.QueryTable("cloud_total_resource").Filter("ResourceDate", tiemdate).Count()
	return num
}

//查询有没有这个日期的数据
func CheckIsDate(tiemdate time.Time, groupid int64) bool {
	o := orm.NewOrm()
	flag := o.QueryTable("cloud_total_resource").Filter("ResourceDate", tiemdate).Filter("Group", groupid).Exist()
	return flag
}

//列表
func CloudResourceIndexList(page int64, pageSize int64, filter map[string]interface{}, userids []int64) (list []orm.Params, count int64, err error) {
	o := orm.NewOrm()
	mResources := new(CloudTotalResource)
	qs := o.QueryTable(mResources)

	var cond *orm.Condition
	cond = orm.NewCondition()

	starttimes := filter["starttime"].(string) //开始时间
	endtimes := filter["edntime"].(string)     //结束时间
	teams, _ := filter["s1"].(string)          //团队

	starttime, _ := time.Parse("2006-01-02", starttimes)
	endtime, _ := time.Parse("2006-01-02", endtimes)
	team, _ := strconv.ParseInt(teams, 10, 64)
	if team != 0 {
		cond = cond.And("Group", team)
	}
	if len(userids) > 0 {
		cond = cond.And("ManageUser__in", userids)
	}
	if len(starttimes) > 0 {
		cond = cond.And("ResourceDate__gte", starttime)
		cond = cond.And("ResourceDate__lte", endtime)
	}
	qs.SetCond(cond).OrderBy("-Id").Limit(pageSize, page).Values(&list)
	count, _ = qs.SetCond(cond).Count()
	return list, count, err
}

//增加
func AddCloudResouIndex(r *CloudTotalResource) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Insert(r)
	return num, err
}

//下发资源更改数据
func ModifyByType(num int64, downtime time.Time, teamid int64) error {
	o := orm.NewOrm()
	list, _ := GetTotalResource(downtime, teamid)
	useresource := list.UsdeResource + num
	notuseresource := list.CanUseResource - num
	respersent := fmt.Sprintf("%.2f", float64(useresource)/float64(list.TotalResource)*100)
	_, err := o.QueryTable("cloud_total_resource").Filter("ResourceDate", downtime).Filter("Group", teamid).Update(orm.Params{
		"UsdeResource":   useresource,
		"CanUseResource": notuseresource,
		"UsePersent":     respersent,
	})
	return err
}

//取得该资源
func GetTotalResource(downtime time.Time, teamid int64) (CloudTotalResource, error) {
	o := orm.NewOrm()
	var list CloudTotalResource
	err := o.QueryTable("cloud_total_resource").Filter("ResourceDate", downtime.Format("2006-01-02")).Filter("Group", teamid).One(&list)
	if err == orm.ErrMultiRows {
		beego.Error("Returned Multi Rows Not One")
	}
	if err == orm.ErrNoRows {
		beego.Error("Not row found")
	}
	return list, err
}
