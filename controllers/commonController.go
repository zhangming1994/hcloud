package controllers

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	m "hcloud/models"
)

type CommonController struct {
	beego.Controller
}

func init() {
	// m.CheckAccessAndRegisterRes()
}

type PageType int64

const (
	USERPAGE PageType = iota + 1
)

var (
	pageTypeMap = map[int64]PageType{
		1: USERPAGE, //用户分页
	}
)

func (this *CommonController) Rsp(status bool, str string) {
	this.Data["json"] = &map[string]interface{}{"status": status, "info": str}
	this.ServeJSON()
}

func (this *CommonController) RspImg(status bool, str string, img string) {
	this.Data["json"] = &map[string]interface{}{"status": status, "info": str, "img": img}
	this.ServeJSON()
}

func (this *CommonController) Alert(msg string, localurl string) {
	alert := fmt.Sprintf("<script>alert('" + msg + "');location.href='" + localurl + "';</script>")
	this.Ctx.WriteString(alert)
}

// Getstring
func (this *CommonController) GetTrimString(key string) string {
	if v := this.Ctx.Input.Query(key); v != "" {
		return strings.TrimSpace(v)
	}
	return ""
}

// 是否POST提交
func (this *CommonController) IsPost() bool {
	return this.Ctx.Request.Method == "POST"
}

func (this *CommonController) GetResList(username string, Id int64) []Tree {
	var cnt, length int = 0, 0
	var nodes []orm.Params
	adminUser := beego.AppConfig.String("admin_user")
	if username == adminUser {
		_, nodes = m.GetAllNode()
	} else {
		nodes, _ = m.GetNodesByRoleId(Id)
	}

	for _, v := range nodes {
		if v["Pid"].(int64) == 0 {
			length = length + 1
		}
	}
	tree := make([]Tree, length)

	for k, v := range nodes {
		if v["Pid"].(int64) == 0 {
			k = cnt
			cnt = cnt + 1
			tree[k].Id = v["Id"].(int64)
			tree[k].Index = cnt
			tree[k].Url = "/" + v["Url"].(string)
			tree[k].Text = v["Name"].(string)
			tree[k].IconCls = v["Ico"].(string)
			// 1代表菜单（目录下面的所有资源）没有把一些不需要的权限去掉
			var childCnt int = 0
			children := make([]map[string]interface{}, 6)
			for _, v3 := range nodes {
				if v3["Pid"].(int64) == v["Id"].(int64) {
					children[childCnt] = v3
					childCnt++
				}
			}
			tree[k].Children = make([]Tree, childCnt)
			for k1, v1 := range children {
				if v1 != nil {
					if v1["Pid"].(int64) == v["Id"].(int64) {
						tree[k].Children[k1].Id = v1["Id"].(int64)
						tree[k].Children[k1].Text = v1["Name"].(string)
						tree[k].Children[k1].Url = "/" + v1["Url"].(string)
						tree[k].Children[k1].IconCls = v1["Ico"].(string)
					}
				}
			}

		}

	}
	return tree
}

func (this *CommonController) GetTree() []Tree {
	nodes, _ := m.GetNodeTree(0, 1)
	tree := make([]Tree, len(nodes))
	for k, v := range nodes {
		tree[k].Id = v["Id"].(int64)
		tree[k].Text = v["Title"].(string)
		tree[k].IconCls = v["Ico"].(string)
		children, _ := m.GetNodeTree(v["Id"].(int64), 2)
		tree[k].Children = make([]Tree, len(children))
		for k1, v1 := range children {
			tree[k].Children[k1].Id = v1["Id"].(int64)
			tree[k].Children[k1].Text = v1["Title"].(string)
			tree[k].Children[k1].Url = "/" + v["Name"].(string) + "/" + v1["Name"].(string)
			tree[k].Children[k1].IconCls = v1["Ico"].(string)
		}
	}
	return tree
}

// 获取功能列表(显示行隐藏菜单)
func (this *CommonController) GetFuncList(username string, Id int64) []string {
	var length int = 0
	var nodes []orm.Params
	adminUser := beego.AppConfig.String("admin_user")
	if username == adminUser {
		_, nodes = m.GetAllNode()
	} else {
		nodes, _ = m.GetNodesByRoleId(Id)
	}

	for _, v := range nodes {
		if v["Pid"].(int64) == 0 {
			length = length + 1
		}
	}
	tree := make([]string, length)
	for _, v := range nodes {
		if v["Pid"].(int64) == 1000 {
			tree = append(tree, v["Key"].(string))
		}
	}
	return tree
}

func (this *CommonController) CommonMenu() {
	userInfo := this.GetSession("userinfo")
	if userInfo == nil {
		this.Ctx.Redirect(302, beego.AppConfig.String("auth_gateway"))
	} else {
		role := m.GetRoleByUserId(userInfo.(m.User).Id)
		tree := this.GetResList(userInfo.(m.User).Username, role.Id)
		functree := this.GetFuncList(userInfo.(m.User).Username, role.Id)

		this.Data["functool"] = &functree
		this.Data["tree"] = &tree
		this.Data["userinfo"] = userInfo.(m.User)
	}
	this.Layout = "layout/layout.html"
}

//获取客户的真是IP地址
func (this *CommonController) GetClientip() string {
	var addr string
	if len(this.Ctx.Request.Header.Get("X-Forwarded-For")) > 0 {
		addr := this.Ctx.Request.Header.Get("X-Forwarded-For")
		return addr
	} else if len(this.Ctx.Request.RemoteAddr) > 0 {
		addr := this.Ctx.Request.RemoteAddr
		return addr
	} else {
		addr := "127.0.0.1:8080"
		return addr
	}
	beego.Debug(addr)
	return ""
}

// slice cap
// arge :
//	0:	start
// 	1:	length
// 	2:	分页类型(1)
//  3:	fid (查找user)
// arge := []{0,10,1,0}
func (this *CommonController) CutSlice(data interface{}, arge ...int64) (dataStruce interface{}, count int) {
	var (
		start  = int(arge[0])
		length = int(arge[1])
	)

	switch pageTypeMap[arge[2]] {
	case USERPAGE:
		dataUser := m.GetUserTree(data.([]*m.User), arge[3])
		count = len(dataUser)
		startLine, endLine := this.CalculatePage(count, start, length)
		dataStruce = dataUser[startLine:endLine]
		break
	default:
	}
	return
}

// 计算页数
func (this *CommonController) CalculatePage(count, start, length int) (objstart, objend int) {
	if length == -1 {
		length = count
	}
	ipagetotal := count / length
	if 0 != count%length {
		ipagetotal = ipagetotal + 1
	}
	if ipagetotal == 1 {
		objstart = 0
		objend = count
	} else {
		objstart = start
		objend = start + length
		if objend > count {
			objend = count
		}
	}
	return
}

func (this *CommonController) CutGroupSlice(data interface{}, arge ...int64) (dataStruce interface{}, count int) {
	var (
		start  = int(arge[0])
		length = int(arge[1])
	)
	switch pageTypeMap[arge[2]] {
	case USERPAGE:
		dataGroup := m.GetTreeAndLv(data.([]*m.Group), arge[3], 0)
		count = len(dataGroup)
		startLine, endLine := this.CalculatePage(count, start, length)
		dataStruce = dataGroup[startLine:endLine]
		break
	default:
	}
	return
}
