package models

import (
	"github.com/astaxie/beego"
	"utils"
)

var UrlNuclearQueue *utils.Queue

func init() {
	var err error
	UrlNuclearQueue, err = utils.QueueNew()
	if err != nil {
		beego.Warning("QueueNew: %s", err.Error())
		// Log.Warningf("QueueNew: %s", err.Error())
	} else {
		UrlNuclearQueue.Start(PushCustomerToUser)
	}
}

type PushPhone struct {
	Phonetype  int64 //资源类型
	DistriTeam int64 //分配团队
	Num        int64 //数量
	ChanData   chan ChanData
}

type ChanData struct {
	PhoneList []int64
	Count     int
	Err       error
}

// phonetype int64, name string, num int64 ,DistriTeam int64
func PushCustomerToUser(ele interface{}) {
	push := ele.(*PushPhone)

	list, count, err := GetAllNotDistriResource(push.Phonetype, push.Num, push.DistriTeam)
	var chandata ChanData
	chandata.Count = count
	chandata.PhoneList = list
	chandata.Err = err

	push.ChanData <- chandata
}
