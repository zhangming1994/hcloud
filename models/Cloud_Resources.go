package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	// "strconv"
	"time"
)

//资源表
type CloudResources struct {
	Id                int64
	PhoneNumber       int64              `orm:"unique"`            //手机号资源
	CloudResourceType *CloudResourceType `orm:"rel(fk)"`           //
	CloudUploadRecord *CloudUploadRecord `orm:"rel(fk)"`           //
	AllocationStatus  int                `orm:"default(0);inedx;"` //分发状态 1:已经分发 0:未分发
	Downtime          int64              `orm:"type(date)"`        //分发时间
	MobilePerson      string             //手机号拥有者姓名
	DistriTeam        int64              //
}

func (r *CloudResources) TableName() string {
	return "cloud_resources"
}
func init() {
	orm.RegisterModel(new(CloudResources))
}

//根据取得所有的数量
func GetAllTeamUse(id int64) int64 {
	o := orm.NewOrm()
	num, _ := o.QueryTable("cloud_resources").Filter("DistriTeam", id).Count()
	return num
}

// 添加资源
//批次id，成功数量。重复数，失败数量，资源类型，总数，上传文件名，上传人
func AddPhoneResource(lastrecordid int64, success int64, repatnumber int64, failenumber int64, total int64) (err error) {
	err = ModifyById(lastrecordid, repatnumber, success, failenumber, total) //上传记录
	beego.Error(err)
	if err != nil {
		beego.Error("uploadrecord insert is fail", err)
		return
	}
	ModifyStatus(2, lastrecordid)
	return err
}

//检查手机号是否重复
func CheckIsExist(phone int64) (flag bool) {
	o := orm.NewOrm()
	flag = o.QueryTable("cloud_resources").Filter("PhoneNumber", phone).Exist()
	return flag
}

//select * from cloud_resources  where  phone_number in(select phone_number from cloud_resources group by phone_number having count(*)>1)
//下载资源
func GetAllNotDistriResource(resourcetype int64, num int64, distriteam int64) (lists []int64, count int, err error) {
	o := orm.NewOrm()
	count = 0
	var list []orm.Params
	var allreadypersonuse int64 //个人已经使用的资源
	var allreadyteamuse int64   //团队已经使用的资源
	timeNow := time.Now().Format("2006-01-02")
	downtime, _ := time.Parse("2006-01-02", timeNow)
	person := GetUserById(distriteam)         //取得请求数据的用户信息
	group, _ := GetGroupById(person.Group.Id) //根据组别id取得组信息

	flag := GetIsPermission(resourcetype, distriteam) //判断权限
	if flag == false {
		err = errors.New("没有下载该资源的权限")
		beego.Error("no permission for this restype")
		return lists, 0, err
	} else {
		resnum := IsEnoughResource(resourcetype) //先判断resource表未分发的该资源是否足够
		if num > resnum {
			err = errors.New("资源不足")
			beego.Error("Insufficient resources")
			return lists, 0, err
		} else {
			if num > group.OneDayLimit { //判断下载是否超过团队
				err = errors.New("资源下载量超过所属团队每日资源上限")
				beego.Error("no permission for this restype")
				return lists, 0, err
			} else {
				//判断团队是否还有资源可供下载
				flag := CheckIsDate(downtime, person.Group.Id) //根据所属团队和日期查询有没有这条数据
				if flag == false {
					allreadyteamuse = GetTeamUseResource(downtime, person.Group.Id) //从资源表取得数据
				} else {
					teaminfo, _ := GetTotalResource(downtime, person.Group.Id) //从已有记录取得数据
					allreadyteamuse = teaminfo.UsdeResource
				}
				Team, _ := GetGroupById(person.Group.Id)
				if Team.OneDayLimit-allreadyteamuse <= 0 {
					err = errors.New("所属团队已经没有资源可供下载")
					beego.Error("on team no resource can down")
					return lists, 0, err
				} else {
					if num > person.OnceLimit { //判断个人下载数量是否超过今天单次拉取值限制
						err = errors.New("个人单次下载量超过限制值")
						beego.Error("Download more than a single value limit")
						return lists, 0, err
					} else {
						flags := CheckIsPerson(downtime, distriteam, person.Group.Id)
						if flags == false {
							allreadypersonuse = GetNumByUserId(distriteam, downtime) //从resource表取得该用户已经使用的资源
						} else {
							info := GetInfo(downtime, distriteam, person.Group.Id) //从下载记录获取
							allreadypersonuse = info.UsdeResource
						}
						if num > person.OneDayLimit-allreadypersonuse {
							err = errors.New("个人下载资源已达今日上限")
							beego.Error("no resource can down")
							return lists, 0, err
						} else {
							err := o.Begin()
							qs := o.QueryTable("cloud_resources")
							_, err = qs.Limit(num).Filter("AllocationStatus", 0).Filter("CloudResourceType__Id", resourcetype).Values(&list, "PhoneNumber")
							if err != nil {
								beego.Error("get data is err", err)
								return nil, count, err
							} else {
								for _, val := range list {
									phonenumber, _ := val["PhoneNumber"].(int64)
									lists = append(lists, phonenumber)
									_, err := o.QueryTable("cloud_resources").Filter("PhoneNumber", phonenumber).Update(orm.Params{
										"AllocationStatus": 1,
										"Downtime":         time.Now().Unix(),
										"DistriTeam":       distriteam,
									})
									if err != nil {
										beego.Error("modify data is err", err)
										o.Rollback()
										return nil, 0, errors.New("资源获取失败")
									} else {
										count++
										o.Commit()
									}
								}
								if count > 0 {
									var resourRecord = new(CloudDistributRecord)
									resourRecord.DownNumber = num
									resourRecord.DistributName = person.Id
									resourRecord.DistriResourceType = resourcetype
									resourRecord.DownTime = downtime
									nums, _ := AddDistributRecord(resourRecord) //添加下载记录
									if nums < 0 {
										return lists, 0, errors.New("下载记录插入失败")
									} else {
										//更改团队资源量
										isteamresord := CheckIsDate(downtime, person.Group.Id) //根据所属团队和日期查询有没有这条数据
										if isteamresord == false {
											var teamtotalrecord = new(CloudTotalResource)
											teamtotalrecord.ResourceDate = downtime
											teamtotalrecord.Group = person.Group.Id
											teamtotalrecord.ManageUser = group.Uid
											teamtotalrecord.TotalResource = group.OneDayLimit
											teamtotalrecord.UsdeResource = num
											teamtotalrecord.CanUseResource = group.OneDayLimit - num
											persentage := fmt.Sprintf("%.1f", float64(num)/float64(group.OneDayLimit)*100)
											teamtotalrecord.UsePersent = persentage
											AddCloudResouIndex(teamtotalrecord)
										} else {
											ModifyByType(num, downtime, person.Group.Id) //修改被请求资源类型的统计(总的)
										}
										ispersonrecord := CheckIsPerson(downtime, distriteam, person.Group.Id)
										if ispersonrecord == false {
											var personuse = new(CloudPersonRecord)
											personuse.Username = person.Id
											personuse.ResourceDate = downtime
											personuse.TotalResource = person.OneDayLimit
											personuse.Team = person.Group.Id
											personuse.UsdeResource = num
											personuse.CanUseResource = person.OneDayLimit - num
											persentage := fmt.Sprintf("%.1f", float64(num)/float64(person.OneDayLimit)*100)
											personuse.UsePersent = persentage
											AddResourcePerson(personuse)
										} else {
											err := ModifyPersonResouce(downtime, num, distriteam) //修改个人资源统计
											if err != nil {
												beego.Error("change personRecord is err", err)
											}
										}
									}
								} else {
									o.Rollback()
									return lists, 0, err
								}
							}
						}
					}
				}
			}
		}
	}
	return lists, count, err
}

//取得指定用户下载的资源量
func GetNumByUserId(id int64, down time.Time) int64 {
	o := orm.NewOrm()
	num, _ := o.QueryTable("cloud_resources").Filter("DistriTeam", id).Filter("Downtime", down).Filter("AllocationStatus", 1).Count()
	return num
}

//取得团队已经使用的资源
func GetTeamUseResource(downtime time.Time, teamid int64) int64 {
	var number int64
	userlist := GetAllTeamUser(teamid)
	for i := 0; i < len(userlist); i++ {
		userid, _ := userlist[i]["Id"].(int64)
		num := GetNumByUserId(userid, downtime)
		number = number + num
	}
	return number
}

//取得所有没有分发出去的资源
func GetNotDistribut() int64 {
	o := orm.NewOrm()
	num, _ := o.QueryTable("cloud_resources").Filter("AllocationStatus", 0).Count()
	return num
}

//判断该资源在resource表中间是否足够
func IsEnoughResource(restype int64) int64 {
	o := orm.NewOrm()
	count, _ := o.QueryTable("cloud_resources").Filter("AllocationStatus", 0).Filter("CloudResourceType__Id", restype).Count()
	return count
}
