package controllers

import (
	"github.com/astaxie/beego"
	m "hcloud/models"

	"strconv"
)

type GroupController struct {
	CommonController
}

// 显示组别
func (this *GroupController) List() {
	if this.IsAjax() {
		sEcho := this.GetString("sEcho")
		iDiaplayStart := this.GetString("iDisplayStart")
		iDisplayLength := this.GetString("iDisplayLength")
		iStart, _ := strconv.Atoi(iDiaplayStart)
		iLength, _ := strconv.Atoi(iDisplayLength)
		status, _ := this.GetInt64("sSearch_0")  // 状态
		groupname := this.GetString("sSearch_1") // 用户名

		user := this.GetSession("userinfo").(m.User)
		tempCate, _ := m.GroupsByCond(status, groupname)

		groups, count := this.CutGroupSlice(tempCate, int64(iStart), int64(iLength), 1, user.Group.Id)

		grouplist := groups.([]*m.Group)
		for _, item := range grouplist {
			user := m.GetUserById(item.Uid)
			if user.Id != 0 {
				item.ManagerName = user.Nickname + " / " + user.Username
			} else {
				item.ManagerName = ""
			}
		}

		data := make(map[string]interface{})
		data["aaData"] = &grouplist
		data["iTotalDisplayRecords"] = count
		data["iTotalRecords"] = iLength
		data["sEcho"] = sEcho
		this.Data["json"] = &data
		this.ServeJSON()
	} else {
		this.CommonController.CommonMenu()
		this.TplName = "group/index.html"
	}
}

// 增加组别
func (this *GroupController) Add() {
	if this.IsPost() {
		fid, _ := this.GetInt64("fid")
		groupname := this.GetTrimString("groupname")
		uid, _ := this.GetInt64("uid")
		status, _ := this.GetInt64("status")
		remark := this.GetTrimString("remark")
		if len(groupname) <= 0 {
			this.Rsp(false, "团队名称不能为空")
			return
		}
		group := new(m.Group)
		group.Name = groupname
		group.Uid = uid
		group.Status = status
		group.Sort = 50
		group.Fid = fid
		group.Remark = remark
		if group.Fid != 1 {
			groupTid, _ := m.GetGroupById(group.Fid)
			if groupTid.Fid == 1 {
				group.Tid = groupTid.Id
			} else {
				group.Tid = groupTid.Tid
			}
		}
		_, err := m.AddGroup(group)
		if err != nil {
			this.Rsp(false, "添加团队失败")
		} else {
			this.Rsp(true, "添加团队成功")
		}
	} else {
		user := this.GetSession("userinfo").(m.User)
		groups := m.GetTree(m.GroupCate(), user.Group.Id)
		group, _ := m.GetGroupById(user.Group.Id)
		if group.Id != 0 {
			groups = append(groups, &group)
			if len(groups) > 2 {
				groups[0], groups[len(groups)-1] = groups[len(groups)-1], groups[0]
			}
		}

		this.Data["groups"] = groups

		var groupids = make([]int64, 0)

		Cate, _ := m.Groups()
		if user.Username == beego.AppConfig.String("admin_user") {
			for i := 0; i < len(Cate); i++ {
				groupids = append(groupids, Cate[i].Id)
			}
		} else {
			groupid := user.Group.Id
			groupids = m.GetTreeGroup(Cate, groupid)
			groupids = append(groupids, group.Id)
		}
		var users []*m.User
		if beego.AppConfig.String("admin_user") == user.Username {
			users, _ = m.GetUsers()
		} else {
			users = m.UserCate(groupids)
		}
		this.Data["adminuser"] = beego.AppConfig.String("admin_user")
		this.Data["users"] = users
		this.CommonMenu()
		this.TplName = "group/add.html"
	}
}

// 更新
func (this *GroupController) Edit() {
	id, _ := this.GetInt64("id")
	group, _ := m.GetGroupById(id)

	if this.IsPost() {
		fid, _ := this.GetInt64("fid")
		groupname := this.GetTrimString("groupname")
		uid, _ := this.GetInt64("uid")
		status, _ := this.GetInt64("status")
		remark := this.GetTrimString("remark")
		group.Name = groupname
		group.Uid = uid
		group.Status = status
		group.Fid = fid
		group.Remark = remark
		if len(groupname) <= 0 {
			this.Rsp(false, "团队名称不能为空")
			return
		}
		_, err := m.UpdateGroup(&group)
		if err != nil {
			this.Rsp(false, "修改团队失败")
		} else {
			this.Rsp(true, "修改团队成功")
		}
	} else {
		user := this.GetSession("userinfo").(m.User)
		groups := m.GetTree(m.GroupCate(), user.Group.Id)
		fgroup, _ := m.GetGroupById(group.Fid)
		if fgroup.Id != 0 {
			groups = append(groups, &fgroup)
			if len(groups) > 2 {
				groups[0], groups[len(groups)-1] = groups[len(groups)-1], groups[0]
			}
		}
		var users []*m.User
		var groupids = make([]int64, 0)

		Cate, _ := m.Groups()
		if user.Username == beego.AppConfig.String("admin_user") {
			for i := 0; i < len(Cate); i++ {
				groupids = append(groupids, Cate[i].Id)
			}
		} else {
			groupid := user.Group.Id
			groupids = m.GetTreeGroup(Cate, groupid)
			groupids = append(groupids, group.Id)
		}
		if beego.AppConfig.String("admin_user") == user.Username {
			users, _ = m.GetUsers()
		} else {
			users = m.UserCate(groupids)
		}
		this.Data["users"] = users
		this.Data["groups"] = groups
		this.Data["group"] = group
		this.Data["adminuser"] = beego.AppConfig.String("admin_user")
		this.CommonMenu()
		this.TplName = "group/edit.html"
	}

}

// 删除组别
func (this *GroupController) Delete() {
	Id, _ := this.GetInt64("id")

	_, err := m.DelGroupById(Id)
	if err != nil {
		this.Rsp(false, "删除团队失败")
	} else {
		this.Rsp(true, "删除团队成功")
	}
}
