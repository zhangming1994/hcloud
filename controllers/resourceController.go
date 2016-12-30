package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	m "hcloud/models"
	"io"
	"mime/multipart"
	"os"
	h "resource"
	"strconv"
	"strings"
	"time"
)

type ResourceController struct {
	CommonController
}

//根据日期来计算的显示列表
func (this *ResourceController) Index() {
	allteam, _ := m.GetAllTeamInfo() //取得所有的团队
	if this.IsAjax() {
		sEcho := this.GetString("sEcho")
		iDisplayStart := this.GetString("iDisplayStart")
		iDisplayLength := this.GetString("iDisplayLength")
		iStart, _ := strconv.ParseInt(iDisplayStart, 10, 64)
		iLength, _ := strconv.ParseInt(iDisplayLength, 10, 64)
		search0 := this.GetString("sSearch_0") //时间
		search1 := this.GetString("sSearch_1") //团队
		search2 := this.GetString("sSearch_2") //管理者
		//根据管理者姓名取得id
		var userids []int64
		user := m.GetUserByNames(search2)
		for i := 0; i < len(user); i++ {
			ids, _ := user[i]["Id"].(int64)
			userids = append(userids, ids)

		}
		var starttime string
		var endtime string
		tradedata := strings.Split(search0, " - ")
		if len(tradedata) == 2 {
			starttime = tradedata[0]
			endtime = tradedata[1]
		}

		filter := make(map[string]interface{})
		filter["starttime"] = starttime
		filter["edntime"] = endtime
		filter["s1"] = search1

		list, count, err := m.CloudResourceIndexList(iStart, iLength, filter, userids)

		if err != nil {
			beego.Error("find the data is err", err)
			return
		}
		for _, val := range list {
			val["ResourceType"] = "全部"
			groupid, _ := val["Group"].(int64)
			group, _ := m.GetGroupById(groupid)
			val["Name"] = group.Name //所属团队
			tamadmin := m.GetUserById(group.Uid).Username
			val["UseAdmin"] = tamadmin //使用者
			usepersent := val["UsePersent"]
			persent := usepersent.(string) + `%`
			times := val["ResourceDate"].(time.Time).Format("2006-01-02")
			val["ResourceDate"] = times
			val["UsePersent"] = persent
		}
		data := make(map[string]interface{})
		data["aaData"] = &list
		data["iTotalDisplayRecords"] = count
		data["iTotalRecords"] = iLength
		data["sEcho"] = sEcho
		this.Data["json"] = &data
		this.ServeJSON()
	}
	this.Data["team"] = allteam
	this.CommonMenu()
	this.TplName = "resource/index.html"
}

//折叠列表
func (this *ResourceController) ListDetail() {
	if this.IsAjax() {
		addtime := this.GetString("AddTime")
		name := this.GetString("Name")
		group, err := m.GetGroupName(name) //根据名字取得团队id
		if err != nil {
			beego.Error("get groupid is err", err)
		}
		times, _ := time.Parse("2006-01-02", addtime)
		alluserinfo := m.GetAllUserTeam(group.Id) //拿出该组别下面所有的组员信息
		var users []int64
		for i := 0; i < len(alluserinfo); i++ {
			userids, _ := alluserinfo[i]["Id"].(int64)
			users = append(users, userids)
		}
		list, _ := m.GetDataByTeam(group.Id, times, users) //取得数据
		for _, val := range list {
			showtime := val["ResourceDate"].(time.Time).Format("2006-01-02")
			val["AddTime"] = showtime
			val["Team"] = name
			val["Username"] = m.GetUserById(group.Uid).Username
		}
		this.Data["json"] = list
		this.ServeJSON()
	}
}

//文件上传界面
func (this *ResourceController) UploadIndex() {
	var resourcetype []orm.Params
	var list []orm.Params
	var count int64
	var err error
	var resourcetypeid []int64 //权限集合
	userinfo := this.GetSession("userinfo").(m.User)
	if userinfo.Username == "admin" {
		resourcetype, _ = m.GetAllResourceType()
	} else {
		permissionres, _ := m.GetPermissionType(userinfo.Id)
		//根据资源id取得资源
		for i := 0; i < len(permissionres); i++ {
			typeid, _ := permissionres[i]["CloudResourceType"].(int64)
			resourcetypeid = append(resourcetypeid, typeid)
		}
		if len(resourcetypeid) > 0 {
			resourcetype, _ = m.GetUserPermissionRes(resourcetypeid)
		}
	}
	team, _ := m.GetAllTeamInfo()
	if this.IsAjax() {
		sEcho := this.GetString("sEcho")
		iDisplayStart := this.GetString("iDisplayStart")
		iDisplayLength := this.GetString("iDisplayLength")
		iStart, _ := strconv.ParseInt(iDisplayStart, 10, 64)
		iLength, _ := strconv.ParseInt(iDisplayLength, 10, 64)
		search0 := this.GetString("sSearch_0") //团队
		search1 := this.GetString("sSearch_1") //资源类型
		search2 := this.GetString("sSearch_2") //日期

		var starttime string
		var endtime string
		tradedata := strings.Split(search2, " - ")
		if len(tradedata) == 2 {
			starttime = fmt.Sprintf("%s 00:00:01", tradedata[0])
			endtime = fmt.Sprintf("%s 23:59:59", tradedata[1])
		}
		filter := make(map[string]interface{})
		filter["starttime"] = starttime
		filter["edntime"] = endtime
		filter["s1"] = search1 //资源类型
		filter["s2"] = search0 //所属团队

		//根据登录的用户来取得相应的上传记录
		//1.超级管理员
		var uploaduser []int64
		if userinfo.Username == "admin" {
			list, count, err = m.ListUploads(iStart, iLength, filter)
		} else {
			//多个团队的管理员
			//1根据userinfo取得团队id
			user := m.GetUserById(userinfo.Id)
			//2.根据团队id取得团队信息
			group, _ := m.GetGroupById(user.Group.Id)
			if userinfo.Id == group.Uid {
				//取出该管理员组别下面所有的成员,
				teamuser := m.GetAllUserTeam(userinfo.Id)

				for i := 0; i < len(teamuser); i++ {
					useid, _ := teamuser[i]["Id"].(int64)
					uploaduser = append(uploaduser, useid)
				}
				list, count, err = m.ListUpload(iStart, iLength, filter, uploaduser, resourcetypeid)
			} else {
				//就是单个的用户
				uploaduser = append(uploaduser, userinfo.Id)
				list, count, err = m.ListUpload(iStart, iLength, filter, uploaduser, resourcetypeid)
			}
			if err != nil {
				beego.Error("find the data is err", err)
				return
			}
		}
		for _, val := range list {
			val["UploadDate"] = val["UploadDate"].(time.Time).Format("2006-01-02 15:04:05")
			typesid, _ := val["CloudResourceType"].(int64)
			typesname := m.GetTypeNameByID(typesid)
			val["CloudResourceType"] = typesname.TypeName
			teamnumber, _ := val["Team"].(int64)
			team, _ := m.GetGroupById(teamnumber) //根据team id 取得名称
			val["Team"] = team.Name
			userid, _ := val["UploadUser"].(int64)
			//根据userid取得user
			users := m.GetUserById(userid)
			val["UploadUser"] = users.Username

		}
		data := make(map[string]interface{})
		data["aaData"] = list
		data["iTotalDisplayRecords"] = count
		data["iTotalRecords"] = iLength
		data["sEcho"] = sEcho
		this.Data["json"] = data
		this.ServeJSON()
	}

	this.CommonMenu()
	this.Data["Team"] = team
	this.Data["List"] = resourcetype
	this.TplName = "resource/upload.html"
}

// upload the file
func (this *ResourceController) UploadPhoneResource() {
	f, fileh, err := this.GetFile("file")
	var json = ""
	if err != nil {
		beego.Error("get the file is error:", err.Error())
		return
	} else {
		defer f.Close()
		beego.Debug("start upload file")
		phonetypes := this.GetString("sourcetype")
		phonetype, _ := strconv.ParseInt(phonetypes, 10, 64)
		if phonetype == 0 {
			json = fmt.Sprintf(`{"msg":"%s","code":"%d"}`, "请选择资源类型", -1)
			this.Ctx.WriteString(json)
			// this.Rsp(false, "请选择资源类型")
			return
		} else {
			localpath := "../static/phone/"
			filename, err := this.Upload(fileh, localpath)
			if err != nil {
				// this.Rsp(false, "文件上传错误")
				json = fmt.Sprintf(`{"msg":"%s","code":"%d"}`, "文件上传错误", -1)
				this.Ctx.WriteString(json)
				return
			}
			//the above had already save to local
			uploaduser := this.GetSession("userinfo").(m.User)
			user := m.GetUserByUname(uploaduser.Username) //根据名称取得所属团队
			hu := new(m.CloudUploadRecord)
			resoutype := new(m.CloudResourceType)
			resoutype.Id = phonetype
			hu.UploadName = filename                   //本地文件名
			hu.CloudResourceType = resoutype           //上传资源类型
			hu.UploadUser = uploaduser.Id              //上传人
			hu.Team = user.Group.Id                    //团队
			lastrecordid, err := m.AddUploadRecord(hu) //得到上传批次id
			if err != nil {
				beego.Error("add upload record is err", err)
				json = fmt.Sprintf(`{"msg":"%s","code":"%d"}`, "文件上传失败", -1)
				// this.Rsp(false, "文件上传失败")
			} else {
				h.RunUpTask(lastrecordid, phonetype, uploaduser.Id, localpath, filename)
				// this.Rsp(true, "上传成功，后台处理中......")
				json = fmt.Sprintf(`{"msg":"%s","code":"%d"}`, "上传成功,后台处理中......", 1)
			}
		}
	}
	this.Ctx.WriteString(json)
}

// 上传文件
func (this *ResourceController) Upload(info *multipart.FileHeader, dist string) (string, error) {
	file, err := info.Open()
	defer file.Close()
	fname := info.Filename
	fname = fmt.Sprintf("%s_%s", time.Now().Format("2006_01_02_15_04_05"), fname)
	//定义文件名
	filename := fmt.Sprintf("%s", fname)
	// 创建目录
	err = os.MkdirAll(dist, 0777)
	if err != nil {
		beego.Error(err.Error())
		return "", err
	}
	out, err := os.OpenFile(dist+filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		beego.Error(err.Error())
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	return filename, err
}

//增加资源类型
func (this *ResourceController) AddResourceType() {
	resourcename := this.GetString("name")
	if this.IsPost() {
		if resourcename == "" {
			this.Rsp(false, "请输入资源类型名称")
			return
		} else {
			flag := m.CheckIsRepat(resourcename)
			if flag == true {
				this.Rsp(false, "您输入的资源类型已经存在!")
				return
			} else {
				userid := this.GetSession("userinfo").(m.User).Id
				resourcetype := new(m.CloudResourceType)
				resourcetype.TypeName = resourcename
				resourcetype.AddUser = userid
				resourcetype.AddTime = time.Now()
				resourcetype.Status = 1
				_, err := m.AddResType(resourcetype)
				if err != nil {
					this.Rsp(false, "添加失败")
				} else {
					this.Rsp(true, "添加成功")
				}
			}
		}
	} else {
		this.CommonMenu()
		this.TplName = "resource/typelist.html"
	}
}

// 修改资源类型
func (this *ResourceController) EditResourceType() {
	typename := this.GetString("typename")
	id, _ := this.GetInt64("id")
	_, restype := m.GetTypeById(id)
	if restype.Id == 0 {
		this.Rsp(false, "获取数据失败")
		return
	}
	restype.TypeName = typename
	_, err := m.EditResType(&restype)
	if err != nil {
		this.Rsp(false, "修改资源类型名称失败")
	} else {
		this.Rsp(true, "修改资源类型名称成功")
	}
}

// 修改资源类型状态
func (this *ResourceController) EditResourceTypeStatus() {
	tids := this.GetStrings("ids[]")
	status, _ := this.GetInt64("status")
	var err error
	var ids []int64
	for i := 0; i < len(tids); i++ {
		id, _ := strconv.ParseInt(tids[i], 10, 64)
		ids = append(ids, id)
	}
	err = m.DelTypeUser(ids)
	if err != nil {
		this.Rsp(false, "设置资源类型状态失败")
		return
	}
	for i := 0; i < len(ids); i++ {
		_, restype := m.GetTypeById(ids[i])
		if restype.Id == 0 {
			this.Rsp(false, "获取数据失败")
			return
		}
		restype.Status = status
		_, err = m.EditResType(&restype)
	}
	if err != nil {
		this.Rsp(false, "设置资源类型状态失败")
	} else {
		this.Rsp(true, "设置资源类型状态成功")
	}
}

// 删除资源类型状态
func (this *ResourceController) DelResourceType() {
	tids := this.GetStrings("tids[]")
	var ids []int64
	for i := 0; i < len(tids); i++ {
		id, _ := strconv.ParseInt(tids[i], 10, 64)
		ids = append(ids, id)
	}
	var err error
	err = m.DelType(ids)
	if err != nil {
		this.Rsp(false, "删除资源类型失败")
		return
	}
	err = m.DelTypeUser(ids)
	if err != nil {
		this.Rsp(false, "删除资源类型失败")
	} else {
		this.Rsp(true, "删除资源类型成功")
	}
}

// 资源类型列表
func (this *ResourceController) ListResourceType() {
	sEcho := this.GetString("sEcho")
	iDisplayStart := this.GetString("iDisplayStart")
	iDisplayLength := this.GetString("iDisplayLength")
	iStart, _ := strconv.ParseInt(iDisplayStart, 10, 64)
	iLength, _ := strconv.ParseInt(iDisplayLength, 10, 64)
	status, _ := this.GetInt64("sSearch_0") // 状态
	typename := this.GetString("sSearch_1") // 资源类型名称
	user := this.GetSession("userinfo").(m.User)
	if this.IsAjax() {
		list, count, _ := m.GetTypeByCond(iStart, iLength, user.Username, user.Id, status, typename)
		for i := 0; i < len(list); i++ {
			user := m.GetUserById(list[i]["AddUser"].(int64))
			list[i]["AddUserName"] = user.Username
		}

		data := make(map[string]interface{})
		data["aaData"] = list
		data["iTotalDisplayRecords"] = count
		data["iTotalRecords"] = iLength
		data["sEcho"] = sEcho

		this.Data["json"] = &data
		this.ServeJSON()
	} else {
		_, allcount, _ := m.GetTypeByCond(iStart, iLength, user.Username, user.Id, 0, typename)
		this.Data["allcount"] = allcount
		_, count1, _ := m.GetTypeByCond(iStart, iLength, user.Username, user.Id, 1, typename)
		this.Data["count1"] = count1
		_, count2, _ := m.GetTypeByCond(iStart, iLength, user.Username, user.Id, 2, typename)
		this.Data["count2"] = count2

		this.CommonMenu()
		this.TplName = "resource/typelist.html"
	}
}

// 获取资源类型
func (this *ResourceController) GetResourceType() {
	id, _ := this.GetInt64("id")
	err, restype := m.GetTypeById(id)
	if err != nil {
		this.Data["json"] = nil
	} else {
		this.Data["json"] = restype
	}
	this.ServeJSON()
}

// 加载团队列表
func (this *ResourceController) GetListGroup() {
	if this.IsAjax() {
		sEcho := this.GetString("sEcho")
		// iDiaplayStart := this.GetString("iDisplayStart")
		iDisplayLength := this.GetString("iDisplayLength")
		// iStart, _ := strconv.Atoi(iDiaplayStart)
		iLength, _ := strconv.Atoi(iDisplayLength)
		status, _ := this.GetInt64("sSearch_0")  // 状态
		groupname := this.GetString("sSearch_1") // 用户名

		user := this.GetSession("userinfo").(m.User)
		tempCate, _ := m.GroupsByCond(status, groupname)

		grouplist := m.GetTreeAndLv(tempCate, user.Group.Id, 0)

		data := make(map[string]interface{})
		data["aaData"] = &grouplist
		data["iTotalDisplayRecords"] = len(grouplist)
		data["iTotalRecords"] = iLength
		data["sEcho"] = sEcho
		this.Data["json"] = &data
		this.ServeJSON()
	} else {
		this.CommonController.CommonMenu()
		this.TplName = "resource/typelist.html"
	}
}

// 加载用户列表
func (this *ResourceController) GetListUser() {

	if this.IsAjax() {
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
		this.CommonMenu()
		this.TplName = "resource/typelist.html"
	}
}

// 加载用户团队角色
func (this *ResourceController) GetResourceTypeUsers() {

	userinfo := this.GetSession("userinfo").(m.User)
	var gs []int64
	var uid, rid int64
	data := make(map[string]interface{})
	_, allcount := m.GetUserList(0, -1, gs, "", uid, userinfo.Username, 0, gs, rid, "")
	data["Allusercount"] = allcount
	_, count1 := m.GetUserList(0, -1, gs, "", uid, userinfo.Username, 1, gs, rid, "")
	data["Usercount1"] = count1
	_, count2 := m.GetUserList(0, -1, gs, "", uid, userinfo.Username, 2, gs, rid, "")
	data["Usercount2"] = count2

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

	data["Roles"] = roleList
	data["Groups"] = groups
	this.Data["json"] = data
	this.ServeJSON()
}

//团队资源设置小界面
func (this *ResourceController) TeamResourceSetup() {
	this.TplName = "resource/teamsetup.html"
}

//个人资源设置小界面
func (this *ResourceController) PersonResourceSetup() {
	this.TplName = "resource/personsetup.html"
}

// 根据团队查询用户
func (this *ResourceController) GetUserByGroup() {
	groupids := this.GetStrings("groupids[]")
	var ids []int64
	for i := 0; i < len(groupids); i++ {
		id, _ := strconv.ParseInt(groupids[i], 10, 64)
		ids = append(ids, id)
	}
	list, _ := m.GetUsersByGroup(ids)
	this.Data["json"] = list
	this.ServeJSON()
}

// 设置共享资源类型
func (this *ResourceController) AddTypeandUser() {
	typeids := this.GetStrings("typeids[]")
	userids := this.GetStrings("userids[]")
	var tids, uids []int64
	for i := 0; i < len(typeids); i++ {
		tid, _ := strconv.ParseInt(typeids[i], 10, 64)
		tids = append(tids, tid)
	}
	for i := 0; i < len(userids); i++ {
		uid, _ := strconv.ParseInt(userids[i], 10, 64)
		uids = append(uids, uid)
	}

	err := m.DelTypeUser(tids)
	if err != nil {
		this.Rsp(false, "设置共享资源类型失败")
		return
	}
	_, err = m.AddTypeUser(uids, tids)
	if err != nil {
		this.Rsp(false, "设置共享资源类型失败")
	} else {
		for i := 0; i < len(tids); i++ {
			_, restype := m.GetTypeById(tids[i])
			if restype.Id == 0 {
				this.Rsp(false, "获取数据失败")
				return
			}
			restype.Status = 2
			_, err = m.EditResType(&restype)
		}
		if err != nil {
			this.Rsp(false, "设置共享资源类型失败")
		} else {
			this.Rsp(true, "设置共享资源类型成功")
		}
	}
}

//资源限制设置
func (this *ResourceController) SetUpGroupResource() {
	json := ""
	ids := this.GetString("Id") //团队或者用户Id
	checkid := strings.Split(ids, ",")
	types := this.GetString("types")                 //1.团队 2.个人
	onceresources := this.GetString("oncedown")      //团队用户单次下载量
	onedaytotals := this.GetString("onedayresource") //一天下载量
	onceresource, _ := strconv.ParseInt(onceresources, 10, 64)
	onedaytotal, _ := strconv.ParseInt(onedaytotals, 10, 64)
	for i := 0; i < len(checkid)-1; i++ {
		id, _ := strconv.ParseInt(checkid[i], 10, 64)
		if types == "1" {
			err := m.ModifyGroupResource(id, onedaytotal, onceresource)
			m.ModifyOncePerson(id, onceresource) //团队id，和单次拉取值
			if err != nil {
				json = fmt.Sprintf(`{"msg":"%s","code":"%d"}`, "团队资源设置失败", -1)
				return
			} else {
				json = fmt.Sprintf(`{"msg":"%s","code":"%d"}`, "团队资源设置成功", 1)
			}
		} else {
			//根据个人用户id取得所属团队，
			user := m.GetUserById(id)
			if onceresource >= user.OnceLimit {
				json = fmt.Sprintf(`{"msg":"%s","code":"%d"}`, "个人单次拉取资源不能超过所属团队", -1)
				this.Ctx.WriteString(json)
				return
			}
			if onedaytotal >= user.OneDayLimit {
				json = fmt.Sprintf(`{"msg":"%s","code":"%d"}`, "个人每日资源不能超过所属团队", -1)
				this.Ctx.WriteString(json)
				return
			}
			err := m.ModifyUserResource(id, onedaytotal, onceresource)
			if err != nil {
				json = fmt.Sprintf(`{"msg":"%s","code":"%d"}`, "个人资源设置失败", -1)
				return
			} else {
				json = fmt.Sprintf(`{"msg":"%s","code":"%d"}`, "个人资源设置成功", 1)
				m.ModifyData(id, onedaytotal)
			}
		}
	}
	this.Ctx.WriteString(json)
}

//团队资源设置列表
func (this *ResourceController) GroupResourceList() {
	total := m.GetNotDistribut()
	resourcetypelist := m.GetAllType()
	shutdownnumber := m.GetShutDownTeam() //得到处于关闭状态的团队数量
	opennumber := m.GetOpenTeam()         //得到处于开启状态的团队数量
	allteam := shutdownnumber + opennumber

	if this.IsAjax() {
		sEcho := this.GetString("sEcho")
		iDisplayStart := this.GetString("iDisplayStart")
		iDisplayLength := this.GetString("iDisplayLength")
		searchStatus, _ := this.GetInt64("sSearch_0")  //状态
		searchAdminName := this.GetString("sSearch_1") //管理者名称
		searchTeamName := this.GetString("sSearch_2")  //团队名称

		var userids []int64
		user := m.GetUserByNames(searchAdminName) //模糊查询得出管理者ID集合
		for i := 0; i < len(user); i++ {
			ids, _ := user[i]["Id"].(int64)
			userids = append(userids, ids)
		}
		filter := make(map[string]interface{})
		filter["s0"] = searchStatus
		filter["s2"] = searchTeamName

		iStart, _ := strconv.ParseInt(iDisplayStart, 10, 64)
		iLength, _ := strconv.ParseInt(iDisplayLength, 10, 64)
		list, count := m.GetAllGroups(iStart, iLength, "Id", filter, userids)
		//已下载数量
		for _, val := range list {
			groupid, _ := val["Id"].(int64)
			count := m.GetGroupPerson(groupid) //团队人数
			val["GroupPerson"] = count
			alreadydown, _ := m.GetGroupTotalDown(groupid)
			val["AlreadyDown"] = alreadydown
			uid, _ := val["Uid"].(int64)
			user := m.GetUserById(uid)
			val["Uid"] = user.Username
		}
		data := make(map[string]interface{})
		data["aaData"] = &list
		data["iTotalDisplayRecords"] = count
		data["iTotalRecords"] = iLength
		data["sEcho"] = sEcho
		this.Data["json"] = &data
		this.ServeJSON()
	} else {
		this.Data["allTeam"] = allteam
		this.Data["openTeam"] = opennumber
		this.Data["shutdownTeam"] = shutdownnumber
		this.Data["ResourceTypeList"] = resourcetypelist
		this.Data["total"] = total
		this.CommonController.CommonMenu()
		this.TplName = "resource/group.html"
	}
}

//个人资源设置列表
func (this *ResourceController) UserResourceList() {
	total := m.GetNotDistribut()
	resourcetypelist := m.GetAllType()
	normal := m.GetNormal()  //正常
	isnormal := m.GetIsNot() //禁用
	Total := normal + isnormal

	userinfo := this.GetSession("userinfo").(m.User)

	if this.IsAjax() {
		sEcho := this.GetString("sEcho")
		iDisplayStart := this.GetString("iDisplayStart")
		iDisplayLength := this.GetString("iDisplayLength")
		iStart, _ := strconv.ParseInt(iDisplayStart, 10, 64)
		iLength, _ := strconv.ParseInt(iDisplayLength, 10, 64)
		searchStatus, _ := this.GetInt64("sSearch_0") // 状态
		searchTeam, _ := this.GetInt64("sSearch_1")   //团队id
		searchName := this.GetString("sSearch_2")     //搜索姓名或者昵称
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

		if searchTeam > 0 {
			Cate, _ := m.Groups()
			usergroup = m.GetTreeGroup(Cate, searchTeam)
			usergroup = append(usergroup, searchTeam)
		}

		filter := make(map[string]interface{})
		filter["s0"] = searchStatus
		filter["s1"] = groups
		filter["s2"] = searchName
		filter["s3"] = usergroup

		list, count, _ := m.GetUserResourceLists(iStart, iLength, "Id", filter)

		for _, val := range list {
			groupid, _ := val["Group"].(int64)
			userid, _ := val["Id"].(int64)
			group, _ := m.GetGroupById(groupid)
			val["Group"] = group.Name
			alreadydown, _ := m.GetGroupTotalDown(userid)
			val["AlreadyDown"] = alreadydown
		}

		data := make(map[string]interface{})
		data["aaData"] = &list
		data["iTotalDisplayRecords"] = count
		data["iTotalRecords"] = iLength
		data["sEcho"] = sEcho
		this.Data["json"] = &data
		this.ServeJSON()
	} else {
		this.Data["allUser"] = Total
		this.Data["normal"] = normal
		this.Data["notnormal"] = isnormal
		this.Data["ResourceTypeList"] = resourcetypelist
		this.Data["total"] = total

		this.CommonController.CommonMenu()
		this.TplName = "resource/group.html"
	}
}
