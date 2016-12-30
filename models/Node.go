package models

import (
	"errors"
	"log"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
)

type Node struct {
	Id          int64
	Name        string  `orm:"size(32)" form:"Name" valid:"Required" `
	Pid         int64   `orm:size(11) form:"Pid" valid:"Required"`
	Key         string  `orm:"size(64)" form :"Key" valid:"Required"`
	Type        string  `orm:"size(10)" form:"Type" `
	Ico         string  `orm: form:"Ico"`
	Url         string  `orm:"size(64)" form:"Url"`
	Fid         int64   // 目录 菜单  按钮
	Level       int64   `orm:"default(1);size(11)" form:"Level"`
	Description string  `orm:"null;size(200)" form:"Description" valid:"MaxSize(200)"`
	Sort        int64   `orm:"default(100);size(11)"`
	Role        []*Role `orm:"rel(m2m)"`

	Typename string `orm:"-" form:"Typename" valid:"Required" `
	Pname    string `orm:"-" form:"Pname" valid:"Required" `
}

func (r *Node) TableName() string {
	return "node"
}

func checkNode(r *Node) (err error) {
	valid := validation.Validation{}
	b, _ := valid.Valid(&r)
	if !b {
		for _, err := range valid.Errors {
			log.Println(err.Key, err.Message)
			return errors.New(err.Message)
		}
	}
	return nil
}

func init() {
	orm.RegisterModel(new(Node))
}

func GetNodes(page int64, page_size int64, sort string) (nodes []orm.Params, count int64) {
	o := orm.NewOrm()
	node := new(Node)
	qs := o.QueryTable(node)
	qs.Limit(page_size, page).OrderBy(sort).Values(&nodes)
	count, err := qs.Count()
	if err != nil {
		beego.Error(err)
		// fmt.Println(err)
	}
	for _, node := range nodes {
		res, _ := GetNode(node["Pid"].(int64))
		node["Pname"] = res.Name
	}
	return nodes, count
}

func GetAllNodes() (nodes []*Node, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("node").All(&nodes)
	return nodes, err
}

func GetAllNode() (count int64, nodes []orm.Params) {
	o := orm.NewOrm()
	node := Node{}
	qs := o.QueryTable(node)
	count, err := qs.OrderBy("level", "sort").Values(&nodes)
	if err != nil {
		// fmt.Println(err)
		beego.Error(err)
	}
	return count, nodes
}

func GetParentNode(pid int64) (count int64, nodes []orm.Params) {
	o := orm.NewOrm()
	node := new(Node)
	qs := o.QueryTable(node).Filter("Pid", pid)
	count, err := qs.Values(&nodes)
	if err != nil {
		// fmt.Println(err)
		beego.Error(err)
	}
	return count, nodes
}

func GetNodeByName(name string) (err error, node Node) {
	o := orm.NewOrm()
	node = Node{Name: name}
	err = o.Read(&node, "Name")
	return err, node
}

func GetNode(id int64) (Node, error) {
	o := orm.NewOrm()
	node := Node{Id: id}
	err := o.Read(&node)
	if err != nil {
		beego.Error(err)
		return node, nil
	}
	return node, nil
}

func AddNode(r *Node) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(r)
	return id, err
}

func UpdateNode(r *Node) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Update(r)
	return id, err
}
func DelNodeById(Id int64) (int64, error) {
	o := orm.NewOrm()
	status, err := o.Delete(&Node{Id: Id})
	return status, err
}

func DelRoleRescoursByRoleId(roleId int64) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("node_roles").Filter("role_id", roleId).Delete()
	return err

}
func AddRoleRescours(roleId int64, resList []string) error {
	o := orm.NewOrm()
	role := Role{Id: roleId}
	for _, v := range resList {
		nodeId, _ := strconv.Atoi(v)
		node := Node{Id: int64(nodeId)}
		m2m := o.QueryM2M(&node, "Role")
		m2m.Add(&role)
	}

	return nil
}

func GetNodesByGroupId(GroupId int64) (nodes []orm.Params, count int64) {
	o := orm.NewOrm()
	node := new(Node)
	count, _ = o.QueryTable(node).Filter("Group", GroupId).Values(&nodes)
	return nodes, count
}

func GetNodeByRoleId(id int64) (nodes []orm.Params, count int64) {
	o := orm.NewOrm()
	node := new(Node)
	count, _ = o.QueryTable(node).Filter("Role__Role__Id", id).Values(&nodes)
	return nodes, count
}

// 无限级分类
func GetNodeTreeAndLv(node []*Node, fid int64) []*Node {
	// 创建一个数组
	var tree []*Node = make([]*Node, 0)
	// 循环
	for _, v := range node {
		// 传来的集合的父id 等于 传来的父id
		if v.Fid == fid {
			// 层级
			child := GetNodeTreeAndLv(node, v.Id)
			tree = append(tree, v)
			tree = append(tree, child...)
		}
	}
	return tree
}

func GetNodeTree(pid int64, resType int64) ([]orm.Params, error) {
	o := orm.NewOrm()
	node := new(Node)
	var nodes []orm.Params
	_, err := o.QueryTable(node).Filter("Pid", pid).Filter("Type", resType).OrderBy("Level").Values(&nodes)
	if err != nil {
		beego.Error(err)
		return nodes, err
	}
	return nodes, nil
}
func GetResSubTree() {

}

//根据url获取节点信息
func GetResByUrl(url string) (res Node) {
	o := orm.NewOrm()
	o.QueryTable("node").Filter("URL", url).One(&res)
	return res
}

//判断是否有权限
func CheckAccessRole(resId int64, roleId int64) bool {
	o := orm.NewOrm()
	var temp []orm.Params
	o.QueryTable("node_roles").Filter("node_id", resId).Filter("role_id", roleId).Values(&temp)
	if len(temp) <= 0 {
		return false
	} else {
		return true
	}
}
