package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	"hcloud/controllers"
	m "hcloud/models"
	_ "hcloud/routers"
	"strings"
)

//验证登录过滤器
var Filter = func(ctx *context.Context) {
	if b, _ := beego.AppConfig.Bool("EnableHttpTLS"); b {
		httpscheme := ctx.Input.Scheme()     //判断http或者https
		httpmethod := ctx.Request.RequestURI //获取控制器加方法
		if httpscheme == "http" && httpmethod != "/api" {
			httphost := ctx.Request.Host //获取主机
			httphostarray := strings.Split(httphost, ":")
			httpurl := "https://" + httphostarray[0] + ":" + beego.AppConfig.String("HttpsPort") + httpmethod //组合URL
			ctx.Redirect(302, httpurl)
		}
	}
	k := ctx.Input.CruSession.Get("userinfo")
	url := ctx.Request.RequestURI
	if url != "/hcloud/login" && !strings.HasPrefix(url, "/hcloud/changepwd") && k == nil {
		ctx.Redirect(302, "/hcloud/login")
	}
}

//权限过滤器
var AccessFilter = func(ctx *context.Context) {
	userinfo := ctx.Input.CruSession.Get("userinfo")
	if userinfo != nil {
		adminUser := beego.AppConfig.String("admin_user")
		//root直接通过认证
		if userinfo.(m.User).Username == adminUser {
			return
		}
		url := ctx.Request.RequestURI
		temp := url[1:]
		res := m.GetResByUrl(temp)
		//只判断已添加的节点
		if res.Id == 0 {
			return
		} else {
			//检测权限
			flag := m.CheckAccessRole(res.Id, userinfo.(m.User).Role.Id)
			if !flag {
				ctx.WriteString("")
			}
		}
	}
}

func main() {
	m.InitDB()
	orm.RunSyncdb("default", false, true)
	//添加过滤器
	beego.InsertFilter("/hcloud/*", beego.BeforeRouter, Filter)
	beego.InsertFilter("/hcloud/*", beego.BeforeRouter, AccessFilter)
	orm.Debug, _ = beego.AppConfig.Bool("ormdebug")
	beego.ErrorController(&controllers.ErrorController{})       //注册错误处理的函数
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600 //设置Session过期时间

	// 设置静态目录路径
	beego.SetStaticPath("/upload", "../upload")
	beego.Run()
}
