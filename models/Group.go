package models

import (
	"fmt"
	// "github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Group struct {
	Id     int64
	Name   string //组名称
	Uid    int64  //组的管理者
	Status int64  `orm:"default(0)"` //状态(1:开启，2关闭)
	Fid    int64  `orm:default(0)`   //默认为0
	Sort   int64  `orm:default(50)`  //排序(默认为50)
	Remark string //组别备注
	Tid    int64  //顶级团队Id

	OneDayLimit int64 `orm:default(0)` //今日资源上限
	OnceLimit   int64 //团队个人单次拉取值

	Level       int64  `orm:"-"`
	ManagerName string `orm:"-"`
}

func (this *Group) TableName() string {
	return "group"
}

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(Group))
}
func GetGroupName(name string) (list Group, err error) {
	o := orm.NewOrm()
	err = o.QueryTable("group").Filter("Name", name).One(&list)
	return list, err
}

//带条件的组别查询
func GetAllGroups(page int64, page_size int64, sort string, filter map[string]interface{}, user []int64) (groups []orm.Params, count int64) {
	o := orm.NewOrm()
	group := new(Group)
	qs := o.QueryTable(group)

	var cond *orm.Condition
	cond = orm.NewCondition()
	s1, _ := filter["s0"].(int64) //状态
	// s2, _ := filter["s1"].(int64) //管理者
	s3 := filter["s2"].(string) //团队名称

	if s1 != 0 {
		cond = cond.And("Status", s1)
	}
	if len(user) > 0 {
		cond = cond.And("Uid__in", user)
	}
	if len(s3) > 0 {
		cond = cond.And("Name__icontains", s3)
	}

	qs.SetCond(cond).OrderBy(sort).Limit(page_size, page).Values(&groups)
	count, _ = qs.SetCond(cond).Count()
	return groups, count
}

func GetGroups(page int64, page_size int64, sort string) (groups []orm.Params, count int64) {
	groupQuerySeter().Limit(page_size, page).OrderBy(sort).Values(&groups)
	count, _ = groupQuerySeter().Count()
	return groups, count
}

func AddGroup(g *Group) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(g)
	return id, err
}

func UpdateGroup(g *Group) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Update(g)
	return num, err
}

//删除组别
func DelGroupById(Id int64) (int64, error) {
	o := orm.NewOrm()
	num, err := o.QueryTable("group").Filter("Id", Id).Delete()

	//将该组下的所有组别归于公司下
	group := make(orm.Params)
	group["Fid"] = 0
	o.QueryTable("group").Filter("Fid", Id).Update(group)

	return num, err
}

/*
*获取所有的分组实体
 */
func Groups() (groups []*Group, count int64) {
	o := orm.NewOrm()
	count, _ = o.QueryTable("group").All(&groups)
	return groups, count
}

func GroupsByCond(status int64, groupname string) (groups []*Group, count int64) {
	o := orm.NewOrm()
	qs := o.QueryTable("group")
	cond := orm.NewCondition()
	if status > 0 {
		cond = cond.And("Status", status)
	}
	if len(groupname) > 0 {
		cond = cond.And("Name__icontains", groupname)
	}
	qs.SetCond(cond).All(&groups)
	count, _ = qs.SetCond(cond).Count()
	return groups, count
}

// 获取所有分组
func GroupCate() (groups []*Group) {
	groupQuerySeter().All(&groups)
	return
}

//按照id获取
func GroupSingleCate(id int64) (groups []*Group) {
	groupQuerySeter().Filter("Id", id).All(&groups)
	return
}

/*
*根据数组获取所有的分组实体
 */
func GetGroupsByArr(arr []int64) (groups []orm.Params, count int64) {
	o := orm.NewOrm()
	count, _ = o.QueryTable("group").Filter("Id__in", arr).Values(&groups)
	return groups, count
}

// 无限极分类
func GetTree(cate []*Group, fid int64) []*Group {
	var tree []*Group = make([]*Group, 0)

	for _, v := range cate {
		if v.Fid == fid {
			child := GetTree(cate, v.Id)
			tree = append(tree, v)
			tree = append(tree, child...)
		}
	}
	return tree
}

// 无限极分类
func GetTreeAndLv(cate []*Group, fid int64, lv int64) []*Group {
	var tree []*Group = make([]*Group, 0)

	for _, v := range cate {
		if v.Fid == fid {
			v.Level = lv
			child := GetTreeAndLv(cate, v.Id, lv+1)
			tree = append(tree, v)
			tree = append(tree, child...)
		}
	}
	return tree
}

// 根据给出的组循环出ID
func GetGroupIdlimit(cate []*Group) []int64 {
	var groupcate = make([]int64, 0)
	for _, v := range cate {
		groupcate = append(groupcate, v.Id)
	}
	return groupcate
}

// 只循环出所有的ID
func GetTreeGroup(cate []*Group, fid int64) []int64 {
	var tree = make([]int64, 0)
	for _, v := range cate {
		if v.Fid == fid {
			childid := GetTreeGroup(cate, v.Id)
			tree = append(tree, v.Id)
			tree = append(tree, childid...)
		}
	}
	return tree
}

// 判断Group是否存在
func CheckGroupIsExist(name string) bool {
	return groupQuerySeter().Filter("Name", name).Exist()
}

// 通过ID进行查询
func GetGroupById(id int64) (Group, error) {
	model := orm.NewOrm()
	group := Group{Id: id}

	err := model.Read(&group)
	return group, err
}

//获取所有下级的gooupid字符串
func GetTreeGroupStr(id int64) (idStr string) {
	group, _ := GetGroupById(id)
	if group.Fid != 0 {
		return fmt.Sprintf("%d", id)
	} else {
		model := orm.NewOrm()
		var table Group
		var list []orm.Params
		model.QueryTable(table).Filter("Fid", id).Values(&list)
		idStr = fmt.Sprintf("%d", id) + ","
		for _, val := range list {
			idStr = idStr + fmt.Sprintf("%d", val["Id"].(int64)) + ","
		}
		idStr = idStr[0 : len(idStr)-1]
		return idStr
	}
}

//通过fid获取所有组别信息
func GetGroupByFid() (groups []orm.Params) {
	o := orm.NewOrm()
	o.QueryTable("group").Values(&groups, "Id", "name")
	return groups
}

//根据组别查询用户
func GetUserByGroup(id int64) (list []orm.Params) {
	o := orm.NewOrm()
	o.QueryTable("user").Filter("Group", id).Filter("Status", 2).Values(&list)
	return list
}

// 判断用户是否是管理
func CheckUserIsAdmin(uid int64) bool {
	return groupQuerySeter().Filter("Fid", 0).Filter("Uid", uid).Exist()
}

//
func groupQuerySeter() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(Group))
}

//根据id修改团队资源信息
func ModifyGroupResource(id, onedaytotal, oncedown int64) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("group").Filter("Id", id).Update(orm.Params{
		"OneDayLimit": onedaytotal,
		"OnceLimit":   oncedown,
	})
	return err
}

//得到处于关闭状态的团队数量
func GetShutDownTeam() int64 {
	o := orm.NewOrm()
	num, _ := o.QueryTable("group").Filter("Status", 2).Count()
	return num
}

//打的处于开启状态的团队数量
func GetOpenTeam() int64 {
	o := orm.NewOrm()
	num, _ := o.QueryTable("group").Filter("Status", 1).Count()
	return num
}

//取得所有的团队信息
func GetAllTeamInfo() (list []orm.Params, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("group").Values(&list)
	return list, err
}

//取得
func GetAllTeamCount() int64 {
	o := orm.NewOrm()
	num, _ := o.QueryTable("group").Count()
	return num
}
