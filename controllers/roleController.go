package controllers

import (
	// "github.com/astaxie/beego"
	m "hcloud/models"
	"strconv"
	"strings"
)

type RoleController struct {
	CommonController
}

// 角色列表
func (this *RoleController) List() {

	if this.IsAjax() {
		iDiaplayStart := this.GetString("iDisplayStart")
		iDisplayLength := this.GetString("iDisplayLength")
		iStart, _ := strconv.Atoi(iDiaplayStart)
		iLength, _ := strconv.Atoi(iDisplayLength)

		status, _ := this.GetInt64("sSearch_0") // 状态
		rolename := this.GetString("sSearch_1") // 用户名

		roles, count := m.GetRoleList(int64(iStart), int64(iLength), status, rolename)

		data := make(map[string]interface{})
		data["aaData"] = &roles
		data["iTotalDisplayRecords"] = count
		data["iTotalRecords"] = iLength
		data["sEcho"] = this.GetString("sEcho")

		this.Data["json"] = &data
		this.ServeJSON()
	} else {
		this.CommonController.CommonMenu()
		this.TplName = "role/index.html"
	}

}

// 添加角色
func (this *RoleController) Add() {
	if this.IsPost() {
		name := this.GetString("rolename")
		description := this.GetString("description")
		status, _ := this.GetInt64("status")
		if len(name) <= 0 {
			this.Rsp(false, "角色名称不能为空")
			return
		}
		_, role := m.GetRoleByName(name)
		if role.Id != 0 {
			this.Rsp(false, "角色名称已经存在")
		} else {
			role.Name = name
			role.Description = description
			role.Status = status

			_, err := m.AddRole(&role)
			if err != nil {
				this.Rsp(false, "添加角色失败")
			} else {
				this.Rsp(true, "添加角色成功")
			}
		}
	} else {
		this.CommonMenu()
		this.TplName = "role/add.html"
	}
}

//  删除角色
func (this *RoleController) Delete() {
	Id, _ := this.GetInt64("id")
	_, err := m.DelRoleById(Id)
	if err != nil {
		this.Rsp(false, "删除角色失败")
	} else {
		this.Rsp(true, "删除角色成功")
	}
}

// 更新角色
func (this *RoleController) Edit() {
	id, _ := this.GetInt64("id")

	_, role := m.GetRoleById(id)
	if role.Id == 0 {
		this.Rsp(false, "获取数据失败")
		return
	}
	if this.IsPost() {
		rolename := this.GetString("rolename")
		description := this.GetString("description")
		status, _ := this.GetInt64("status")

		role.Name = rolename
		role.Description = description
		role.Status = status
		_, err := m.UpdateRole(&role)
		if err != nil {
			this.Rsp(false, "修改职位失败")
		} else {
			this.Rsp(true, "修改职位成功")
		}
	} else {
		this.CommonMenu()
		this.Data["role"] = role
		this.TplName = "role/edit.html"
	}
}

// 给用户分配角色
func (this *RoleController) AllotNode() {
	Id, _ := this.GetInt64("id")
	if this.IsPost() {

		err := m.DelRoleRescoursByRoleId(Id)
		if err != nil {
			this.Rsp(false, "分配权限失败")
		} else {
			ids := this.GetString("ids")
			idsstr := strings.Split(ids, ",")
			err = m.AddRoleRescours(Id, idsstr)
			if err != nil {
				this.Rsp(false, "分配权限失败")
			} else {
				this.Rsp(true, "分配权限成功")
			}
		}
	} else {
		var cnt, length int = 0, 0
		ns, _ := m.GetAllNodes()
		nodes := m.GetNodeTreeAndLv(ns, 0)

		for _, v := range nodes {
			if v.Fid == 0 {
				length = length + 1
			}
		}
		tree := make([]Tree, length)
		for k, v := range nodes {
			if v.Fid == 0 {
				k = cnt
				cnt = cnt + 1
				tree[k].Id = v.Id
				tree[k].Url = "/" + v.Url
				tree[k].Text = v.Name
				tree[k].Fid = v.Fid
				// 1代表菜单（目录下面的所有资源）没有把一些不需要的权限去掉
				var childCnt int = 0
				children := make([]Tree, 16)
				for _, v3 := range nodes {
					if v3.Fid == v.Id {
						children[childCnt].Id = v3.Id
						children[childCnt].Text = v3.Name
						children[childCnt].Url = v3.Url
						children[childCnt].Fid = v3.Fid
						childCnt++
					}
				}
				tree[k].Children = make([]Tree, childCnt)
				for k1, v1 := range children {
					if v1.Fid == v.Id {
						tree[k].Children[k1].Fid = v1.Fid
						tree[k].Children[k1].Id = v1.Id
						tree[k].Children[k1].Text = v1.Text
						tree[k].Children[k1].Url = "/" + v1.Url
					}
				}
			}
		}
		checknode, _ := m.GetNodesByRoleId(Id)
		this.Data["nodes"] = checknode
		this.Data["roleid"] = Id
		this.Data["nodetree"] = tree
		this.CommonMenu()
		this.TplName = "role/permission.html"
	}
}

// 给角色分配节点
func (this *RoleController) AllocationRes() {
	Id, _ := this.GetInt64("Id")
	res, _ := m.GetNodeByRoleId(Id)

	// data := make(map[string]interface{})
	this.Data["json"] = &res
	this.ServeJSON()
}
