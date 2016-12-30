package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"hcloud/models"
	"strconv"
)

type DeviceController struct {
	CommonController
}

// Mgr设备列表
func (this *DeviceController) List() {
	sEcho := this.GetString("sEcho")
	sSearch_0 := this.GetString("sSearch_0")     // 状态
	clouduserid, _ := this.GetInt64("sSearch_1") // 云控账户
	nickname := this.GetString("sSearch_2")      // 所有者/用户名

	iDisplayStart := this.GetString("iDisplayStart")
	iDisplayLength := this.GetString("iDisplayLength")
	iStart, _ := strconv.Atoi(iDisplayStart)
	iLength, _ := strconv.Atoi(iDisplayLength)

	userInfo := this.GetSession("userinfo").(models.User)
	if this.IsAjax() {

		var statues int64
		if len(sSearch_0) <= 0 {
			statues = -1
		} else {
			statues, _ = strconv.ParseInt(sSearch_0, 10, 64)
		}

		mgrud, count, _ := deviceList(int64(iStart), int64(iLength), &userInfo, statues, clouduserid, nickname)
		data := make(map[string]interface{})
		data["aaData"] = mgrud
		data["iTotalDisplayRecords"] = count
		data["iTotalRecords"] = iLength
		data["sEcho"] = sEcho
		this.Data["json"] = &data
		this.ServeJSON()
	} else {
		this.CommonMenu()
		list, _ := models.GetGroupById(0)
		this.Data["list"] = list
		userInfo := this.GetSession("userinfo").(models.User)
		var groups = make([]int64, 0)
		groupid := userInfo.Group.Id
		groups = append(groups, groupid)
		if group, err := models.GetGroupById(groupid); err == nil {
			if group.Uid == userInfo.Id { //判断组别中是否存在当前用户ID
				groups = models.GetTreeGroup(models.GroupCate(), groupid)
				groups = append(groups, groupid)
			}
		}

		couldlist, _ := models.GetCouldRoleUserByGroup(groups)

		_, allcount, _ := deviceList(int64(iStart), int64(iLength), &userInfo, -1, clouduserid, nickname)
		this.Data["allcount"] = allcount

		_, count1, _ := deviceList(int64(iStart), int64(iLength), &userInfo, 1, clouduserid, nickname)
		this.Data["count1"] = count1

		_, count2, _ := deviceList(int64(iStart), int64(iLength), &userInfo, 2, clouduserid, nickname)
		this.Data["count2"] = count2

		this.Data["couldlist"] = couldlist
		userlist, _ := models.GetUserByFid(groups, &userInfo)
		this.Data["userlist"] = userlist
		this.TplName = "device/list.html"
	}
}

// 分配设备
func (this *DeviceController) AllotDevice() {
	userInfo := this.GetSession("userinfo").(models.User)
	if this.IsPost() {

	} else {

		var groups = make([]int64, 0)
		groupid := userInfo.Group.Id
		groups = append(groups, groupid)
		if group, err := models.GetGroupById(groupid); err == nil {
			if group.Uid == userInfo.Id { //判断组别中是否存在当前用户ID
				groups = models.GetTreeGroup(models.GroupCate(), groupid)
				groups = append(groups, groupid)
			}
		}

		list, _ := models.GetUserByFid(groups, &userInfo)
		this.Data["list"] = list
		this.TplName = "device/allot.html"
	}
}

func (this *DeviceController) GetDevices() {
	userInfo := this.GetSession("userinfo").(models.User)
	var groups = make([]int64, 0)
	groupid := userInfo.Group.Id
	groups = append(groups, groupid)
	if group, err := models.GetGroupById(groupid); err == nil {
		if group.Uid == userInfo.Id { //判断组别中是否存在当前用户ID
			groups = models.GetTreeGroup(models.GroupCate(), groupid)
			groups = append(groups, groupid)
		}
	}

	list, _ := models.GetUserByFid(groups, &userInfo)
	this.Data["json"] = list
	this.ServeJSON()
}

// 分配设备处理
func (this *DeviceController) AllotDeviceHandle() {
	uid, _ := this.GetInt64("uid")
	id := this.GetStrings("ids[]")

	userInfo := this.GetSession("userinfo").(models.User)

	beego.Debug("id and uid", id, uid, len(id))
	var query bool = false
	for i := 0; i < len(id); i++ {
		idInt, _ := strconv.ParseInt(id[i], 10, 64)
		b, err := models.AllotDevice(idInt, uid, &userInfo)
		if !b {
			beego.Error("update the allot error", err)
			query = true
			break
		}
	}
	if query {
		this.Rsp(false, "分配设备失败")
	} else {
		this.Rsp(true, "分配设备成功")
	}
}

// 获取设备列表
// page 当前页数
// pageSize 每页数量
// uid 自己的ID
// statues 设备发放状态
// group 设备属于那个组
// nickname 设备所属姓名
func deviceList(page, pageSize int64, user *models.User, statues, clouduserid int64, nickname string) (mgrud []orm.Params, count int64, err error) {

	mgrud, count, err = models.GetSearchUserList(page, pageSize, user, statues, clouduserid, nickname)
	if err != nil {
		beego.Error("get user list error", err)
	}

	for i := 0; i < len(mgrud); i++ {
		deviceinfo, err := models.GetDeviceListById(mgrud[i]["Did"].(int64))
		if err != nil {
			mgrud[i]["Uuid"] = "-"
			mgrud[i]["DeviceName"] = "-"
			mgrud[i]["Imei"] = "-"
			mgrud[i]["Version"] = "-"
			mgrud[i]["Remark"] = "-"
			mgrud[i]["SdkVersion"] = "-"
		}
		user := models.GetUserById(mgrud[i]["Uid"].(int64))
		mgrud[i]["CloudName"] = user.Username
		mgrud[i]["Uuid"] = deviceinfo.Uuid
		mgrud[i]["DeviceName"] = deviceinfo.NickName
		mgrud[i]["Imei"] = deviceinfo.Imei
		mgrud[i]["Version"] = deviceinfo.Version
		mgrud[i]["Remark"] = deviceinfo.Remark
		mgrud[i]["SdkVersion"] = deviceinfo.SdkVersion
		userinfo := models.GetUserById(mgrud[i]["ToUid"].(int64))
		mgrud[i]["Username"] = userinfo.Username
		mgrud[i]["NickName"] = userinfo.Nickname
		mgrud[i]["Order"] = i + 1
	}
	return
}
