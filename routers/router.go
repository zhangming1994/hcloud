package routers

import (
	"github.com/astaxie/beego"
	"hcloud/controllers"
	"hcloud/controllers/api"
	"hcloud/models"
	"mime"
	"os"
)

func init() {
	initialize()
	router()
	if ok, _ := beego.AppConfig.Bool("testdebug"); ok {
		testrouter()
	}
}

func initialize() {
	mime.AddExtensionType(".css", "text/css")
	args := os.Args
	for _, v := range args {
		if v == "InitDB" {
			// 创建表 插入数据等等，只在执行./hcloud InitDB 时候执行一次
			models.InitDB()
			os.Exit(0)
		}
	}
	// 链接数据库
	// models.Connect()
}
func router() {
	beego.Router("/", &controllers.MainController{}, "*:Index")
	beego.Router("/hcloud/index", &controllers.MainController{}, "*:Index")
	beego.Router("/hcloud/login", &controllers.MainController{}, "*:Login")
	beego.Router("/hcloud/logout", &controllers.MainController{}, "*:Logout")
	beego.Router("/hcloud/main", &controllers.MainController{}, "*:MainFrame")
	beego.Router("/hcloud/about", &controllers.MainController{}, "*:About")

	beego.Router("/hcloud/background/top", &controllers.MainController{}, "*:Top")
	beego.Router("/hcloud/background/center", &controllers.MainController{}, "*:Center")
	beego.Router("/hcloud/background/left", &controllers.MainController{}, "*:Left")
	beego.Router("/hcloud/background/tab", &controllers.MainController{}, "*:Tab")

	beego.Router("/hcloud/user/list", &controllers.UserController{}, "*:List")
	beego.Router("/hcloud/user/edit", &controllers.UserController{}, "*:Edit")
	beego.Router("/hcloud/user/delete", &controllers.UserController{}, "*:Delete")
	beego.Router("/hcloud/user/add", &controllers.UserController{}, "*:Add")
	beego.Router("/hcloud/user/editpwd", &controllers.UserController{}, "*:EditPassword")
	beego.Router("/hcloud/user/resetpwd", &controllers.UserController{}, "*:ResetPassword")
	beego.Router("/hcloud/user/view", &controllers.UserController{}, "*:ViewUser")
	beego.Router("/hcloud/user/person", &controllers.UserController{}, "*:GetPerson")

	beego.Router("/hcloud/user/allocation", &controllers.UserController{}, "*:Allocation")
	beego.Router("/hcloud/user/edituserstatus", &controllers.UserController{}, "*:EditUserStatus")
	beego.Router("/hcloud/user/uploadImg", &controllers.UserController{}, "*:UploadImg")
	beego.Router("/hcloud/user/export", &controllers.UserController{}, "*:ExportUserList")

	beego.Router("/hcloud/role/list", &controllers.RoleController{}, "*:List")
	beego.Router("/hcloud/role/add", &controllers.RoleController{}, "*:Add")
	beego.Router("/hcloud/role/delete", &controllers.RoleController{}, "*:Delete")
	beego.Router("/hcloud/role/edit", &controllers.RoleController{}, "*:Edit")
	beego.Router("/hcloud/role/resource", &controllers.RoleController{}, "*:AllocationRes")
	beego.Router("/hcloud/role/allotnode", &controllers.RoleController{}, "*:AllotNode")

	beego.Router("/hcloud/group/list", &controllers.GroupController{}, "*:List")
	beego.Router("/hcloud/group/add", &controllers.GroupController{}, "*:Add")
	beego.Router("/hcloud/group/edit", &controllers.GroupController{}, "*:Edit")
	beego.Router("/hcloud/group/delete", &controllers.GroupController{}, "*:Delete")

	beego.Router("/hcloud/device/list", &controllers.DeviceController{}, "*:List")
	beego.Router("/hcloud/device/allot", &controllers.DeviceController{}, "*:AllotDevice")
	beego.Router("/hcloud/device/allothandle", &controllers.DeviceController{}, "*:AllotDeviceHandle")
	beego.Router("/hcloud/device/getdevices", &controllers.DeviceController{}, "*:GetDevices")
	//resource

	beego.Router("/hcloud/resource/list", &controllers.ResourceController{}, "*:Index")                            //资源首页
	beego.Router("/hcloud/resource/uploadindex", &controllers.ResourceController{}, "*:UploadIndex")               //文件上传页面
	beego.Router("/hcloud/resource/groupset", &controllers.ResourceController{}, "*:GroupResourceList")            //资源分发页面
	beego.Router("/hcloud/resource/peronset", &controllers.ResourceController{}, "*:UserResourceList")             //个人资源设置
	beego.Router("/hcloud/resource/upload", &controllers.ResourceController{}, "*:UploadPhoneResource")            //文件上传
	beego.Router("/hcloud/resource/typedeta", &controllers.ResourceController{}, "*:ListDetail")                   //折叠详细
	beego.Router("/hcloud/resource/addtype", &controllers.ResourceController{}, "*:AddResourceType")               //增加资源类型
	beego.Router("/hcloud/resource/edittype", &controllers.ResourceController{}, "*:EditResourceType")             //修改资源类型
	beego.Router("/hcloud/resource/edittypestatus", &controllers.ResourceController{}, "*:EditResourceTypeStatus") //修改资源类型状态
	beego.Router("/hcloud/resource/deltype", &controllers.ResourceController{}, "*:DelResourceType")               //删除资源类型
	beego.Router("/hcloud/resource/gettypeusers", &controllers.ResourceController{}, "*:GetResourceTypeUsers")     //获取筛选资源类型参数
	beego.Router("/hcloud/resource/typelist", &controllers.ResourceController{}, "*:ListResourceType")             //资源类型列表
	beego.Router("/hcloud/resource/getrestype", &controllers.ResourceController{}, "*:GetResourceType")            //获取资源类型
	beego.Router("/hcloud/resource/grouplist", &controllers.ResourceController{}, "*:GetListGroup")                //加载团队列表
	beego.Router("/hcloud/resource/userlist", &controllers.ResourceController{}, "*:GetListUser")                  //加载用户列表
	beego.Router("/hcloud/resource/getgroupuser", &controllers.ResourceController{}, "*:GetUserByGroup")           //根据团队查询用户列表
	beego.Router("/hcloud/resource/addtypeuser", &controllers.ResourceController{}, "*:AddTypeandUser")            //根据团队查询用户列表
	beego.Router("/hcloud/resource/resourcelimit", &controllers.ResourceController{}, "*:SetUpGroupResource")
	beego.Router("/hcloud/resource/groupresourcesetup", &controllers.ResourceController{}, "*:TeamResourceSetup")
	beego.Router("/hcloud/resource/personresourcesetup", &controllers.ResourceController{}, "*:PersonResourceSetup")
	//api

	// 接口请求
	beego.Router("/api", &controllers.CloudApiController{}, "*:ApiIndex")
	beego.Router("/api/v1/login", &api.LoginController{}, "*:Login")
	beego.Router("/api/v1/setcould", &api.LoginController{}, "POST:SetCouldUserName")
	beego.Router("/api/v1/device", &api.DeviceController{}, "GET:DeviceList")
	beego.Router("/api/v1/deviceorder", &api.DeviceController{}, "GET:DeviceOrder")
	beego.Router("/api/v1/devicedel", &api.DeviceController{}, "DELETE:DeviceDel")
	beego.Router("/api/v1/deviceedit", &api.DeviceController{}, "POST:DeviceEdit")
	beego.Router("/api/v1/dataapi", &controllers.CloudApiController{}, "*:Index")
	//请求资源
	beego.Router("/api/v1/requestresource", &controllers.CloudApiController{}, "*:RequestDownloadResource") //请求资源
	beego.Router("/api/v1/restype", &api.RestypeController{}, "*:RequestResJson")                           //获取资源类型
}
func testrouter() {
	// 测试
	beego.Router("/test", &controllers.TestController{}, "*:Index")
	beego.Router("/test/postapi", &controllers.TestController{}, "post:PostApi")
}
