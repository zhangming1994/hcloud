package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

//上传记录表
type CloudUploadRecord struct {
	Id                int64
	CloudResourceType *CloudResourceType `orm:"rel(fk)"` //类型
	RepatNumber       int64              //重复数量
	SuccessNumber     int64              //成功数量
	FailuedNumber     int64              //失败数量
	TotalNumber       int64              //上传总数
	UploadStatus      int                `orm:"default(0)"` //上传状态 2:成功 1：失败 0:上传资源处理中,
	UploadName        string             //上传文件名
	UploadUser        int64              //上传人
	Team              int64              //所属团队
	Persentage        string             //进度
	UploadDate        time.Time          `orm:"type(datetime);auto_now_add"` //上传日期
}

func (m *CloudUploadRecord) TableName() string {
	return "cloud_upload_record"
}
func init() {
	orm.RegisterModel(new(CloudUploadRecord))
}

// 上传记录列表
func ListUpload(page int64, pageSize int64, filter map[string]interface{}, uploaduser []int64, resourcetypeid []int64) (list []orm.Params, count int64, err error) {
	o := orm.NewOrm()
	mUpload := new(CloudUploadRecord)
	qs := o.QueryTable(mUpload)
	starttimes := filter["starttime"].(string) //开始时间
	endtimes := filter["edntime"].(string)     //结束时间
	restypes, _ := filter["s1"].(string)       //资源类型
	teams, _ := filter["s2"].(string)          //所属团队

	restype, _ := strconv.ParseInt(restypes, 10, 64)
	team, _ := strconv.ParseInt(teams, 10, 64)
	starttime, _ := time.Parse("2006-01-02", starttimes)
	endtime, _ := time.Parse("2006-01-02", endtimes)

	var cond *orm.Condition
	cond = orm.NewCondition()

	if len(resourcetypeid) > 0 {
		cond = cond.And("CloudResourceType__Id__in", resourcetypeid)
	}
	if len(uploaduser) > 0 {
		cond = cond.And("UploadUser__in", uploaduser)
	}
	if len(starttimes) > 0 {
		cond = cond.And("UploadDate__gte", starttime)
		cond = cond.And("UploadDate__lte", endtime)
	}
	if restype != 0 {
		cond = cond.And("CloudResourceType__Id", restype)
	}
	if team != 0 {
		cond = cond.And("Team", team)
	}
	qs.SetCond(cond).OrderBy("-Id").Limit(pageSize, page).Values(&list)
	count, _ = qs.SetCond(cond).Count()
	return list, count, err
}

//管理员记录列表
func ListUploads(page int64, pageSize int64, filter map[string]interface{}) (list []orm.Params, count int64, err error) {
	o := orm.NewOrm()
	mUpload := new(CloudUploadRecord)
	qs := o.QueryTable(mUpload)

	starttimes := filter["starttime"].(string) //开始时间
	endtimes := filter["edntime"].(string)     //结束时间
	restypes, _ := filter["s1"].(string)       //资源类型
	teams, _ := filter["s2"].(string)          //所属团队

	restype, _ := strconv.ParseInt(restypes, 10, 64)
	team, _ := strconv.ParseInt(teams, 10, 64)
	starttime, _ := time.Parse("2006-01-02 15:04:05", starttimes)
	endtime, _ := time.Parse("2006-01-02 15:04:05", endtimes)

	var cond *orm.Condition
	cond = orm.NewCondition()

	if len(starttimes) > 0 {
		cond = cond.And("UploadDate__gte", starttime)
		cond = cond.And("UploadDate__lte", endtime)
	}
	if restype != 0 {
		cond = cond.And("CloudResourceType__Id", restype)
	}
	if team != 0 {
		cond = cond.And("Team", team)
	}
	qs.SetCond(cond).OrderBy("-Id").Limit(pageSize, page).Values(&list)
	count, _ = qs.SetCond(cond).Count()
	return list, count, err
}

//增加
func AddUploadRecord(r *CloudUploadRecord) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Insert(r)
	return num, err
}

//根据id更改值
func ModifyById(id, repat, success, fail, total int64) (err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("cloud_upload_record").Filter("Id", id).Update(orm.Params{
		"RepatNumber":   repat,
		"SuccessNumber": success,
		"FailuedNumber": fail,
		"TotalNumber":   total,
	})
	return err
}

//根据上传批次id修改状态
func ModifyStatus(status int, id int64) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("cloud_upload_record").Filter("Id", id).Update(orm.Params{
		"UploadStatus": status,
	})
	return err
}

//根据批次id更新Bar
func UpdateBar(id, total, success, failenumber, repatnumber int64, num float64) error {
	o := orm.NewOrm()
	persentage := fmt.Sprintf("%.2f", num)
	_, err := o.QueryTable("cloud_upload_record").Filter("Id", id).Update(orm.Params{
		"Persentage":    persentage,
		"RepatNumber":   repatnumber,
		"SuccessNumber": success,
		"FailuedNumber": failenumber,
		"TotalNumber":   total,
	})
	return err
}

//根据上传人取得个人上传总数
func GetPersonTotal(userid int64) int64 {
	o := orm.NewOrm()
	var count int64 = 0
	var uploadrecord []orm.Params
	_, err := o.QueryTable("cloud_upload_record").Filter("UploadUser", userid).Values(&uploadrecord)
	if err != nil {
		beego.Error("get data is err", err)
		return 0
	}
	for i := 0; i < len(uploadrecord); i++ {
		num, _ := uploadrecord[i]["SuccessNumber"].(int64)
		count = count + num
	}
	return count
}
