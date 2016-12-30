package models

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

// 用户设备表
type MgrUserDevice struct {
	Id      int64 // 标识
	Fid     int64 // 父级别ID
	Uid     int64 // 用户ID
	Did     int64 // 设备ID
	ToUid   int64 // 分配的Uid
	Statues int   // 是否分配[0：未分配,1：分配]

	Order      int64  `orm:"-"` //序号
	GroupName  string `orm:"-"` //组别名
	NickName   string `orm:"-"` //用户姓名
	UserName   string `orm:"-"` //用户帐号
	Uuid       string `orm:"-"` //设备UUID
	DeviceName string `orm:"-"` //设备名称
}

// 初始化表
func init() {
	orm.RegisterModel(new(MgrUserDevice))
}

// 定义表名
func (this *MgrUserDevice) TableName() string {
	return "mgr_user_device"
}

// 添加
func AddUserDevice(ud *MgrUserDevice) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Insert(ud)
	return id, err
}

// 按照id查询
func GetUserDeviceById(id int64) (ud MgrUserDevice, err error) {
	o := orm.NewOrm()
	ud = MgrUserDevice{Id: id}
	err = o.Read(&ud)
	return ud, err
}

// 修改
func EditUserDevice(ud *MgrUserDevice) (int64, error) {
	o := orm.NewOrm()
	id, err := o.Update(ud)
	return id, err
}

// 根据用户获取设备
func GetUserDeviceByUser(uid int64) (list []orm.Params, err error) {
	o := orm.NewOrm()

	userinfo := GetUserById(uid)
	if userinfo.Fid == 0 {
		_, err = o.QueryTable("mgr_user_device").Filter("Uid", 0).Values(&list)
	} else {
		_, err = o.QueryTable("mgr_user_device").Filter("Uid", uid).Values(&list)
	}
	return list, err
}

// 根据父id获取用户
func GetUserDeviceByFid(fid int64) (list []orm.Params, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("mgr_user_device").Filter("Fid", fid).Values(&list)
	return list, err
}

//获取设备列表
func GetSearchUserList(page int64, page_size int64, user *User, statues, uid int64, nickname string) (list []orm.Params, count int64, err error) {
	var (
		uidList []int
	)
	o := orm.NewOrm()
	qs := o.QueryTable("mgr_user_device")
	cond := orm.NewCondition()

	if statues >= 0 {
		cond = cond.And("Statues", statues)
	}

	//姓名必须是分配后的
	if len(nickname) > 0 && statues != 0 {
		list, _, err = GetUserListByNickname(nickname)
		for i := 0; i < len(list); i++ {
			uidList = append(uidList, int(list[i]["Id"].(int64)))
		}
	}

	if len(uidList) > 0 {
		cond = cond.And("to_uid__in", uidList)
	}

	if uid > 0 {
		cond = cond.And("Uid", uid)
	}

	if user.Username == beego.AppConfig.String("admin_user") {
		cond = cond.And("Fid", 0)
	} else {
		cond = cond.And("Uid", user.Id)
	}

	// beego.Debug("cond", didList, uidList)
	_, err = qs.Limit(page_size, page).SetCond(cond).Values(&list)
	count, _ = qs.SetCond(cond).Count()
	return list, count, err

}

// 查询名下所有设备
func GetAllStatusDevice(status []int) (list []*MgrUserDevice, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("mgr_user_device").Filter("statues__in", status).All(&list)
	return list, err
}

// 分配设备
// id mgr_user_device的主键
// touid 分发到用户的id
// user 登录用户的user
func AllotDevice(id, touid int64, user *User) (b bool, err error) {
	ud, err := GetUserDeviceById(id)
	o := orm.NewOrm()
	err = o.Begin()
	if err != nil {
		err = o.Rollback()
		return false, errors.New("get the userdevice error")
	} else {
		// 设备已分配
		if ud.Statues == 1 {
			_, err = o.QueryTable("mgr_user_device").Filter("Id__gt", id).Filter("Did", ud.Did).Delete()
			if err != nil {
				err = o.Rollback()
				beego.Error("Rollback DB insrt mgr_userdevice", id, err.Error())
				return false, err
			}
		}

		if ud.Uid != touid {
			mgrud := new(MgrUserDevice)
			mgrud.Fid = id
			mgrud.Uid = touid
			mgrud.Did = ud.Did
			mgrud.Statues = 0
			_, err = o.Insert(mgrud)
			if err != nil {
				err = o.Rollback()
				beego.Error("Rollback DB insrt mgr_userdevice", id, err.Error())
				return false, err
			} else {
				_, err = o.QueryTable("mgr_user_device").Filter("Id", id).Update(orm.Params{"Statues": 1, "ToUid": touid})
				if err != nil {
					err = o.Rollback()
					beego.Error("Rollback DB update mgr_userdevice", id, err.Error())
					return false, err
				} else {
					_, err = UpdateDeviceUid(ud.Did, touid)
					err = o.Commit()
					return true, nil
				}
			}
		} else {
			_, err = o.QueryTable("mgr_user_device").Filter("Id", id).Update(orm.Params{"Statues": 0, "ToUid": 0})
			if err != nil {
				err = o.Rollback()
				beego.Error("Rollback DB update mgr_userdevice", id, err.Error())
				return false, err
			} else {
				_, err = UpdateDeviceUid(ud.Did, touid)
				err = o.Commit()
				return true, nil
			}
		}

	}
	err = o.Rollback()
	return false, errors.New("unkown error")
}

// 更新设备的分配状态
func UpdateDeviceStatus(id int64, status int) (num int64, err error) {
	o := orm.NewOrm()
	num, err = o.QueryTable("mgr_user_device").Filter("Id", id).Update(orm.Params{"Statues": status})
	return num, err
}

// 根据设备删除Did
func DelUserDeviceByDid(did int64) (num int64, err error) {
	o := orm.NewOrm()
	num, err = o.QueryTable("mgr_user_device").Filter("Did", did).Delete()
	return num, err
}
