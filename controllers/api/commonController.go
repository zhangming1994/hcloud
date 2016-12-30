package api

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"hcloud/controllers"
)

type CommonController struct {
	controllers.CommonController
}

type RequestJson struct {
	RequestJsonC
	UUID
	RequestResource
	Page     int64  `json:"page"`
	Password string `json:"password"`
	From     int    `json:"from"`
}

type RequestJsons struct {
	RequestJsonC
	UUIDS
	RequestResource
	Order    []int  `json:"order"`
	Nickname string `json"nickname"`
	GroupInt int64  `json:"group"`
	Remark   string `json:"remark"`
}

type RequestJsonC struct {
	Respond
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}

type UUID struct {
	Uuid string `json:"uuid"`
}

type UUIDS struct {
	Uuid []string `json:"uuid"`
}

type UserLogin struct {
	Username string `json:"username"`
}

// 设备JSON结构
type DeviceList struct {
	Uuid      string
	Usb       string
	Online    int
	NickName  string
	Order     int
	Remark    string
	GroupName string
	Version   string
}

func (this *CommonController) Beejson(data string) (RequestJson, error) {
	var request RequestJson
	err := json.Unmarshal([]byte(data), &request)
	if err != nil {
		beego.Error("parse json is", err.Error())
	}
	return request, err
}

func (this *CommonController) BeejsonS(data string) (RequestJsons, error) {
	var requests RequestJsons
	err := json.Unmarshal([]byte(data), &requests)
	if err != nil {
		beego.Error("parse json is", err.Error())
	}
	return requests, err
}

func (this *CommonController) ToJSON(v interface{}) {
	this.Data["json"] = v
	this.ServeJSON()
}

func (this *CommonController) Rsp(status bool, str string) {
	this.Data["json"] = &map[string]interface{}{"status": status, "info": str}
	this.ServeJSON()
}

// 获取令牌
func (this *CommonController) GetToken() string {
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(time.Now().Unix(), 10))
	token := fmt.Sprintf("%x", h.Sum(nil))
	return token
}
