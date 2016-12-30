package api

import (
	"fmt"
	"hcloud/models"

	"github.com/astaxie/beego"
)

type DeviceController struct {
	CommonController
}

// 查询设备列表或者单个设备
// uuid 	设备名称[string](列表为空)
// username 用户名[string]
// token 	Token[string]
func (this *DeviceController) DeviceList() {
	data := this.GetString("data")
	if len(data) <= 0 {
		this.Rsp(false, "data为空")
		this.ServeJSON()
		return
	}
	RspJson, err := this.Beejson(data)

	beego.Debug("get the Rspjson:", RspJson)

	uuid := RspJson.Uuid         // js.Get("uuid").MustString("")
	username := RspJson.Username // js.Get("username").MustString("")
	token := RspJson.Token       //js.Get("token").MustString("")
	userinfo, err := models.GetUserByToken(username, token)
	if err != nil {
		this.Data["json"] = nil
		this.ServeJSON()
		return
	}

	if len(uuid) > 0 { //单个设备请求
		devicelist := new(DeviceList)
		device, err := models.GetDeviceByUUID(uuid)
		if err != nil {
			beego.Error("get the device error:", err.Error())
			this.ServeJSON()
			return
		}
		devicelist.Uuid = device.Uuid
		devicelist.Usb = device.Usb
		devicelist.Online = device.Online
		devicelist.Order = device.Order
		devicelist.Remark = device.Remark
		devicelist.Version = device.SdkVersion
		if device.GroupId == 0 {
			devicelist.GroupName = "未分组"
		} else {
			group, _ := models.GetGroupByIdOne(device.GroupId)
			devicelist.GroupName = group.Name
		}
		if device.NickName == "" {
			devicelist.NickName = fmt.Sprintf("%d", device.Id)
		} else {
			devicelist.NickName = device.NickName
		}
		devicelist.Online = 1

		this.Data["json"] = devicelist
		this.ServeJSON()
		return
	}

	var devicelists []*DeviceList
	userdevicelist, err := models.GetIndexUserDeviceByUser(userinfo.Id)
	if err != nil {
		beego.Error("not find the userdevicelist:", userinfo.Id)
		this.ServeJSON()
		return
	}
	for i := 0; i < len(userdevicelist); i++ {
		did := userdevicelist[i]["Did"].(int64)
		device, err := models.GetDeviceListById(did)
		if err != nil {
			beego.Error("get the device error:", err.Error())
		} else {
			devicelist := new(DeviceList)
			devicelist.Uuid = device.Uuid
			devicelist.Usb = device.Usb
			devicelist.NickName = device.NickName
			devicelist.Order = device.Order
			devicelist.Remark = device.Remark
			devicelist.Version = device.SdkVersion
			if device.GroupId == 0 {
				devicelist.GroupName = "未分组"
			} else {
				group, _ := models.GetGroupByIdOne(device.GroupId)
				devicelist.GroupName = group.Name
			}
			if device.NickName == "" {
				devicelist.NickName = fmt.Sprintf("%d", device.Id)
			} else {
				devicelist.NickName = device.NickName
			}
			devicelist.Online = 1
			devicelists = append(devicelists, devicelist)
		}
	}
	this.Data["json"] = devicelists
	this.ServeJSON()
}

// 设备排序
// uuidarr		用户ID[array]
// orderarr		排序位置[array]
func (this *DeviceController) DeviceOrder() {
	data := this.GetString("data")

	requestJson, err := this.BeejsonS(data)

	if err != nil {
		beego.Error("get the data error:", err.Error())
		return
	}

	uuidarr := requestJson.Uuid   // .Get("uuid").MustArray()
	orderarr := requestJson.Order //.Get("order").MustArray()

	var query = false
	if len(uuidarr) != len(orderarr) {
		beego.Error("the uuid lenght don't eq order length")
		query = true
		return
	} else {
		for i := 0; i < len(uuidarr); i++ {
			_, err := models.UpdateDeviceOrder(uuidarr[i], orderarr[i])
			if err != nil {
				beego.Error("order the device error:", err.Error())
				query = true
			}
		}
	}
	if query {
		this.Rsp(false, "排序是失败")
	} else {
		this.Rsp(true, "排序完成")
	}
}

// 删除设备
// uuid		设备UUid[string]
func (this *DeviceController) DeviceDel() {
	data := this.GetString("data")

	requestJson, err := this.BeejsonS(data)
	if err != nil {
		this.Rsp(false, "删除失败")
		beego.Error("get the data error:", err.Error())
		return
	}
	uuid := requestJson.Uuid //js.Get("uuid").MustArray()
	if len(uuid) <= 0 {
		this.Rsp(false, "设备获取失败")
		return
	}

	var flag = false
	for i := 0; i < len(uuid); i++ {
		_, err = models.DeleteDevice(uuid[i])
		if err != nil {
			flag = true
		}
	}

	if flag {
		this.Rsp(false, "删除失败")
	} else {
		this.Rsp(true, "删除成功")
	}
}

// 设备编辑(有值表示修改，没值不修改)
// nickname		设备名称[string]
// groupInt		组别[int64]
// remark		手机备注(string)
func (this *DeviceController) DeviceEdit() {
	data := this.GetString("data")
	if len(data) <= 0 {
		this.Rsp(false, "data为空")
		return
	}

	requestJson, err := this.BeejsonS(data)
	if err != nil {
		beego.Error("get the data error:", err.Error())
		return
	}
	uuid := requestJson.Uuid         //js.Get("uuid").MustArray()
	nickname := requestJson.Nickname //js.Get("nickname").MustString("")
	groupInt := requestJson.GroupInt //js.Get("group").MustInt64()
	remark := requestJson.Remark     // js.Get("remark").MustString("")
	if len(uuid) <= 0 {
		this.Rsp(false, "设备获取失败")
		return
	} else {
		for i := 0; i < len(uuid); i++ {
			device, err := models.GetDeviceByUUID(uuid[i])
			if err != nil {
				this.Rsp(false, "数据获取失败")
				return
			}
			if device.Id != 0 {
				device.Uuid = uuid[i]
				if len(nickname) > 0 {
					device.NickName = nickname
				}
				if groupInt > 0 {
					device.GroupId = groupInt
				}
				if len(remark) > 0 {
					device.Remark = remark
				}
			} else {
				this.Rsp(false, "数据获取失败")
				return
			}
			_, err = models.UpdateMgrDevice(&device)
		}

		if err != nil {
			this.Rsp(false, "修改失败")
		} else {
			this.Rsp(true, "修改成功")
		}
	}
}

func (this *DeviceController) ListDeviceHandler() {
	data := this.GetString("data")

	requestJson, err := this.Beejson(data)
	if err != nil {
		beego.Error("get the data is error:", err.Error())
		return
	}

	// 当前页
	page := requestJson.Page //js.Get("page").MustInt64()
	if page <= 0 {
		page = 0
	}
	// 每页长度
	// pagesize := js.Get("pagesize").MustInt64()
	var pagesize int64 = 5
	// 每页的数据
	num := (page - 1) * pagesize
	devices, pagecount, count, err := models.GetDeviceListByPager(num, pagesize)
	beego.Info(devices)
	if err != nil {
		this.Rsp(false, "查询失败")
	} else {
		datajson := make(map[string]interface{})
		for i := 0; i < len(devices); i++ {
			switch devices[i]["Online"].(int64) {
			case 0:
				devices[i]["OnlineName"] = "不在线"
				devices[i]["Online"] = 0
			case 1:
				devices[i]["OnlineName"] = "在线"
				devices[i]["Online"] = 1
			}
			devices[i]["Fans"] = 0
			g, err := models.GetGroupByIdOne(devices[i]["Group"].(int64))
			if err != nil {
				this.Rsp(false, "查询失败")
				return
			}
			devices[i]["GroupName"] = g.Name
		}
		datajson["devices"] = devices
		datajson["count"] = count
		datajson["pagecount"] = pagecount
		this.Data["json"] = datajson
		this.ServeJSON()
	}
}
