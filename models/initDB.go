package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"database/sql"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

/*
* 初始化数据库
*包括穿件数据库，表，以及插入部分数据
 */
func InitDB() {
	createdb()
	Connect()
	inserData()
	// insertUser()
	// insertGroup()
	// insertRole()
	// insertNode()
	// fmt.Println("database init is complete.\nPlease restart the application")
}

/**
* 创建数据库
 */
func createdb() {
	var sqlstring string
	dns, db_name := getConfig(0)
	sqlstring = fmt.Sprintf("CREATE DATABASE  if not exists `%s` CHARSET utf8 COLLATE utf8_general_ci", db_name)
	db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err.Error())
	}
	r, err := db.Exec(sqlstring)
	// ExecSqlFile(db)
	if err != nil {
		beego.Error("err is :", err.Error(), "and r is:", r)
	} else {
		beego.Info("Database: ", db_name, " created succes")
	}
	defer db.Close()
}
func Connect() {
	dns, _ := getConfig(1)
	beego.Info("数据库is %s", dns)
	err := orm.RegisterDataBase("default", "mysql", dns)
	if err != nil {
		beego.Error("数据库连接失败")
	} else {
		beego.Info("数据库连接sucess ")
		// writeSiteConf()
	}
}

/*
* 获取配置
	flag ==1 表示 只链接
	==0 创建 加链接
*/
func getConfig(flag int) (string, string) {
	var dns string
	db_host := beego.AppConfig.String("db_host")
	db_port := beego.AppConfig.String("db_port")
	db_user := beego.AppConfig.String("db_user")
	db_pass := beego.AppConfig.String("db_pass")
	db_name := beego.AppConfig.String("db_name")
	if flag == 1 {
		// fmt.Println("链接数据库")
		orm.RegisterDriver("mysql", orm.DRMySQL)
		dns = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&loc=Local", db_user, db_pass, db_host, db_port, db_name)
	} else {
		// fmt.Println("创建数据库")
		dns = fmt.Sprintf("%s:%s@tcp(%s:%s)?charset=utf8", db_user, db_pass, db_host, db_port)
	}
	return dns, db_name
}

// Id            int64     // 标识
// Username      string    // 用户名
// Password      string    // 密码
// Nickname      string    // 昵称
// Status        int       `orm:"default(1)"` // 状态 [1: 正常，2：封停，3:冻结]
// Empnum        string    // 工号
// Companyname   string    // 公司名称
// Fid           int64     // 父级用户
// Token         string    // Token
// Remark        string    `orm:"null;type(text)" `
// Role          *Role     `orm:"rel(fk);"`
// Group         *Group    `orm:"rel(fk);"`
// Protected     int64     `orm:"size(1);default(0)"`          // 保护字段 0不保护 1保护 保护状态下不可编辑
// Createtime    time.Time `orm:"type(datetime);auto_now_add"` // 创建时间
// Lastlogintime time.Time `orm:"type(datetime);auto_now_add"` // 最后登录时间
func inserData() {
	beego.Debug("insert group start")
	group := new(Group)
	group.Id = 1
	group.Name = "管理组"
	group.Uid = 1
	group.Status = 1
	group.Remark = "系统初始默认组"
	AddGroup(group)
	beego.Debug("insert group end")
	beego.Debug("instrt user start")
	user := new(User)
	user.Id = 1
	user.Username = "admin"
	user.Password = "21232f297a57a5a743894a0e4a801fc3"
	user.Empnum = "admin"
	user.Companyname = "皓月科技"
	user.Remark = "系统默认超级管理员"
	AddUser(user)
}

func ExecSqlFile(db *sql.DB) {
	sqlpath := "static/sql/hcloud.sql"
	_, err := os.Stat(sqlpath)
	if err != nil {
		beego.Error("stat", err)
	} else {
		filepath, _ := os.Open(sqlpath)
		file, err := ioutil.ReadAll(filepath)
		if err != nil {
			beego.Error("read sql file error", err)
		}
		requests := strings.Split(string(file), ";")
		for _, request := range requests {
			beego.Debug("request", request)
			result, err := db.Exec(request)
			if err != nil {
				beego.Error("exec sql file error", err)
			}
			beego.Debug("result:", result)
		}
	}
}

// func insertUser() {
// 	model := orm.NewOrm()
// 	err := model.Begin()
// 	userextend := new(UserExtend)
// 	userextend.Id = 1
// 	userextend.Uname = 6666
// 	userextend.LoginIp = "127.0.0.1"
// 	uid, err := model.Insert(userextend)
// 	user := new(User)
// 	user.Id = 1
// 	user.Uname = 6666
// 	user.Pwd = EncodeUserPwd("admin")
// 	user.Nickname = "超级管理员"
// 	user.Email = "admin@ihaoyue.com"
// 	user.Group = &Group{Id: 1}
// 	user.UserExtend = &UserExtend{Id: uid}
// 	_, err = model.Insert(user)
// 	if err != nil {
// 		beego.Error(err)
// 		// fmt.Printf("num is %d and err is %s", num, err.Error())
// 	} else {
// 		// fmt.Println("insert user end")
// 	}
// }
// func insertGroup() {
// 	group := new(Group)
// 	group.Id = 1
// 	group.Name = "管理员组"
// 	group.Sort = 50
// 	group.Market = 0
// 	group.Fid = 0
// 	group.Remark = "管理员组"
// 	_, err := AddGroup(group)
// 	if err != nil {
// 		// fmt.Println("添加资源失败")
// 		beego.Error(err)
// 		return
// 	}
// }
/*func insertNodes() {
	fmt.Println("insert node ...")
	g := new(Group)
	g.Id = 1
	//nodes := make([20]Node)
	nodes := [24]Node{
		{Name: "rbac", Title: "RBAC", Remark: "", Level: 1, Pid: 0, Status: 2, Group: g},
		{Name: "node/index", Title: "Node", Remark: "", Level: 2, Pid: 1, Status: 2, Group: g},
		{Name: "index", Title: "node list", Remark: "", Level: 3, Pid: 2, Status: 2, Group: g},
		{Name: "AddAndEdit", Title: "add or edit", Remark: "", Level: 3, Pid: 2, Status: 2, Group: g},
		{Name: "DelNode", Title: "del node", Remark: "", Level: 3, Pid: 2, Status: 2, Group: g},
		{Name: "user/index", Title: "User", Remark: "", Level: 2, Pid: 1, Status: 2, Group: g},
		{Name: "Index", Title: "user list", Remark: "", Level: 3, Pid: 6, Status: 2, Group: g},
		{Name: "AddUser", Title: "add user", Remark: "", Level: 3, Pid: 6, Status: 2, Group: g},
		{Name: "UpdateUser", Title: "update user", Remark: "", Level: 3, Pid: 6, Status: 2, Group: g},
		{Name: "DelUser", Title: "del user", Remark: "", Level: 3, Pid: 6, Status: 2, Group: g},
		{Name: "group/index", Title: "Group", Remark: "", Level: 2, Pid: 1, Status: 2, Group: g},
		{Name: "index", Title: "group list", Remark: "", Level: 3, Pid: 11, Status: 2, Group: g},
		{Name: "AddGroup", Title: "add group", Remark: "", Level: 3, Pid: 11, Status: 2, Group: g},
		{Name: "UpdateGroup", Title: "update group", Remark: "", Level: 3, Pid: 11, Status: 2, Group: g},
		{Name: "DelGroup", Title: "del group", Remark: "", Level: 3, Pid: 11, Status: 2, Group: g},
		{Name: "role/index", Title: "Role", Remark: "", Level: 2, Pid: 1, Status: 2, Group: g},
		{Name: "index", Title: "role list", Remark: "", Level: 3, Pid: 16, Status: 2, Group: g},
		{Name: "AddAndEdit", Title: "add or edit", Remark: "", Level: 3, Pid: 16, Status: 2, Group: g},
		{Name: "DelRole", Title: "del role", Remark: "", Level: 3, Pid: 16, Status: 2, Group: g},
		{Name: "Getlist", Title: "get roles", Remark: "", Level: 3, Pid: 16, Status: 2, Group: g},
		{Name: "AccessToNode", Title: "show access", Remark: "", Level: 3, Pid: 16, Status: 2, Group: g},
		{Name: "AddAccess", Title: "add accsee", Remark: "", Level: 3, Pid: 16, Status: 2, Group: g},
		{Name: "RoleToUserList", Title: "show role to userlist", Remark: "", Level: 3, Pid: 16, Status: 2, Group: g},
		{Name: "AddRoleToUser", Title: "add role to user", Remark: "", Level: 3, Pid: 16, Status: 2, Group: g},
	}
	for _, v := range nodes {
		n := new(Node)
		n.Name = v.Name
		n.Title = v.Title
		n.Remark = v.Remark
		n.Level = v.Level
		n.Pid = v.Pid
		n.Status = v.Status
		n.Group = v.Group
		o.Insert(n)
	}
	fmt.Println("insert node end")
}
*/
