package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

//资源类型
type CloudResourceType struct {
	Id                int64
	TypeName          string               `orm:"unique"` //资源类型名称
	AddUser           int64                //资源类型添加者
	AddTime           time.Time            `orm:"type(datetime);auto_now_add"` //资源类型添加日期
	Status            int64                `orm:"default(1)"`                  // 状态 [1： 私有  2： 公开]
	CloudResources    []*CloudResources    `orm:"reverse(many)"`               //一对多资源
	CloudUploadRecord []*CloudUploadRecord `orm:"reverse(many)"`
	AddUserName       string               `orm:"-"` // 用户名称
	Users             []*User              `orm:"rel(m2m)"`
}

func (r *CloudResourceType) TableName() string {
	return "cloud_resource_type"
}
func init() {
	orm.RegisterModel(new(CloudResourceType))
}

// 分页查询资源分类
func GetTypeByCond(page int64, page_size int64, username string, userid int64, status int64, typename string) (list []orm.Params, count int64, err error) {
	o := orm.NewOrm()
	cond := orm.NewCondition()
	qs := o.QueryTable("cloud_resource_type")
	// 状态
	if status > 0 {
		cond = cond.And("Status", status)
	}
	// 类型名称
	if len(typename) > 0 {
		cond = cond.And("TypeName__icontains", typename)
	}
	// 用户
	if username != beego.AppConfig.String("admin_user") {
		cond = cond.And("AddUser", userid)
	}
	qs.Limit(page_size, page).OrderBy("Id").SetCond(cond).Values(&list)
	count, err = qs.SetCond(cond).Count()
	return list, count, err
}

// 通过Id获取资源分类
func GetTypeById(id int64) (err error, rt CloudResourceType) {
	o := orm.NewOrm()
	rt = CloudResourceType{Id: id}
	err = o.Read(&rt)
	return err, rt
}

//取得所有的分类
func GetAllType() (list []orm.Params) {
	o := orm.NewOrm()
	_, err := o.QueryTable("cloud_resource_type").Values(&list)
	if err != nil {
		beego.Error("get data is err", err)
	}
	return list
}
func GetAllTypeNum() int64 {
	o := orm.NewOrm()
	num, _ := o.QueryTable("cloud_resource_type").Count()
	return num
}

//取得所有的资源类型
func GetAllResourceType() (list []orm.Params, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("cloud_resource_type").Values(&list)
	return list, err
}

//增加资源类型
func AddResType(r *CloudResourceType) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(r)
	return id, err
}

//修改资源类型
func EditResType(r *CloudResourceType) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Update(r)
	return id, err
}

//核实增加的新的资源类型是否重复
func CheckIsRepat(name string) (flag bool) {
	o := orm.NewOrm()
	flag = o.QueryTable("cloud_resource_type").Filter("TypeName", name).Exist()
	return flag
}

//根据id取得资源名称
func GetTypeNameByID(id int64) CloudResourceType {
	o := orm.NewOrm()
	var types CloudResourceType
	err := o.QueryTable("cloud_resource_type").Filter("Id", id).One(&types, "TypeName")
	if err != nil {
		beego.Error("get data is err", err)
	}
	return types
}

// 添加
func AddTypeUser(uids []int64, tids []int64) (num int64, err error) {
	o := orm.NewOrm()
	err = o.Begin()
	for i := 0; i < len(tids); i++ {
		for j := 0; j < len(uids); j++ {
			rt := CloudResourceType{Id: tids[i]}
			user := User{Id: uids[j]}
			m2m := o.QueryM2M(&rt, "Users")
			num, err = m2m.Add(&user)
			num += num
			beego.Debug(num)
		}
	}
	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	return num, err
}

// 删除资源类型
func DelType(ids []int64) (err error) {
	o := orm.NewOrm()
	err = o.Begin()
	for i := 0; i < len(ids); i++ {
		_, err = o.QueryTable("cloud_resource_type").Filter("Id", ids[i]).Delete()
	}
	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	return err
}

// 删除
func DelTypeUser(tids []int64) (err error) {
	o := orm.NewOrm()
	err = o.Begin()
	for i := 0; i < len(tids); i++ {
		_, err = o.QueryTable("cloud_resource_type_users").Filter("cloud_resource_type_id", tids[i]).Delete()
	}
	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	return err
}

//查询这个人有没有这个资源类型的权限
func GetIsPermission(restype int64, user int64) bool {
	o := orm.NewOrm()
	flag := o.QueryTable("cloud_resource_type_users").Filter("user_id", user).Filter("cloud_resource_type_id", restype).Exist()
	return flag
}

//根据用户id取得有权限的资源类型
func GetPermissionType(id int64) (list []orm.Params, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("cloud_resource_type_users").Filter("user_id", id).Values(&list)
	return list, err
}

//取得有权限信息
func GetUserPermissionRes(res []int64) (list []orm.Params, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("cloud_resource_type").Filter("Id__in", res).Values(&list)
	return list, err
}
