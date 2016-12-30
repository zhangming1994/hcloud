package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	m "hcloud/models"
	// "github.com/liudng/godump"
)

type MainController struct {
	CommonController
}
type Tree struct {
	Id       int64
	Index    int    `json:"index"`
	Text     string `json:"id"`
	IconCls  string `json:"text"`
	Checked  string `json:"iconCls"`
	State    string `json:"checked"`
	Children []Tree `json:"state"`
	Url      string `json:"url"`
	Fid      int64
}

//主页
func (this *MainController) Index() {
	userinfo := this.GetSession("userinfo")
	if userinfo == nil {
		this.Ctx.Redirect(302, beego.AppConfig.String("auth_gateway"))
	}

	tree := this.GetTree()
	if this.IsAjax() {
		this.Data["json"] = &tree
		this.ServeJSON()
	} else {
		this.Data["userinfo"] = userinfo
		this.Data["groups"] = m.GroupCate()
		this.Data["tree"] = &tree

		this.CommonMenu()
		this.TplName = "index.html"
	}
}

func (this *MainController) Login() {
	if this.IsPost() {
		username := this.GetString("username") //用户名或工号
		password := this.GetString("password")
		user, err := m.Login(username, password)
		if err == nil {
			if user.Status == 2 {
				err := "用户被禁用！"
				code := -1
				data := fmt.Sprintf(`{"code":%d,"err":"%s"}`, code, err)
				this.Ctx.WriteString(data)
				return
			} else {
				this.SetSession("userinfo", user)
				accessList, _ := m.GetAccessList(user.Id)
				this.SetSession("accessList", accessList)
				m.UpdateLoginTime(user.Id)

				err := ""
				code := 0
				data := fmt.Sprintf(`{"code":%d,"err":"%s"}`, code, err)
				this.Ctx.WriteString(data)
				return
			}
		} else {
			err := "用户或密码错误！"
			code := -1
			data := fmt.Sprintf(`{"code":%d,"err":"%s"}`, code, err)
			this.Ctx.WriteString(data)
			return
		}

	}
	userInfo := this.GetSession("userinfo")
	if userInfo != nil {
		this.Ctx.Redirect(302, "index.html")
	}
	this.TplName = "login.html"
}
func (this *MainController) Logout() {
	this.DelSession("userinfo")
	this.Ctx.Redirect(302, "login.html")
}

//主页
func (this *MainController) MainFrame() {
	this.Left()
	this.Layout = "template/layer.tpl"

	//LEFT 个人信息
	userinfo := this.GetSession("userinfo").(m.User)
	groupinfo, _ := m.GetGroupById(userinfo.Group.Id)
	role := m.GetRoleByUserId(userinfo.Id)
	this.Data["role"] = role
	this.Data["group"] = groupinfo
	this.Data["session"] = userinfo

	groupid := userinfo.Group.Id
	groupstr := m.GetTreeGroup(m.GroupCate(), groupid)
	groupstr = append(groupstr, groupid)

	// MAIN 七天数据
	//ids := []int{0, 2}
	//alluser := m.GetCountNum(ids)
	alluser := m.GetCountNum(groupstr)
	if alluser == 0 {
		alluser = 1
	}

	// this.TplName = "main.html"
	this.TplName = "template/layer.tpl"
	this.TplName = "newIndex.html"
}
func (this *MainController) Center() {
	this.TplName = "background/center.html"
}

func (this *MainController) Top() {
	userInfo := this.GetSession("userinfo")

	if userInfo == nil {
		this.Ctx.Redirect(302, beego.AppConfig.String("auth_gateway"))
	} else {
		role := m.GetRoleByUserId(userInfo.(m.User).Id)
		roleName := role.Name
		this.Data["userinfo"] = userInfo
		this.Data["roleName"] = roleName
	}
	this.TplName = "background/top.html"
}

func (this *MainController) Left() {
	userInfo := this.GetSession("userinfo")
	if userInfo == nil {
		this.Ctx.Redirect(302, beego.AppConfig.String("auth_gateway"))
	} else {
		role := m.GetRoleByUserId(userInfo.(m.User).Id)
		tree := this.GetResList(userInfo.(m.User).Username, role.Id)

		this.Data["tree"] = &tree

		// fmt.Println(&tree)
	}
	this.TplName = "template/layer.tpl"
}
func (this *MainController) Tab() {
	this.TplName = "background/tab.html"
}

//版本更新页面
func (this *MainController) About() {
	this.TplName = "about/update.html"
}
