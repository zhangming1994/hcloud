package resource

import (
	"cron"
	"github.com/astaxie/beego"
)

var mgr *Mgr

type Mgr struct {
	up   *cron.Cron
	down *cron.Cron
}

func init() {
	mgr = &Mgr{}
	mgr.up = cron.New(1024)
	mgr.down = cron.New(1024)
	mgr.up.Start()
	mgr.down.Start()
	beego.Debug("mgr:", mgr)
}

func GetMgr() *Mgr {
	return mgr
}

func (m *Mgr) UpMgr() *cron.Cron {
	return m.up
}

func (m *Mgr) DownMgr() *cron.Cron {
	return m.down
}
