package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/tealeg/xlsx"
	m "hcloud/models"
	"os"
	"strconv"
	"strings"
	"time"
	"tool"
)

type UserController struct {
	CommonController
}

var downloadurl = beego.AppConfig.String("downloadImg")

func (this *UserController) List() {
	sEcho := this.GetString("sEcho")
	status, _ := this.GetInt("sSearch_0")       // 状态
	username := this.GetString("sSearch_1")     // 用户名
	sUsergroup, _ := this.GetInt64("sSearch_2") // 团队
	roleid, _ := this.GetInt64("sSearch_3")     // 角色
	companyname := this.GetString("sSearch_4")  // 公司名
	iDisplayStart := this.GetString("iDisplayStart")
	iDisplayLength := this.GetString("iDisplayLength")
	iStart, _ := strconv.Atoi(iDisplayStart)
	iLength, _ := strconv.Atoi(iDisplayLength)
	userinfo := this.GetSession("userinfo").(m.User)
	if this.IsAjax() {

		//以下判断组别
		var groups = make([]int64, 0)
		var usergroup = make([]int64, 0)

		Cate, _ := m.Groups()
		if userinfo.Username == beego.AppConfig.String("admin_user") {
			for i := 0; i < len(Cate); i++ {
				groups = append(groups, Cate[i].Id)
			}
		} else {
			groupid := userinfo.Group.Id
			groups = m.GetTreeGroup(Cate, groupid)
			groups = append(groups, groupid)
		}

		if sUsergroup > 0 {
			Cate, _ := m.Groups()
			usergroup = m.GetTreeGroup(Cate, sUsergroup)
			usergroup = append(usergroup, sUsergroup)
		}
		//以下判断组别 结束
		//判断组别 结束
		users, count := m.GetUserList(int64(iStart), int64(iLength), groups, username, userinfo.Id, userinfo.Username, status, usergroup, roleid, companyname)
		this.SetSession("userList", users)
		for _, user := range users {

			u := m.GetUserById(user["Fid"].(int64))
			if u.Id == 0 {
				user["Createname"] = ""
			} else {
				user["Createname"] = u.Nickname
			}
			// 组别显示
			role := m.GetRoleByUserId(user["Id"].(int64))
			user["Rolename"] = role.Name
			group, _ := m.GetGroupById(user["Group"].(int64))
			user["Groupname"] = group.Name
			user["Deviceconut"] = m.GetDeviceCount(user["Id"].(int64))
			user["CreateTimeStr"] = user["Createtime"].(time.Time).Format("2006-01-02 15:04:05")
			user["LastLoginTimeStr"] = user["Lastlogintime"].(time.Time).Format("2006-01-02 15:04:05")

		}
		data := make(map[string]interface{})
		data["aaData"] = users
		data["iTotalDisplayRecords"] = count
		data["iTotalRecords"] = iLength
		data["sEcho"] = sEcho
		this.Data["json"] = &data
		this.ServeJSON()
	} else {
		var gs []int64
		var uid, rid int64

		_, allcount := m.GetUserList(int64(iStart), int64(iLength), gs, "", uid, userinfo.Username, 0, gs, rid, "")
		this.Data["allCount"] = allcount
		_, count1 := m.GetUserList(int64(iStart), int64(iLength), gs, "", uid, userinfo.Username, 1, gs, rid, "")
		this.Data["count1"] = count1
		_, count2 := m.GetUserList(int64(iStart), int64(iLength), gs, "", uid, userinfo.Username, 2, gs, rid, "")
		this.Data["count2"] = count2
		this.CommonController.CommonMenu()

		userinfo := this.GetSession("userinfo").(m.User)
		groupid := userinfo.Group.Id
		var groups = make([]*m.Group, 0)

		group, err := m.GetGroupById(groupid) //进行组别的查询
		if err != nil {
			beego.Error(err)
		}
		if group.Uid == userinfo.Id { //判断组别中是否存在当前用户ID
			groups = m.GetTree(m.GroupCate(), groupid)
			groups = append(groups, &group)
		} else if beego.AppConfig.String("admin_user") == userinfo.Username {
			groups = m.GetTree(m.GroupCate(), 0)
		} else {
			beego.Error("不是管理")
		}

		roleList, _ := m.GetAllRole()

		this.Data["roleList"] = roleList
		this.Data["groups"] = groups
		this.CommonMenu()
		this.TplName = "user/index.html"
	}
}

//更新用户
func (this *UserController) Edit() {
	id, _ := this.GetInt64("id")
	userInfo := m.GetUserById(id)
	if this.IsPost() {
		username := this.GetString("username")
		nickname := this.GetString("nickname")
		status, _ := this.GetInt("status")
		empnum := this.GetString("empnum")
		companyname := this.GetString("companyname")
		remark := this.GetString("remark")
		groupid, _ := this.GetInt64("groupid")
		roleid, _ := this.GetInt64("roleid")
		officephone := this.GetString("officephone")
		familyphone := this.GetString("familyphone")
		telphone := this.GetString("telphone")
		email := this.GetString("email")
		wxorqq := this.GetString("wxorqq")
		zipcode := this.GetString("zipcode")
		address := this.GetString("address")
		if len(username) <= 0 {
			this.Rsp(false, "用户名称不能为空")
			return
		}
		if len(nickname) <= 0 {
			this.Rsp(false, "用户姓名不能为空")
			return
		}
		if groupid <= 0 {
			this.Rsp(false, "请选择所属团队")
			return
		}
		if roleid <= 0 {
			this.Rsp(false, "请选择员工职位")
			return
		}
		userInfo.Username = username
		userInfo.Nickname = nickname
		userInfo.Status = status
		userInfo.Empnum = empnum
		userInfo.Companyname = companyname
		userInfo.Fid = this.GetSession("userinfo").(m.User).Id
		userInfo.Remark = remark
		userInfo.Officephone = officephone
		userInfo.Familyphone = familyphone
		userInfo.Telphone = telphone
		userInfo.Email = email
		userInfo.Wxorqq = wxorqq
		userInfo.Zipcode = zipcode
		userInfo.Address = address
		userInfo.Group = &m.Group{Id: groupid}
		userInfo.Role = &m.Role{Id: roleid}
		if userInfo.Protected == 1 {
			this.Rsp(false, "该账号处于保护状态下，无法更改！")
		} else {
			_, err := m.UpdateUser(&userInfo)
			if err != nil {
				this.Rsp(false, "修改用户失败")
			} else {
				this.Rsp(true, "修改用户成功")
			}
		}
	} else {
		roles, _ := m.GetAllRole()
		roleList := make([]m.Role, len(roles))
		for k, role := range roles {
			roleList[k].Id = role["Id"].(int64)
			roleList[k].Name = role["Name"].(string)
			roleList[k].Description = role["Description"].(string)
			roleList[k].Status = role["Status"].(int64)
			if roleList[k].Status == 2 {
				roleList[k].Isnormal = true
			} else {
				roleList[k].Isnormal = false
			}

		}
		this.Data["RoleList"] = &roleList

		groups := m.GetTree(m.GroupCate(), 0)
		this.Data["groups"] = groups
		this.Data["groupname"], _ = this.GetInt64("Groupname")

		this.CommonMenu()
		if userInfo.Imgurl == "" {
			this.Data["imgurl"] = ""
		} else {
			this.Data["imgurl"] = downloadurl + userInfo.Imgurl
		}
		this.Data["user"] = userInfo
		this.Data["RoleId"] = userInfo.Role.Id
		this.Data["GroupId"] = userInfo.Group.Id
		this.TplName = "user/edit.html"
	}

}

// 添加用户
func (this *UserController) Add() {
	if this.IsPost() {
		username := this.GetString("username")
		password := this.GetString("password")
		nickname := this.GetString("nickname")
		status, _ := this.GetInt("status")
		empnum := this.GetString("empnum")
		companyname := this.GetString("companyname")
		remark := this.GetString("remark")
		groupid, _ := this.GetInt64("groupid")
		createtime := time.Now()
		lastlogintime := time.Now()
		roleid, _ := this.GetInt64("roleid")
		officephone := this.GetString("officephone")
		familyphone := this.GetString("familyphone")
		telphone := this.GetString("telphone")
		email := this.GetString("email")
		wxorqq := this.GetString("wxorqq")
		zipcode := this.GetString("zipcode")
		address := this.GetString("address")
		if len(username) <= 0 {
			this.Rsp(false, "用户名称不能为空")
			return
		}
		if len(nickname) <= 0 {
			this.Rsp(false, "用户姓名不能为空")
			return
		}
		if len(password) <= 0 {
			this.Rsp(false, "用户密码不能为空")
			return
		}
		if groupid <= 0 {
			this.Rsp(false, "请选择所属团队")
			return
		}
		if roleid <= 0 {
			this.Rsp(false, "请选择员工职位")
			return
		}
		user := m.GetUserByUname(username)
		if user.Id == 0 {
			user := new(m.User)
			user.Username = username
			user.Password = tool.EncodeUserPwd(password)
			user.Nickname = nickname
			user.Status = status
			user.Empnum = empnum
			user.Companyname = companyname
			user.Fid = this.GetSession("userinfo").(m.User).Id
			user.Remark = remark
			user.Officephone = officephone
			user.Familyphone = familyphone
			user.Telphone = telphone
			user.Email = email
			user.Wxorqq = wxorqq
			user.Zipcode = zipcode
			user.Address = address
			user.Group = &m.Group{Id: groupid}
			user.Role = &m.Role{Id: roleid}
			user.Createtime = createtime
			user.Lastlogintime = lastlogintime
			_, err := m.AddUser(user)
			if err != nil {
				beego.Error(err)
				this.Rsp(false, "添加用户失败")
			} else {
				this.Rsp(true, "添加用户成功")
			}
		} else {
			this.Rsp(false, "用户已经存在")
		}
	} else {
		groups := m.GetTree(m.GroupCate(), 0)
		this.Data["groups"] = groups

		roles, _ := m.GetAllRole()
		roleList := make([]m.Role, len(roles))
		for k, role := range roles {
			roleList[k].Id = role["Id"].(int64)
			roleList[k].Name = role["Name"].(string)
			roleList[k].Description = role["Description"].(string)
			roleList[k].Status = role["Status"].(int64)
			if roleList[k].Status == 2 {
				roleList[k].Isnormal = true
			} else {
				roleList[k].Isnormal = false
			}

		}

		this.Data["RoleList"] = &roleList

		this.CommonMenu()
		this.TplName = "user/add.html"
	}

}

// 删除用户
func (this *UserController) Delete() {
	ids := this.GetStrings("ids[]")
	err := m.DelUser(ids)
	if err != nil {
		this.Rsp(false, "删除用户失败")
	} else {
		this.Rsp(true, "删除用户成功")
	}

}

func (this *UserController) Allocation() {
	userId, _ := this.GetInt64("userId")
	roleId, _ := this.GetInt64("roleId")

	m.DelUserRoleByUserId(userId)

	_, err := m.AddRoleUser(roleId, userId)

	if err != nil {
		this.Rsp(false, "分配权限失败")
	} else {
		this.Rsp(true, "分配权限成功")
	}
}

// 禁用/ 启用用户
func (this *UserController) EditUserStatus() {
	user := this.GetSession("userinfo").(m.User)
	ids := this.GetStrings("ids[]")
	status, _ := this.GetInt64("status")
	for i := 0; i < len(ids); i++ {
		id, _ := strconv.ParseInt(ids[i], 10, 64)
		if user.Id == id {
			this.Rsp(false, "无法改变自己状态！")
			return
		}
	}
	// 禁用 / 启用
	err := m.EditUserStatus(ids, status)
	if err != nil {
		this.Rsp(false, "状态修改失败")
	} else {
		this.Rsp(true, "状态修改成功")
	}

}

// 重置密码
func (this *UserController) ResetPassword() {
	ids := this.GetStrings("ids[]")
	password := beego.AppConfig.String("defaultpassword")
	var err error
	for i := 0; i < len(ids); i++ {
		id, _ := strconv.ParseInt(ids[i], 10, 64)
		user := m.GetUserById(id)
		if user.Id == 0 {
			this.Rsp(false, "获取用户失败")
			return
		}
		user.Password = tool.EncodeUserPwd(password)
		_, err = m.UpdateUser(&user)
	}
	if err != nil {
		this.Rsp(false, "重置密码失败")
	} else {
		this.Rsp(true, "重置密码成功")
	}
}

// 修改密码
func (this *UserController) EditPassword() {
	id, _ := this.GetInt64("id")
	if this.IsPost() {
		password := this.GetString("password")
		user := m.GetUserById(id)
		if user.Id == 0 {
			this.Rsp(false, "获取用户失败")
			return
		}
		user.Password = tool.EncodeUserPwd(password)
		_, err := m.UpdateUser(&user)
		beego.Debug(err)
		if err != nil {
			this.Rsp(false, "修改密码失败")
		} else {
			this.Rsp(true, "修改密码成功")
		}
	} else {
		this.TplName = "user/password.html"
	}
}

// 查看用户
func (this *UserController) ViewUser() {
	id, _ := this.GetInt64("id")
	user := m.GetUserById(id)
	switch int64(user.Status) {
	case 1:
		user.Statusname = "正常"
	case 2:
		user.Statusname = "禁用"
	}
	user.Createtimestr = user.Createtime.Format("2006-01-02 15:04:05")
	user.Lastlogintimestr = user.Lastlogintime.Format("2006-01-02 15:04:05")
	// 所属团队
	group, err := m.GetGroupById(user.Group.Id)
	if err != nil {
		beego.Error(err)
		this.Data["GroupName"] = ""
	} else {
		this.Data["GroupName"] = group.Name
	}
	// 职位
	err, role := m.GetRoleById(user.Role.Id)
	if err != nil {
		beego.Error(err)
		this.Data["RoleName"] = ""
	} else {
		this.Data["RoleName"] = role.Name
	}

	if user.Imgurl == "" {
		this.Data["imgurl"] = ""
	} else {
		this.Data["imgurl"] = downloadurl + user.Imgurl
	}
	this.Data["user"] = user
	this.Data["Deviceconut"] = m.GetDeviceCount(user.Id)      // 已有设备数
	this.Data["DownTotalCount"] = m.GetAllDownNum(user.Id)    // 已使用资源
	this.Data["UploadTotalCount"] = m.GetPersonTotal(user.Id) // 资源总数
	this.Data["StatusName"] = user.Statusname
	this.Data["CreateTime"] = user.Createtimestr
	this.Data["LastLoginTime"] = user.Lastlogintimestr
	this.CommonMenu()
	this.TplName = "user/view.html"
}

// 个人信息
func (this *UserController) GetPerson() {
	if this.IsAjax() {
		sEcho := this.GetString("sEcho")
		iDisplayStart := this.GetString("iDisplayStart")
		iDisplayLength := this.GetString("iDisplayLength")
		iStart, _ := strconv.ParseInt(iDisplayStart, 10, 64)
		iLength, _ := strconv.ParseInt(iDisplayLength, 10, 64)

		// users, count := this.CutSlice(m.UserCate(), iStart, iLength, 1, this.GetSession("userinfo").(m.User).Id)
		list, count, _ := m.GetUsersByFid(iStart, iLength, this.GetSession("userinfo").(m.User).Id)
		for i := 0; i < len(list); i++ {
			list[i]["Createtimestr"] = list[i]["Createtime"].(time.Time).Format("2006-01-02 15:04:05")
			list[i]["Lastlogintimestr"] = list[i]["Lastlogintime"].(time.Time).Format("2006-01-02 15:04:05")
		}
		data := make(map[string]interface{})
		data["aaData"] = list
		data["iTotalDisplayRecords"] = count
		data["iTotalRecords"] = iLength
		data["sEcho"] = sEcho
		this.Data["json"] = &data
		this.ServeJSON()
	} else {
		id, _ := this.GetInt64("id")
		user := m.GetUserById(id)
		if user.Id == 0 {
			beego.Debug("The user is not exist.")
			return
		}
		switch int64(user.Status) {
		case 1:
			user.Statusname = "正常"
		case 2:
			user.Statusname = "禁用"
		}
		user.Createtimestr = user.Createtime.Format("2006-01-02 15:04:05")
		user.Lastlogintimestr = user.Lastlogintime.Format("2006-01-02 15:04:05")
		// 团队
		group, err := m.GetGroupById(user.Group.Id)
		if err != nil {
			beego.Error(err)
			this.Data["Groupname"] = ""
		} else {
			this.Data["Groupname"] = group.Name
		}
		// 职位
		err, role := m.GetRoleById(user.Role.Id)
		if err != nil {
			beego.Error(err)
			this.Data["Rolename"] = ""
		} else {
			this.Data["Rolename"] = role.Name
		}
		this.Data["user"] = user
		if user.Imgurl == "" {
			this.Data["imgurl"] = ""
		} else {
			this.Data["imgurl"] = downloadurl + user.Imgurl
		}
		this.Data["StatusName"] = user.Statusname
		this.Data["CreateTime"] = user.Createtimestr
		this.Data["LastLoginTime"] = user.Lastlogintimestr
		this.CommonMenu()
		this.TplName = "user/person.html"
	}

}

//上传头像文件
func (this *UserController) UploadImg() {
	userid, _ := this.GetInt64("userid")
	_, fh, err := this.GetFile("uploadFile")
	ext := fh.Filename[strings.LastIndex(fh.Filename, ".")+1 : len(fh.Filename)]
	var fileName, imgurl string
	dest := beego.AppConfig.String("uploadImg")
	download := beego.AppConfig.String("downloadImg")
	if ext != "gif" && ext != "jpg" && ext != "jpeg" && ext != "png" && ext != "bmp" {
		this.RspImg(false, "图片类型不正确", imgurl)
		beego.Error("the file ext don't math")
		return
	} else {
		if err != nil {
			imgurl = ""
			this.RspImg(false, "上传图片出错", imgurl)
			beego.Error(err)
			return
		} else {
			nowtime := time.Now().Unix()
			fileName = fmt.Sprintf("%d.%s", nowtime, ext)

			os.MkdirAll(dest, 0755)
			path := dest + fileName
			this.SaveToFile("uploadFile", path)

			_, err := m.UploadImg(fileName, userid)
			imgurl = download + fileName
			if err != nil {
				this.RspImg(false, "上传图片失败", "")
			} else {
				this.RspImg(true, "上传图片成功", imgurl)
			}
		}
	}
}

//导出用户列表
func (this *UserController) ExportUserList() {
	list := this.GetSession("userList")

	if list == nil {
		this.Rsp(false, "无查询数据！")
	} else {

		users := list.([]orm.Params)

		outexcel := "./upload/Excel/user/"
		os.MkdirAll(outexcel, 0755)
		outexcel = outexcel + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
		var file *xlsx.File
		var sheet *xlsx.Sheet
		var row *xlsx.Row
		var cell *xlsx.Cell
		file = xlsx.NewFile()
		sheet, _ = file.AddSheet("Sheet1")

		row = sheet.AddRow()
		cell = row.AddCell()
		cell.Value = "工号"
		cell = row.AddCell()
		cell.Value = "姓名"
		cell = row.AddCell()
		cell.Value = "密码"

		for _, v := range users {
			row = sheet.AddRow()

			cell = row.AddCell()
			cell.Value = strconv.FormatInt(v["Username"].(int64), 10)

			cell = row.AddCell()
			cell.Value = v["Nickname"].(string)

			cell = row.AddCell()
		}

		err := file.Save(outexcel)
		this.Ctx.Output.Download(outexcel)
		os.Remove(outexcel) //删除文件
		if err != nil {
			beego.Error(err)
		}
	}
}
