## hcloud

基于beego的一个群控服务器云端，系统作为htrans的服务器端

VERSION = "V0.1.1208"

## 其它相关

## 获取安装
go get git@git.haoyue.me:haphupan/hcloud.git

执行以下命令
首先你应该先有beego 环境。
 
1.然后把源码放在gopath的src目录下。

2.利用go run 运行程序，或bee run (若无法执行bee run ,请下载 go get [github.com/beego/bee](https://github.com/beego/bee))

## 接口文档

具体的接口参考 [Wiki](http://git.haoyue.me:8080/haphupan/hcloud/wiki)

## 功能模块
1. 用户管理
2. 设备管理
3. 数据运营
4. 资源管理
5. 权限管理

## 命名方式
1. 后台管理数据库的命名为`cloud_`前缀,models文件命名不需要使用，其他端命名再议
2. `CommonController`函数命名使用大写命名，并且为当前ComonController的方法
3. 数据库关联必要时候必须关联，当涉及到有`默认`等类别时候可以不使用关联
4. `Router`路径命名，方便以后权限控制，路径命名使用`\hcloud\模块名\方法`,路径请求全部使用小写，公共方法使用`\方法`
5. 模板文件以文件夹的形式管理，一个模块一个文件夹，公共文件放置跟目录`index.html,login.html,test.html`
6. 第三方扩展的库使用`vendor`方式存放，使用方式`import "tools"`即可
7. `static`静态目录中js第三方库存放至文件夹进行管理
8. 注意命名规则，每个函数至少有一个函数注释# hcloud
