package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
	"hcloud/models"
	"tool"
)

type CloudApiController struct {
	CommonController
}

type REQUESTTYPE int

const (
	device REQUESTTYPE = iota + 1
	deldevice

	deviceInfo = "上传设备"
)

var (
	requestTypeMap = map[string]REQUESTTYPE{
		"device":    device,    //设备请求
		"deldevice": deldevice, //删除设备
	}

	textErrorMsg = map[int]string{
		0:  "上传成功",
		-1: "数据库插入失败",
		-2: "数据库更新失败",
		-3: "数据解析失败",
		-4: "群控不存在",
		-5: "不是云控账户",
		-6: "未知错误",
		-7: "删除设备失败",
		-8: "更新设备失败",
	}

	serveraddress = beego.AppConfig.String("serveraddress")
)

func (this *CloudApiController) Error(err string) {
	beego.Error("Error:" + err)
	this.Ctx.WriteString("protocol error")
}

const (
	ActionUploadAndStore = "uploadphonedevice" //上传的设备信息
)

func (this *CloudApiController) Index() {
	beego.Debug("download start.............")
	data := this.GetString("data")
	if len(data) <= 0 {
		this.Ctx.WriteString("empty data")
		return
	}
	jsonStr := []byte(data)
	js, err := simplejson.NewJson(jsonStr)
	if err != nil {
		this.Error(fmt.Sprintf("parse json error:%v", err))
		this.Ctx.WriteString("parse json error")
		return
	}
	action := js.Get("action").MustString("")
	switch action {
	case ActionUploadAndStore:
		// this.UploadAndStoreDeviceInfo(js)
	default:
	}
}

func (this *CloudApiController) RequestDownloadResource() {
	num := this.GetString("num")          //请求数量
	types := this.GetString("requestype") //请求资源类型
	team := this.GetString("user")        //请求人
	requestnum, _ := strconv.ParseInt(num, 10, 64)
	resourcetype, _ := strconv.ParseInt(types, 10, 64)
	distriteam, _ := strconv.ParseInt(team, 10, 64)
	var (
		Push      models.PushPhone
		phonelist []int64
	)

	Push.Phonetype = resourcetype
	Push.Num = requestnum
	Push.DistriTeam = distriteam
	Push.ChanData = make(chan models.ChanData)

	e := models.UrlNuclearQueue.Put(&Push)
	if e != nil {
		beego.Error("request data error:", e)
		this.Rsp(false, "推送队列错误")
		return
	} else {
		Cd := <-Push.ChanData
		if Cd.Err != nil {
			beego.Error("phone request error:", Cd.Err.Error())
			this.Rsp(false, Cd.Err.Error())
			return
		} else {
			if Cd.Count <= 0 {
				beego.Error("phone is empty :")
				this.Rsp(false, "资源不足")
				return
			} else {
				for _, cust := range Cd.PhoneList {
					phonelist = append(phonelist, cust)
				}
			}
		}
	}

	this.Data["json"] = phonelist
	this.ServeJSON()
}

func (this *CloudApiController) ApiIndex() {
	respond := new(tool.Respond)
	data := this.GetString("data")
	dataJson, _ := tool.NewRequestJson(data)

	sysList, userInfo, err := models.GetSyslistByUserToken(dataJson.Username, dataJson.Token)
	if err != nil && sysList.Id <= 0 && userInfo.Id <= 0 {
		respond.Status = false
		respond.Info = textErrorMsg[-4]
	} else {
		status, code := false, -6
		switch tool.GetRequestType(dataJson.Action) {
		case tool.Device:
			mgrdevice := this.beeJsonToDevice(dataJson.Data)
			status, code = this.deviceHandle(&mgrdevice, sysList.Id, &userInfo) // .Id, userInfo.Group.Tid, userInfo.Group.Uid
			break
		case tool.Deldevice:
			device, _ := dataJson.DelDeviceJson()
			status, code = this.deldeviceHandle(device.Uuid)
			break
		case tool.Editdevice:
			mgrdevice := this.beeJsonToDevice(dataJson.Data)
			status, code = this.editdeviceHandle(&mgrdevice)
		default:
		}
		respond.Status = status
		respond.Info = textErrorMsg[code]
	}

	this.Data["json"] = respond
	this.ServeJSON()
	return
}

func (this *CloudApiController) deviceHandle(mgrdevice *models.MgrDevice, sysid int64, userinfo *models.User) (bool, int) {
	if rolecould, _ := beego.AppConfig.Int64("couldrole"); userinfo.Role.Id != rolecould {
		return false, -5
	} else {
		beego.Debug("userinfo", userinfo.Group)
		if userinfo.Group.Tid == 0 {
			mgrdevice.Muid = userinfo.Group.Uid
		} else {
			usergroup, err := models.GetGroupById(userinfo.Group.Tid)
			if err != nil {
				beego.Debug("get the top group error", err)
				return false, -1
			}
			mgrdevice.Muid = usergroup.Uid
		}
	}

	mgrdevice.Uid = userinfo.Id
	mgrdevice.Syslist = &models.SysList{Id: sysid}
	if models.CheckDeviceIsExist(mgrdevice.Uuid) {
		_, err := models.UpdateMgrDevice(mgrdevice)
		if err != nil {
			beego.Error("update the device error:", err.Error(), "device uuid is:", mgrdevice.Uuid)
			return false, -2
		}
	} else {
		insertid, err := models.AddMgrDevice(mgrdevice)
		if err != nil {
			return false, -1
		} else {
			userdevice := new(models.MgrUserDevice)
			userdevice.Uid = userinfo.Id
			userdevice.Did = insertid
			models.AddUserDevice(userdevice)
			return true, 0
		}
	}
	return false, -6
}

func (this *CloudApiController) deldeviceHandle(uuid string) (bool, int) {
	_, err := models.DeleteDevice(uuid)
	if err != nil {
		return false, -7
	}
	return true, 0
}

func (this *CloudApiController) editdeviceHandle(mgrdevice *models.MgrDevice) (bool, int) {
	Uuid, NickName, GroupId, Remark := mgrdevice.Uuid, mgrdevice.NickName, mgrdevice.GroupId, mgrdevice.Remark
	beego.Debug("...,", Uuid, NickName, GroupId, Remark)
	_, err := models.UpdateDeviceInfo(Uuid, NickName, Remark, GroupId)
	if err != nil {
		return false, -8
	}
	return true, 0

}

func (this *CloudApiController) beeJsonToDevice(data string) models.MgrDevice {
	var device models.MgrDevice
	err := json.Unmarshal([]byte(data), &device)
	if err != nil {
		beego.Error("parse json is", err.Error())
	}
	return device
}

func (this *CloudApiController) clientError(errinfo string, code int) string {
	return errinfo + textErrorMsg[code]
}
