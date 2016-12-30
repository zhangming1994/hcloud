package api

import (
	"github.com/astaxie/beego"
	"hcloud/models"
	"tool"
)

type LoginController struct {
	CommonController
}

type FROMTYPE int

const (
	index FROMTYPE = iota + 1
	mgr
)

var (
	fromMap = map[int]FROMTYPE{
		0: index,
		1: mgr,
	}
)

type LoginRespond struct {
	Respond
	Token   string
	Isadmin bool
}

type Respond struct {
	Status bool
	Info   string
}

var (
	InfoMap = map[int]string{
		0:   "登录成功",
		-1:  "获取data数据为空",
		-2:  "用户名为空",
		-3:  "密码为空",
		-4:  "Ip地址为空",
		-5:  "用户校验失败，请重新登录",
		-6:  "您无权限登录",
		-7:  "参数解析失败",
		-8:  "用户名或者密码错误",
		-9:  "用户登录失败",
		-10: "用户设置失败",
		-11: "获取相应资源失败",
		-12: "token值为空",
	}
)

// 登录
func (this *LoginController) Login() {
	data := this.GetString("data")
	lognRespond := new(LoginRespond)
	if len(data) <= 0 {
		lognRespond.Info = InfoMap[-1]
		this.ToJSON(lognRespond)
		return
	}

	RequestJson, err := this.Beejson(data)
	if err != nil {
		lognRespond.Info = InfoMap[-7]
		this.ToJSON(lognRespond)
		beego.Error("get the data is error:", err.Error())
		return
	}
	username := RequestJson.Username //js.Get("username").MustString("")
	if len(username) <= 0 {
		lognRespond.Info = InfoMap[-2]
		this.ToJSON(lognRespond)
		return
	}
	password := RequestJson.Password //js.Get("password").MustString("")
	if len(password) <= 0 {
		lognRespond.Info = InfoMap[-3]
		this.ToJSON(lognRespond)
		return
	}
	userInfo, err := models.Login(username, password)
	if err != nil {
		lognRespond.Info = InfoMap[-8]
		this.ToJSON(lognRespond)
		return
	} else {
		token := this.GetToken()
		userInfo.Token = token
		_, err = models.UpdateUser(&userInfo)
		if err != nil {
			beego.Error(err)
			lognRespond.Info = InfoMap[-9]
			this.ToJSON(lognRespond)
			return
		}

		isadmin := models.CheckUserIsAdmin(userInfo.Id)
		if from := RequestJson.From; mgr == fromMap[from] {
			if !isadmin {
				lognRespond.Info = InfoMap[-6]
				this.ToJSON(lognRespond)
				return
			}
		}
		lognRespond.Info = InfoMap[0]
		lognRespond.Status = true
		lognRespond.Token = userInfo.Token
		lognRespond.Isadmin = isadmin
		this.ToJSON(lognRespond)
	}
	this.ServeJSON()
}

// 设置群控用户密码
func (this *LoginController) SetCouldUserName() {
	data := this.GetString("data")
	beego.Debug("get the data info:", data)
	respond := new(tool.Respond)
	if len(data) <= 0 {
		respond.Info = InfoMap[-1]
		this.ToJSON(respond)
		return
	}

	request, err := tool.NewRequestJson(data)
	if err != nil {
		respond.Info = InfoMap[-7]
		this.ToJSON(respond)
		return
	}

	dataJson, err := request.SetCouldJson()
	if err != nil {
		respond.Info = InfoMap[-7]
		this.ToJSON(respond)
		return
	}

	username := dataJson.Username
	if len(username) <= 0 {
		respond.Info = InfoMap[-2]
		this.ToJSON(respond)
		return
	}
	password := dataJson.Password
	if len(password) <= 0 {
		respond.Info = InfoMap[-3]
		this.ToJSON(respond)
		return
	}
	ip := dataJson.Ip
	if len(ip) <= 0 {
		respond.Info = InfoMap[-4]
		this.ToJSON(respond)
		return
	}

	userInfo, err := models.Login(username, password)
	if err != nil {
		this.Data["json"] = nil
		this.ServeJSON()
		return
	}
	beego.Debug("get the userInfo", userInfo, userInfo.Role.Id)
	if couldrole, _ := beego.AppConfig.Int64("couldrole"); userInfo.Role.Id != couldrole {
		respond.Info = InfoMap[-5]
		this.ToJSON(respond)
		return
	} else {
		syslist := new(models.SysList)
		syslist.Token = this.GetToken()
		syslist.Uid = userInfo.Id
		syslist.Ip = ip

		_, err := models.AddSysList(syslist)
		if err != nil {
			respond.Info = InfoMap[-10]
			this.ToJSON(respond)
			return
		}
		respond.Info = InfoMap[0]
		respond.Status = true
		respond.Token = syslist.Token
		this.ToJSON(respond)
		return
	}
	this.Ctx.WriteString("")
}
