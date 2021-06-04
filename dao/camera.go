package dao

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/yumrano/rtsp2rtmp/rlog"
)

func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(Camera))
}

type Camera struct {
	Id                string `orm:"pk;column(id)" json:"id"`
	Code              string `orm:"column(code)" json:"code"`
	RtspURL           string `orm:"column(rtsp_url)" json:"rtspURL"`
	RtmpURL           string `orm:"column(rtmp_url)" json:"rtmpURL"`
	AuthCodeTemp      string `orm:"column(auth_code_temp)" json:"authCodeTemp"`
	AuthCodePermanent string `orm:"column(auth_code_permanent)" json:"authCodePermanent"`
	OnlineStatus      int    `orm:"column(online_status)" json:"onlineStatus"`
}

func CameraInsert(e Camera) (i int64, err error) {
	o := orm.NewOrm()
	i, err = o.Insert(&e)
	if err != nil {
		rlog.Log.Printf("camera insert error : %v", err)
		return i, err
	}
	return i, nil
}

func CameraDelete(e Camera) (i int64, err error) {
	o := orm.NewOrm()
	i, err = o.Delete(&e)
	if err != nil {
		rlog.Log.Printf("camera delete error : %v", err)
		return 0, err
	}
	return i, nil
}

func CameraUpdate(e Camera) (i int64, err error) {
	o := orm.NewOrm()
	i, err = o.Update(&e)
	if err != nil {
		rlog.Log.Printf("camera update error : %v", err)
		return 0, err
	}
	return i, nil
}

func CameraSelectById(id string) (e Camera, err error) {
	o := orm.NewOrm()
	e = Camera{Id: id}

	err = o.Read(&e)

	if err == orm.ErrNoRows {
		rlog.Log.Println("查询不到")
		return e, err
	} else if err == orm.ErrMissPK {
		rlog.Log.Println("找不到主键")
		return e, err
	} else if err != nil {
		rlog.Log.Printf("错误: %v", err)
		return e, err
	} else {
		return e, nil
	}
}

func CameraSelectOne(q Camera) (e Camera, err error) {
	o := orm.NewOrm()
	err = o.QueryTable(new(Camera)).Filter("code", q.Code).One(&e)
	if err != nil {
		rlog.Log.Printf("查询出错：%v", err)
		return e, err
	}
	return e, nil
}

func CameraCountByCode(code string) (count int64, err error) {
	o := orm.NewOrm()
	count, err = o.QueryTable(new(Camera)).Filter("code", code).Count()
	if err != nil {
		rlog.Log.Printf("查询出错：%v", err)
		return count, err
	}
	return count, nil
}

func CameraSelectAll() (es []Camera, err error) {
	o := orm.NewOrm()
	num, err := o.QueryTable(new(Camera)).All(&es)
	if err != nil {
		rlog.Log.Printf("查询出错：%v", err)
		return es, err
	}
	rlog.Log.Printf("查询到%d条记录", num)
	return es, nil
}
