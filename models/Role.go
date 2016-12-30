package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Role struct {
	Id          int64
	Name        string  `orm:"size(64);unique"`
	Status      int64   `orm:"default(1)"` //Status  	1 正常    2 禁用
	Description string  `orm:"type(text)" `
	Statusname  string  `orm:"-" `
	Isnormal    bool    `orm:"-"`
	User        []*User `orm:"reverse(many)"`
	Node        []*Node `orm:"reverse(many)"`
}

func (r *Role) TableName() string {
	return "role"
}

/*
*注册model
 */
func init() {
	orm.RegisterModel(new(Role))
}

func GetRoleList(page int64, page_size int64, status int64, rolename string) (roles []orm.Params, count int64) {
	o := orm.NewOrm()
	cond := orm.NewCondition()
	qs := o.QueryTable("role")
	if status > 0 {
		cond = cond.And("Status", status)
	}
	if len(rolename) > 0 {
		cond = cond.And("Name__icontains", rolename)
	}
	qs.Limit(page_size, page).OrderBy("Id").SetCond(cond).Values(&roles)
	count, _ = qs.SetCond(cond).Count()
	return roles, count
}

// 添加角色
func AddRole(r *Role) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(r)
	return id, err
}

// 更新角色
func UpdateRole(r *Role) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Update(r)
	return id, err
}

// 删除角色
func DelRoleById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Role{Id: Id})
	return status, err
}

// 通过角色ID获取用户的节点
func GetNodesByRoleId(Id int64) (nodes []orm.Params, count int64) {
	o := orm.NewOrm()
	node := new(Node)
	count, err := o.QueryTable(node).Filter("Role__Role__Id", Id).Values(&nodes)
	if err != nil {
		beego.Error(err)
	}
	return nodes, count
}

// 获取所有角色列表
func GetAllRole() (roles []orm.Params, count int64) {
	o := orm.NewOrm()
	role := new(Role)
	qs := o.QueryTable(role)
	count, err := qs.Values(&roles)
	if err != nil {
		beego.Error("get the role error:", err)
	}
	return roles, count
}

// 增加角色节点
func AddRoleNode(roleId int64, nodeId int64) (int64, error) {
	o := orm.NewOrm()
	role := Role{Id: roleId}
	node := Node{Id: nodeId}
	m2m := o.QueryM2M(&node, "Role")
	num, err := m2m.Add(&role)
	return num, err
}

// 通过角色名获取角色
func GetRoleByName(Name string) (err error, role Role) {
	o := orm.NewOrm()
	role = Role{Name: Name}
	err = o.Read(&role, "Name")
	return err, role
}

// 通过Id获取角色
func GetRoleById(Id int64) (err error, role Role) {
	o := orm.NewOrm()
	role = Role{Id: Id}
	err = o.Read(&role)
	return err, role
}

// 通过角色Id删除角色
func DelUserRole(roleId int64) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("user_roles").Filter("role_id", roleId).Delete()
	return err
}

// 通过用户id删除角色
func DelUserRoleByUserId(userId int64) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("user_roles").Filter("user_id", userId).Delete()
	return err
}

// 增加角色
func AddRoleUser(roleid int64, userid int64) (int64, error) {
	o := orm.NewOrm()
	role := Role{Id: roleid}
	user := User{Id: userid}
	m2m := o.QueryM2M(&user, "Role")
	num, err := m2m.Add(&role)
	return num, err
}

// 根据用户id获取角色列表
func AccessList(uid int64) (list []orm.Params, err error) {
	o := orm.NewOrm()
	user := GetUserById(uid)
	var nodes []orm.Params
	_, err = o.QueryTable("node").Filter("Role__Id", user.Role.Id).Values(&nodes)
	if err != nil {
		beego.Error("get the node error:", err)
		return nil, err
	}
	for _, n := range nodes {
		list = append(list, n)
	}
	return list, nil
}
