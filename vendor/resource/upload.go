package resource

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/axgle/mahonia"
	m "hcloud/models"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type UpTask struct {
	id            int64  //批次Id
	localpath     string //存储地址
	phonetype     int64  //资源类型
	uploaduser    int64
	localfilename string //上传文件名
	bar           int64
}

func (t *UpTask) Run() {
	o := orm.NewOrm()
	file, err := os.Open(t.localpath + t.localfilename)
	defer file.Close()
	if err != nil {
		beego.Error("read the file is error: ", err.Error())
		return
	}
	var repatnumber int64 //重复数量
	var failenumber int64 //失败数量
	var success int64     //成功数量
	var successNum int64  //实际插入数据库成功数量
	var total int64       //总数
	var list []string     //每次取出来的数组
	// uploadname := t.localfilename
	//得出文件大小
	files, _ := file.Stat()
	fileSize := files.Size()
	//循环次数
	sizelangth := beego.AppConfig.String("size")
	beego.Debug(sizelangth)
	size, _ := strconv.Atoi(sizelangth)
	num := (int(fileSize) + size - 1) / size
	//百分比进度条
	var persentage float64 = 0
	for i := 0; i < num; i++ {
		persentage = persentage + 100/float64(num)
		buffer := make([]byte, size)
		_, err = file.Read(buffer)
		if err != nil {
			fmt.Println(err)
			return
		}
		decoder := mahonia.NewDecoder("gb18030")
		if decoder == nil {
			fmt.Println("the code is not exist: ", err.Error())
			return
		}
		decodelist := decoder.ConvertString(string(buffer))
		liststring := strings.Split(decodelist, "\n")
		length := len(liststring)
		ju := []byte(liststring[length-1])
		count := bytes.Count(ju, nil)

		list = liststring[:length-1]
		newposition, _ := file.Seek(-int64(count), 1)
		beego.Info(newposition)
		total = total + int64(len(list))
		var phone, name string
		var sliceRes []m.CloudResources
		for i := 0; i < len(list); i++ {
			listone := list[i]
			phonenames := strings.Split(listone, " ")
			if len(phonenames) == 2 {
				phone = phonenames[0]
				name = phonenames[1]
			} else {
				phone = phonenames[0]
				name = ""
			}
			//虚拟号段 联通：1709 移动：1705 电信1700
			//1[3|5|7|8|][\d]{9}
			reg := regexp.MustCompile(`(13[0-9]|15[012356789]|17[678]|18[0-9])[0-9]{8}`)
			checkphone := reg.MatchString(phone)
			if checkphone == false {
				failenumber++
				continue
			}

			phonenumber, _ := strconv.ParseInt(strings.TrimSpace(phone), 10, 64)

			if !m.CheckIsExist(phonenumber) {
				var pr m.CloudResources
				ty := new(m.CloudResourceType)
				rec := new(m.CloudUploadRecord)
				rec.Id = t.id
				ty.Id = t.phonetype
				pr.CloudUploadRecord = rec
				pr.PhoneNumber = phonenumber
				pr.CloudResourceType = ty
				pr.Downtime = 0
				pr.MobilePerson = name
				flags := CheckIS(sliceRes, phonenumber)
				if flags == true {
					repatnumber++
				} else {
					sliceRes = append(sliceRes, pr)
					success = success + 1 //成功数量
				}
			} else {
				repatnumber++
			}
		}
		beego.Debug(len(sliceRes))
		getnum := int(success)
		successNums, err := o.InsertMulti(getnum, sliceRes)
		if err != nil {
			beego.Error(err)
		}
		successNum = successNum + successNums
		beego.Info(successNum)
		m.UpdateBar(t.id, total, success, failenumber, repatnumber, persentage)
	}
	failenumber = total - success
	err = m.AddPhoneResource(t.id, success, repatnumber, failenumber, total)
	m.UpdateBar(t.id, total, success, failenumber, repatnumber, persentage)
}

//判断结构体是否存在这个号码
func CheckIS(cloures []m.CloudResources, phonenumber int64) bool {
	for _, val := range cloures {
		res := val.PhoneNumber
		if res == phonenumber {
			return true
		}
	}
	return false
}

// 任务
func RunUpTask(id, phonetype, user int64, path, filename string) {
	t := &UpTask{
		id:            id,
		localpath:     path,
		localfilename: filename,
		phonetype:     phonetype,
		bar:           0,
		uploaduser:    user,
	}
	GetMgr().UpMgr().AddJob(t)
}
