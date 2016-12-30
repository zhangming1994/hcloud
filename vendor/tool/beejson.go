package tool

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
)

//
type RequestJson struct {
	Action   string `json:"action"`
	Username string `json:"username"`
	Token    string `json:"token"`
	Data     string `json:"data"`
}

type SetCouldJson struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Ip       string `json:"ip"`
}

type DelDeviceJson struct {
	Uuid string `json:"uuid"`
}

type Respond struct {
	Status  bool   `json:"status"`
	Info    string `json:"info"`
	Token   string `json:"token"`
	Isadmin bool   `json:"isadmin"`
}

type RequestType int

const (
	Setcould RequestType = iota
	Device
	Deldevice
	Editdevice
)

var (
	requestTypeMap = map[string]RequestType{
		"setcould":   Setcould,   //设置Could
		"device":     Device,     //设备上传
		"deldevice":  Deldevice,  //删除设备
		"editdevice": Editdevice, //编辑设备
	}
)

func NewRequestJson(data string) (RequestJson, error) {
	var request RequestJson
	err := json.Unmarshal([]byte(data), &request)
	if err != nil {
		beego.Error("parse json is", err.Error())
	}
	return request, err
}

func GetRequestType(t string) RequestType {
	return requestTypeMap[t]
}

func (r *RequestJson) SetCouldJson() (SetCouldJson, error) {
	dataJson, err := r.Beejson()
	data := dataJson.(SetCouldJson)
	return data, err
}

func (r *RequestJson) DelDeviceJson() (DelDeviceJson, error) {
	dataJson, err := r.Beejson()
	data := dataJson.(DelDeviceJson)
	return data, err
}

func (r *RequestJson) Beejson() (interface{}, error) {
	switch requestTypeMap[r.Action] {
	case Setcould:
		var request SetCouldJson
		err := json.Unmarshal([]byte(r.Data), &request)
		if err != nil {
			beego.Error("parse json is", err.Error())
		}
		return request, err
		break
	case Device:
		break
	case Deldevice:
		var request DelDeviceJson
		err := json.Unmarshal([]byte(r.Data), &request)
		if err != nil {
			beego.Error("parse json is", err.Error())
		}
		return request, err
		break
	case Editdevice:
		break
	default:
	}
	return nil, errors.New("can't find the requesttype type")
}
