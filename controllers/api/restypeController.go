package api

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"hcloud/models"
)

type RestypeController struct {
	CommonController
}
type RequestResource struct {
	Respond
	RequestNum    int64 `json:"num"`
	RequestType   int64 `json:"restype"`
	RequestPerson int64 `json:"person"`
}

var (
	ResMap = map[int]string{
		-1: "获取data数据为空",
		-2: "用户名为空",
		-3: "获取相应资源失败",
		-4: "token值为空",
		-5: "解析数据失败",
		-6: "请求数量为空",
		-7: "请求人为空",
		-8: "请求资源类型为空",
	}
)

func (this *RestypeController) RequestResJson() {
	data := this.GetString("data") //username和token
	requestjson := new(RequestJsonC)
	if len(data) <= 0 {
		requestjson.Info = ResMap[-1]
		this.ToJSON(requestjson)
		return
	}

	RequestJson, err := this.Beejson(data)
	if err != nil {
		requestjson.Info = ResMap[-5]
		this.ToJSON(requestjson)
		beego.Error("parsing the data is error:", err.Error())
		return
	}
	if RequestJson.Username == "" {
		requestjson.Info = ResMap[-2]
		this.ToJSON(requestjson)
		beego.Error("username is null:", err.Error())
		return
	}
	if RequestJson.Token == "" {
		requestjson.Info = ResMap[-4]
		this.ToJSON(requestjson)
		beego.Error("token is null", err.Error())
		return
	}
	username := RequestJson.Username            //用户名称
	userinfo := models.GetUserByUname(username) //用户信息
	var reslist []orm.Params
	if username == "admin" {
		reslist, err = models.GetAllResourceType()
		if err != nil {
			requestjson.Info = ResMap[-3]
			beego.Error("get data is error", err.Error())
			this.ToJSON(requestjson)
			return
		}
	} else {
		reslist, err = models.GetPermissionType(userinfo.Id)
		if err != nil {
			requestjson.Info = ResMap[-3]
			beego.Error("get data is error", err.Error())
			this.ToJSON(requestjson)
			return
		}
	}
	this.ToJSON(reslist)
}

//请求资源
func (this *RestypeController) RequestPhoneResource() {
	data := this.GetString("data") //username和token
	requestresource := new(RequestResource)
	if len(data) <= 0 {
		requestresource.Info = ResMap[-1]
		this.ToJSON(requestresource)
		return
	}

	RequestJson, err := this.Beejson(data)
	if err != nil {
		requestresource.Info = ResMap[-5]
		this.ToJSON(requestresource)
		beego.Error("parsing the data is error:", err.Error())
		return
	}
	if RequestJson.RequestNum == 0 {
		requestresource.Info = ResMap[-6]
		this.ToJSON(requestresource)
		beego.Error("request num is null:", err.Error())
		return
	}
	if RequestJson.RequestPerson == 0 {
		requestresource.Info = ResMap[-7]
		this.ToJSON(requestresource)
		beego.Error("request num is null:", err.Error())
		return
	}
	if RequestJson.RequestType == 0 {
		requestresource.Info = ResMap[-8]
		this.ToJSON(requestresource)
		beego.Error("request num is null:", err.Error())
		return
	}
}
