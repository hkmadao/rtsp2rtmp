package service

import (
	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
)

func CameraShareInsert(e entity.CameraShare) (i int64, err error) {
	o := orm.NewOrm()
	i, err = o.Insert(&e)
	if err != nil && err != orm.ErrLastInsertIdUnavailable {
		logs.Error("camera insert error : %v", err)
		return i, err
	}
	return i, nil
}

func CameraShareDelete(e entity.CameraShare) (i int64, err error) {
	o := orm.NewOrm()
	i, err = o.Delete(&e)
	if err != nil {
		logs.Error("camera delete error : %v", err)
		return 0, err
	}
	return i, nil
}

func CameraShareUpdate(e entity.CameraShare) (i int64, err error) {
	o := orm.NewOrm()
	i, err = o.Update(&e)
	if err != nil {
		logs.Error("camera update error : %v", err)
		return 0, err
	}
	return i, nil
}

func CameraShareSelectById(id string) (e entity.CameraShare, err error) {
	o := orm.NewOrm()
	e = entity.CameraShare{Id: id}

	err = o.Read(&e)

	if err == orm.ErrNoRows {
		logs.Error("查询不到")
		return e, err
	} else if err == orm.ErrMissPK {
		logs.Error("找不到主键")
		return e, err
	} else if err != nil {
		logs.Error("错误: %v", err)
		return e, err
	} else {
		return e, nil
	}
}

func CameraShareSelectOne(q entity.CameraShare) (e entity.CameraShare, err error) {
	o := orm.NewOrm()
	err = o.QueryTable(new(entity.CameraShare)).Filter("CameraId", q.CameraId).Filter("AuthCode", q.AuthCode).One(&e)
	if err != nil {
		logs.Error("查询出错：%v", err)
		return e, err
	}
	return e, nil
}

func CameraShareCountByCode(code string) (count int64, err error) {
	o := orm.NewOrm()
	count, err = o.QueryTable(new(entity.CameraShare)).Filter("code", code).Count()
	if err != nil {
		logs.Error("查询出错：%v", err)
		return count, err
	}
	return count, nil
}

func CameraShareSelectAll() (es []entity.CameraShare, err error) {
	o := orm.NewOrm()
	num, err := o.QueryTable(new(entity.CameraShare)).All(&es)
	if err != nil {
		logs.Error("查询出错：%v", err)
		return es, err
	}
	logs.Debug("查询到%d条记录", num)
	return es, nil
}

func CameraShareSelectByCameraId(cameraId string) (es []entity.CameraShare, err error) {
	o := orm.NewOrm()
	num, err := o.QueryTable(new(entity.CameraShare)).Filter("CameraId", cameraId).All(&es)
	if err != nil {
		logs.Error("查询出错：%v", err)
		return es, err
	}
	logs.Info("查询到%d条记录", num)
	return es, nil
}
