package models

import (
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

/*
*用户表
 */
type User struct {
	Id                 int64                // 标识
	Username           string               // 用户名
	Password           string               // 密码
	Nickname           string               // 昵称
	Status             int                  `orm:"default(1)"` // 状态 [1: 正常，2：禁用]
	Empnum             string               // 工号
	Companyname        string               // 公司名称
	Fid                int64                // 父级用户
	Token              string               // Token
	Remark             string               `orm:"null;type(text)" `
	Role               *Role                `orm:"rel(fk);"`
	Group              *Group               `orm:"rel(fk);"`
	Protected          int64                `orm:"size(1);default(0)"`          // 保护字段 0不保护 1保护 保护状态下不可编辑
	Createtime         time.Time            `orm:"type(datetime);auto_now_add"` // 创建时间
	Lastlogintime      time.Time            `orm:"type(datetime);auto_now_add"` // 最后登录时间
	CloudResourceTypes []*CloudResourceType `orm:"reverse(many)"`
	Imgurl             string               // 用户头像名称

	OneDayLimit int64 `orm:"default(30)"` //每日资源上限
	OnceLimit   int64 `orm:default(1)`    //单次拉取上限

	Officephone string // 办公电话
	Familyphone string // 家庭电话
	Telphone    string // 手机号
	Email       string // 电子邮箱
	Wxorqq      string // 微信或者QQ
	Zipcode     string // 邮政编码
	Address     string // 联系地址

	Createtimestr    string `orm:"-"`
	Lastlogintimestr string `orm:"-"`
	Createname       string `orm:"-"` // 创建者
	Deviceconut      int    `orm:"-"` // 设备个数
	Teamname         string `orm:"-"` // 团队名称
	Positionname     string `orm:"-"` // 职位名称
	Statusname       string `orm:"-"` // 状态名称
}

func (u *User) TableName() string {
	return "user"
}

func init() {
	orm.RegisterModel(new(User))
}

//根据团队id取得所有该团队下面的用户
func GetAllTeamUser(id int64) (list []orm.Params) {
	o := orm.NewOrm()
	_, err := o.QueryTable("user").Filter("Group__Id", id).Values(&list)
	if err != nil {
		beego.Error("get team user is err", err)
	}
	return list
}

func GetUserList(page int64, page_size int64, groups []int64, uname string, userid int64, username string, sState int, usergroup []int64, sUserrole int64, companyname string) (users []orm.Params, count int64) {
	omodel := orm.NewOrm()
	user := new(User)
	qs := omodel.QueryTable(user)
	//登录用户判断身份

	var cond *orm.Condition
	cond = orm.NewCondition()

	admin_user := beego.AppConfig.String("admin_user")
	if admin_user == username {
		cond = cond.And("Id__gt", 1)
	} else {
		if len(groups) > 0 {
			cond = cond.And("group_id__in", groups)
		} else {
			cond = cond.And("Id", userid)
		}
	}
	if len(uname) > 0 {
		cond = cond.And("Username__icontains", uname)
	}
	if len(companyname) > 0 {
		cond = cond.And("Companyname__icontains", companyname)
	}
	if sState > 0 {
		cond = cond.And("status", sState)
	}
	if sUserrole > 0 {
		cond = cond.And("Role__Id", sUserrole)
	}
	if len(usergroup) > 0 {
		cond = cond.And("group_id__in", usergroup)
	}

	qs.Limit(page_size, page).OrderBy("Username").SetCond(cond).Values(&users)

	count, _ = qs.SetCond(cond).Count()
	return users, count
}

// 根据Group查询用户
func GetUsersByGroup(groups []int64) (list []orm.Params, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("user").Filter("group_id__in", groups).Values(&list)
	return list, err
}

// 根据fid 查询用户
func GetUsersByFid(page int64, page_size int64, fid int64) (list []orm.Params, count int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable("user").Filter("Fid", fid)
	qs.Limit(page_size, page).Values(&list)
	count, _ = qs.Count()
	return list, count, err

}

/*
*获取所有用户列表
 */
func GetUserListAll(groups []int64, userid int64, username int64, readalluser int64) (users []orm.Params, count int64) {
	omodel := orm.NewOrm()
	user := new(User)
	qs := omodel.QueryTable(user)

	var cond *orm.Condition
	cond = orm.NewCondition()

	if readalluser > 0 || beego.AppConfig.DefaultInt64("admin_user", 0) == username {
		// cond = cond.And("Id__gt", 0)
	} else if len(groups) > 0 {
		cond = cond.And("group_id__in", groups)
	} else {
		cond = cond.And("Id", userid)
	}

	qs.SetCond(cond).Values(&users)

	count, _ = qs.SetCond(cond).Count()
	return users, count
}

func userQuerySeter() orm.QuerySeter {
	return orm.NewOrm().QueryTable(new(User))
}

// 获取所有用户
func UserCate(groupids []int64) (users []*User) {
	cond := orm.NewCondition()
	if len(groupids) > 0 {
		cond = cond.And("group_id__in", groupids)
	}
	userQuerySeter().SetCond(cond).All(&users)
	return
}

// 无限极分类
func GetUserTree(cate []*User, fid int64) []*User {
	var tree []*User = make([]*User, 0)

	for _, v := range cate {
		if v.Fid == fid {
			child := GetUserTree(cate, v.Id)
			tree = append(tree, v)
			tree = append(tree, child...)
		}
	}
	return tree
}

/*
*根据条件获取所有用户列表
 */
func GetUserListByCond(groups []int64, userid int64) (users []orm.Params, count int64) {
	omodel := orm.NewOrm()
	user := new(User)
	qs := omodel.QueryTable(user)

	var cond *orm.Condition
	cond = orm.NewCondition()

	if len(groups) > 0 {
		cond = cond.And("group_id__in", groups)
	}
	if userid > 0 {
		cond = cond.And("Id", userid)
	}

	qs.SetCond(cond).Values(&users)

	count, _ = qs.SetCond(cond).Count()
	return users, count
}

func GetReadUser(Id int64) (error, User) {
	o := orm.NewOrm()
	user := User{Id: Id}
	err := o.Read(&user, "Id")
	return err, user
}

// 添加用户
func AddUser(u *User) (int64, error) {
	model := orm.NewOrm()
	id, err := model.Insert(u)
	return id, err
}

// 修改用户
func UpdateUser(u *User) (int64, error) {
	model := orm.NewOrm()
	id, err := model.Update(u)
	return id, err
}

// 批量删除用户
func DelUser(ids []string) error {
	omodel := orm.NewOrm()
	err := omodel.Begin()
	if len(ids) > 0 {
		for i := 0; i < len(ids); i++ {
			id, _ := strconv.ParseInt(ids[i], 10, 64)

			_, err = omodel.Delete(&User{Id: id})
		}
		if err != nil {
			omodel.Rollback()
		} else {
			omodel.Commit()
		}
	} else {
		err = errors.New("the userid is not exist.")
	}
	return err
}

// 删除用户
func DelUserById(id int64) (int64, error) {
	omodel := orm.NewOrm()
	num, err := omodel.Delete(&User{Id: id})
	return num, err
}

//根据ID获取user
func GetRoleByUserId(userId int64) (roles Role) {
	o := orm.NewOrm()
	role := new(Role)
	err := o.QueryTable(role).Filter("User__Id", userId).One(&roles)
	if err != nil {
		beego.Error(err)
	}
	return roles
}

//根据Uname获取user
func GetUserByUname(username string) (user User) {
	o := orm.NewOrm()
	user = User{Username: username}
	err := o.Read(&user, "Username")
	if err == orm.ErrNoRows {
		beego.Error("not find the username", username)
	} else if err == orm.ErrMissPK {
		beego.Error("not find the key")
	} else {
		beego.Error(err)
	}
	return user
}

//模糊根据name获取userid
func GetUserByNames(username string) (list []orm.Params) {
	o := orm.NewOrm()
	_, err := o.QueryTable("user").Filter("Username__icontains", username).Values(&list)
	beego.Debug(err)
	return list
}

// 根据id获取用户
func GetUserById(id int64) (user User) {
	o := orm.NewOrm()
	user = User{Id: id}
	err := o.Read(&user, "Id")
	if err == orm.ErrNoRows {
		beego.Error("not find the username id:")
	} else if err == orm.ErrMissPK {
		beego.Error("not find the username key")
	}
	return user
}

func Users(groups []int64, userid int64, username string, readalluser int64, status []int) (cout int64, users []orm.Params) {
	o := orm.NewOrm()
	user := new(User)

	// 用户限定
	var cond *orm.Condition
	cond = orm.NewCondition()

	if readalluser > 0 || beego.AppConfig.String("admin_user") == username {
		cond = cond.And("Id__gt", 2)
	} else if len(groups) > 0 {
		cond = cond.And("group_id__in", groups)
		cond = cond.And("Id__gt", 2)
	} else {
		cond = cond.And("Id", userid)
		cond = cond.And("Id__gt", 2)
	}
	cond = cond.And("status__in", status)
	qs := o.QueryTable(user)
	count, _ := qs.SetCond(cond).Values(&users, "id", "user_name", "nickname", "remark", "password")

	return count, users
}

// 获取在职员工
func GetInUsers(groups []int64, userid int64, username string, readalluser int64) (cout int64, users []orm.Params) {
	o := orm.NewOrm()
	user := new(User)

	// 用户限定
	var cond *orm.Condition
	cond = orm.NewCondition()

	if readalluser > 0 || beego.AppConfig.String("admin_user") == username {
		cond = cond.And("Id__gt", 2)
	} else if len(groups) > 0 {
		cond = cond.And("group_id__in", groups)
		cond = cond.And("Id__gt", 2)
	} else {
		cond = cond.And("Id", userid)
		cond = cond.And("Id__gt", 2)
	}
	cond = cond.And("Status", 2)
	qs := o.QueryTable(user)
	count, _ := qs.SetCond(cond).OrderBy("Username").Values(&users, "id", "user_name", "nickname", "remark", "password", "group_id")

	return count, users
}

// 用户登录校验并返回用户表的信息和设备信息
func LoginUserCheck(username string, password string, macaddress string, times string, clientip string) (*User, bool) {
	o := orm.NewOrm()
	user := User{Username: username, Password: password}

	if err := o.Read(&user, "Username", "Password"); err != nil {
		beego.Error(err)
		return nil, false
	}

	user.Lastlogintime = time.Now()
	_, err := o.Update(&user)
	if err != nil {
		beego.Error(err)
	}

	return &user, true
}

//检查用户是否存在
func CheckUserInfo(username string, password string) (int64, *User) {
	model := orm.NewOrm()
	if len(username) > 0 && len(password) == 0 {
		count, err := model.QueryTable("user").Filter("Username", username).Count()
		return count, nil
		beego.Error(err)
	} else if len(username) > 0 && len(password) == 32 {
		var user User
		if err := model.QueryTable("user").Filter("Username", username).Filter("Password", password).One(&user); err != nil {
			beego.Error(err)
			return 1, &user
		}
	}
	return 0, nil
}

// 查询条件的员工数量
func GetCountNum(groups []int64) int64 {
	model := orm.NewOrm()
	count, _ := model.QueryTable("user").Filter("status", 2).Filter("group_id__in", groups).Count()
	return count
}

// 禁用/启用账户
func EditUserStatus(ids []string, status int64) error {
	omodel := orm.NewOrm()
	err := omodel.Begin()
	var table User
	user := make(orm.Params)
	user["Status"] = status
	if len(ids) > 0 {
		for i := 0; i < len(ids); i++ {
			id, _ := strconv.ParseInt(ids[i], 10, 64)
			_, err = omodel.QueryTable(table).Filter("Id", id).Update(user)
		}
		if err != nil {
			omodel.Rollback()
		} else {
			omodel.Commit()
		}
	} else {
		err = errors.New("the userid is not exist.")
	}
	return err
}

//更新用户登录时间
func UpdateLoginTime(id int64) (int64, error) {
	time := time.Now().Format("2006-01-02 15:04:05")
	omodel := orm.NewOrm()
	user := make(orm.Params)
	user["Lastlogintime"] = time

	var table User
	num, err := omodel.QueryTable(table).Filter("Id", id).Update(user)
	return num, err
}

//上传图片
func UploadImg(url string, id int64) (int64, error) {
	omodel := orm.NewOrm()
	user := make(orm.Params)
	user["Imgurl"] = url

	var table User
	num, err := omodel.QueryTable(table).Filter("Id", id).Update(user)
	return num, err
}

//获取所有的在职员工
func GetAllActiveUsers() (list []orm.Params) {
	o := orm.NewOrm()
	/*sql := "select id,uname from user where status=2 and id>1"
	o.Raw(sql).Values(&list)*/
	o.QueryTable("user").Filter("status", 2).Filter("id__gt", 2).Values(&list, "Id", "Username", "Nickname")
	return list
}

//获取所有的在职业务员 业务员角色id=11
func GetAllActiveSaleman() (list []orm.Params) {
	o := orm.NewOrm()
	sql := "select id,username from user where status=2 and id>2 and id in (select user_id from user_roles where role_id=11)"
	o.Raw(sql).Values(&list)
	return list
}

func GetUserByUnameOrNick(name string) (user User) {
	o := orm.NewOrm()
	cond := orm.NewCondition()
	cond = cond.And("Username", name).Or("Nickname", name)
	o.QueryTable("user").SetCond(cond).Limit(1).One(&user)
	return user
}

//删除组别下的所有用户
func DeleteUserByGroup(id int64) {
	o := orm.NewOrm()
	var users []orm.Params
	o.QueryTable("user").Filter("group_id", id).Values(&users)
	for _, item := range users {
		DelUserById(item["Id"].(int64))
	}
}

//获取用户
func ReadUserByField(e *User, fields ...string) (*User, error) {
	o := orm.NewOrm()
	if err := o.Read(e, fields...); err != nil {
		return e, err
	}
	return nil, nil
}

//根据组别获取在职用户
func GetActiveByGroup(group string) (list []orm.Params) {
	o := orm.NewOrm()
	arr := strings.Split(group, ",")
	o.QueryTable("user").Filter("status", 2).Filter("group__id__in", arr).Filter("id__gt", 2).Filter("Role__role_id", 11).Values(&list, "Id", "Username", "Nickname")
	return list
}

//修改工号
func ChangeUname(old_uname int64, new_uname string, new_group int64, new_name string) int64 {
	tempUser := GetUserByUname(new_uname)
	if tempUser.Id > 0 {
		//工号已存在
		return -1
	}
	o := orm.NewOrm()
	user := make(orm.Params)
	user["Username"] = new_uname
	user["Group"] = new_group
	user["Nickname"] = new_name
	num, err := o.QueryTable("user").Filter("Username", old_uname).Update(user)
	if err == nil {
		return num
	} else {
		//更新失败
		return 0
	}
}

//获取用户权限id
func GetUserRoleID(id int64) (roleId string) {
	o := orm.NewOrm()
	sql := "select role_id from user_roles where user_id=?"
	var temp []orm.ParamsList
	o.Raw(sql, id).ValuesList(&temp)
	roleId = temp[0][0].(string)
	return roleId
}

// 根据姓名查找
func GetUserListByNickname(nickname string) (list []orm.Params, count int64, err error) {
	cond := orm.NewCondition()
	cond = cond.And("NickName__icontains", nickname).Or("Username__icontains", nickname)
	count, err = userQuerySeter().SetCond(cond).OrderBy("id").Values(&list)
	return list, count, err
}

// 根据Fid查询所有用户
func GetUserByFid(groups []int64, userinfo *User) (users []orm.Params, count int64) {
	omodel := orm.NewOrm()
	user := new(User)
	qs := omodel.QueryTable(user)

	var cond *orm.Condition
	cond = orm.NewCondition()

	if userinfo.Username == beego.AppConfig.String("admin_user") {
		cond = cond.And("Id__gt", 1)
	} else if len(groups) > 0 {
		cond = cond.And("group_id__in", groups)
	} else {
		cond = cond.And("Id", userinfo.Id)
	}

	qs.OrderBy("-Id").SetCond(cond).Values(&users)
	count, _ = qs.SetCond(cond).Count()
	return users, count
}

// 根据Username 、token 查询
func GetUserByToken(username string, token string) (user User, err error) {
	o := orm.NewOrm()
	err = o.QueryTable("user").Filter("Username", username).Filter("Token", token).One(&user)
	return user, err
}

func GetUsers() (users []*User, err error) {
	o := orm.NewOrm()
	_, err = o.QueryTable("user").All(&users)
	return users, err
}

//根据id修改个人资源限制
func ModifyUserResource(id, total, once int64) error {
	o := orm.NewOrm()
	_, err := o.QueryTable("user").Filter("Id", id).Update(orm.Params{
		"OneDayLimit": total,
		"OnceLimit":   once,
	})
	return err
}

//团队人数
func GetGroupPerson(groupid int64) int64 {
	o := orm.NewOrm()
	num, _ := o.QueryTable("user").Filter("Group__Id", groupid).Count()
	return num
}

//个人列表
func GetUserResourceList(page int64, page_size int64, sort string, filter map[string]interface{}, grouids int64) (groups []orm.Params, count int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable("user")

	var cond *orm.Condition
	cond = orm.NewCondition()
	se0, _ := filter["s0"].(int64)
	se1, _ := filter["s1"].(int64)
	se2 := filter["s2"].(string)
	if grouids > 0 {
		cond = cond.And("Group__Id", grouids)
	}
	if se0 > 0 {
		cond = cond.And("Status", se0)
	}
	if se1 != 0 {
		cond = cond.And("Group", se1)
	}
	if len(se2) > 0 {
		cond = cond.And("Username__icontains", se2)
	}
	qs.SetCond(cond).OrderBy(sort).Limit(page_size, page).Values(&groups)
	count, err = qs.SetCond(cond).Count()
	return groups, count, err
}

//管理员个人列表
func GetUserResourceLists(page int64, page_size int64, sort string, filter map[string]interface{}) (groups []orm.Params, count int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable("user")

	var cond *orm.Condition
	cond = orm.NewCondition()
	if filter != nil {
		se0, _ := filter["s0"].(int64)
		se1, _ := filter["s1"].([]int64)
		se2 := filter["s2"].(string)
		se3 := filter["s3"].([]int64)
		if se0 > 0 {
			cond = cond.And("Status", se0)
		}
		if len(se1) > 0 {
			cond = cond.And("group_id__in", se1)
		}
		if len(se2) > 0 {
			cond = cond.And("Username__icontains", se2).Or("Nickname__icontains", se2)
		}
		if len(se3) > 0 {
			cond = cond.And("group_id__in", se3)
		}
	}
	qs.SetCond(cond).OrderBy(sort).Limit(page_size, page).Values(&groups)
	count, err = qs.SetCond(cond).Count()
	return groups, count, err
}

//正常
func GetNormal() int64 {
	o := orm.NewOrm()
	num, _ := o.QueryTable("user").Filter("Status", 1).Count()
	return num
}

//禁用
func GetIsNot() int64 {
	o := orm.NewOrm()
	num, _ := o.QueryTable("user").Filter("Status", 2).Count()
	return num
}

// 获取角色为群控账户的用户
func GetCouldRoleUserByGroup(groupid []int64) (user []*User, err error) {
	o := orm.NewOrm()
	couldrole, _ := beego.AppConfig.Int64("couldrole")
	_, err = o.QueryTable("user").Filter("Group__Id__in", groupid).Filter("Role__Id", couldrole).All(&user)
	return
}

//团队资源设置的时候修改单次拉取值
//团队id，数量
func ModifyOncePerson(groupid, num int64) {
	o := orm.NewOrm()
	//根据id取得该组别下所有的用户
	list := GetAllUserTeam(groupid)
	for i := 0; i < len(list); i++ {
		userid := list[i]["Id"]
		number, _ := o.QueryTable("user").Filter("Id", userid).Update(orm.Params{
			"OnceLimit": num,
		})
		beego.Debug(number)
	}

}

//根据团队id取得所有的成员
func GetAllUserTeam(groupid int64) (list []orm.Params) {
	o := orm.NewOrm()
	_, err := o.QueryTable("user").Filter("Group__Id", groupid).Values(&list)
	if err != nil {
		beego.Error("get this team user is err", err)
	}
	return list
}
