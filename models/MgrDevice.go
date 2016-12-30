package models

import (
	"encoding/json"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 设备信息
type MgrDevice struct {
	Id         int64     //主键
	Uid        int64     //使用者Id
	Uuid       string    `orm:"unique"` //手机UUID
	Serial     string    //手机序号
	Usb        string    `orm:"index"` //Devpath
	Model      string    //手机型号
	SdkVersion string    //SDK版本
	Version    string    //版本
	Abi        string    //二进制接口
	Height     int       //高度
	Width      int       //宽度
	Imei       string    //IMEI
	Online     int       //设备是否在线 0 不在线 1 在线
	NickName   string    //手机备注
	Order      int       //排序
	UpdateTime time.Time //更新时间
	CreateTime time.Time `orm:"type(datetime);auto_now_add"` //时间
	Remark     string    `orm:"type(text)"`                  //备注
	GroupId    int64     `orm:"index" json:"GroupId"`        //组别
	Muid       int64     //管理者Id
	Syslist    *SysList  `orm:"rel(fk)"` //群控系统Id
}

func init() {
	orm.RegisterModel(new(MgrDevice))
}

// add MgrDevice
func AddMgrDevice(MgrDevice *MgrDevice) (int64, error) {
	o := orm.NewOrm()
	num, err := o.Insert(MgrDevice)
	return num, err
}

// update MgrDevice
func UpdateMgrDevice(MgrDevice *MgrDevice) (int64, error) {
	MgrDeviceOrm := make(orm.Params)
	MgrDeviceOrm["Serial"] = MgrDevice.Serial
	MgrDeviceOrm["Usb"] = MgrDevice.Usb
	MgrDeviceOrm["Model"] = MgrDevice.Model
	MgrDeviceOrm["SdkVersion"] = MgrDevice.SdkVersion
	MgrDeviceOrm["Version"] = MgrDevice.Version
	MgrDeviceOrm["Abi"] = MgrDevice.Abi
	MgrDeviceOrm["Height"] = MgrDevice.Height
	MgrDeviceOrm["Width"] = MgrDevice.Width
	MgrDeviceOrm["Imei"] = MgrDevice.Imei
	MgrDeviceOrm["Online"] = MgrDevice.Online
	MgrDeviceOrm["UpdateTime"] = MgrDevice.UpdateTime
	MgrDeviceOrm["NickName"] = MgrDevice.NickName
	MgrDeviceOrm["Remark"] = MgrDevice.Remark
	MgrDeviceOrm["GroupId"] = MgrDevice.GroupId
	MgrDeviceOrm["Muid"] = MgrDevice.Muid
	MgrDeviceOrm["Syslist"] = MgrDevice.Syslist.Id
	num, err := mgrDeviceQuerySeter().Filter("Uuid", MgrDevice.Uuid).Update(MgrDeviceOrm)
	return num, err
}

func UpdateDeviceInfo(uuid, nickname, remark string, groupid int64) (int64, error) {
	MgrDeviceOrm := make(orm.Params)
	MgrDeviceOrm["NickName"] = nickname
	MgrDeviceOrm["Remark"] = remark
	MgrDeviceOrm["GroupId"] = groupid
	num, err := mgrDeviceQuerySeter().Filter("Uuid", uuid).Update(MgrDeviceOrm)
	return num, err
}

// 根据Id查询
func GetMgrDeviceListById(id int64) (MgrDevice MgrDevice, err error) {
	err = mgrDeviceQuerySeter().Filter("Id", id).One(&MgrDevice)
	return MgrDevice, err
}

// query group list by id
func GetMgrDeviceListByGroup(id int64) (list []*MgrDevice, num int64, err error) {
	// if user.UserName == beego.AppConfig.String("admin") {
	num, err = mgrDeviceQuerySeter().Filter("GroupId", id).All(&list)
	// }
	return list, num, err
}

// TO JSON
func ReadMgrDeviceJson(data []byte) MgrDevice {
	var MgrDevice MgrDevice
	err := json.Unmarshal(data, &MgrDevice)
	if err != nil {
		beego.Error("parse json is", err.Error())
	}
	return MgrDevice
}

// query group list by id
func GetDeviceListByGroup(id int64) (list []*MgrDevice, num int64, err error) {
	if id <= 0 {
		num, err = mgrDeviceQuerySeter().All(&list)
	} else {
		num, err = mgrDeviceQuerySeter().Filter("GroupId", id).Limit(-1).All(&list)
	}
	return list, num, err
}

// 根据Id查询
func GetDeviceListById(id int64) (device MgrDevice, err error) {
	err = mgrDeviceQuerySeter().Filter("Id", id).One(&device)
	return device, err
}

// 查询设备是否存在
func CheckDeviceIsExist(uuid string) bool {
	return mgrDeviceQuerySeter().Filter("Uuid", uuid).Exist()
}

// get device by uuid
func GetDeviceByUUID(uuid string) (device MgrDevice, err error) {
	err = mgrDeviceQuerySeter().Filter("Uuid", uuid).Limit(1).One(&device)
	return device, err
}

// 前台用户获取设备
func GetIndexUserDeviceByUser(uid int64) (list []orm.Params, err error) {
	userinfo := GetUserById(uid)
	if userinfo.Username == beego.AppConfig.String("admin_user") {
		_, err = mgrDeviceQuerySeter().Values(&list)
	} else {
		_, err = mgrDeviceQuerySeter().Filter("Uid", uid).Values(&list)
	}
	return list, err
}

// change the device order by uuid
func UpdateDeviceOrder(uuid string, order int) (num int64, err error) {
	num, err = mgrDeviceQuerySeter().Filter("Uuid", uuid).Update(orm.Params{"Order": order})
	return num, err
}

// delete device
func DeleteDevice(uuid string) (num int64, err error) {
	deviceinfo, _ := GetDeviceByUUID(uuid)
	if err == nil {
		_, err = DelUserDeviceByDid(deviceinfo.Id)
	}
	num, err = mgrDeviceQuerySeter().Filter("Uuid", uuid).Delete()
	return num, err
}

// 更新设备所属
func UpdateDeviceUid(id, uid int64) (int64, error) {
	return mgrDeviceQuerySeter().Filter("Id", id).Update(orm.Params{"Uid": uid})
}

// 分页查询所有的设备
func GetDeviceListByPager(num int64, pagesize int64) (list []orm.Params, pagenum int64, count int64, err error) {
	mgrDeviceQuerySeter().Limit(pagesize, num).OrderBy("id").Values(&list)
	count, err = mgrDeviceQuerySeter().Count()
	// 页数
	pagenum = (count + pagesize - 1) / pagesize
	return list, pagenum, count, err
}

// Query Seter
func mgrDeviceQuerySeter() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(MgrDevice))
}

// 根据 uid 查询 设备
func GetDeviceCount(uid int64) int64 {
	qs := mgrDeviceQuerySeter().Filter("Uid", uid)
	count, _ := qs.Count()
	return count
}
