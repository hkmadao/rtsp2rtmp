package service

import (
	"fmt"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dyn_query"
)

func CameraCreate(e entity.Camera) (i int64, err error) {
	o := orm.NewOrm()
	i, err = o.Insert(&e)
	if err != nil && err != orm.ErrLastInsertIdUnavailable {
		logs.Error("camera insert error : %v", err)
		return i, err
	}
	return i, nil
}

func CameraUpdateById(e entity.Camera) (i int64, err error) {
	o := orm.NewOrm()
	i, err = o.Update(&e)
	if err != nil {
		logs.Error("camera update error : %v", err)
		return 0, err
	}
	return i, nil
}

func CameraDelete(e entity.Camera) (i int64, err error) {
	o := orm.NewOrm()
	i, err = o.Delete(&e)
	if err != nil {
		logs.Error("camera delete error : %v", err)
		return 0, err
	}
	return i, nil
}

func CameraSelectById(id string) (e entity.Camera, err error) {
	o := orm.NewOrm()
	e = entity.Camera{Id: id}

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

func CameraSelectOne(q entity.Camera) (e entity.Camera, err error) {
	o := orm.NewOrm()
	err = o.QueryTable(new(entity.Camera)).Filter("code", q.Code).One(&e)
	if err != nil {
		logs.Error("查询出错：%v", err)
		return e, err
	}
	return e, nil
}

func CameraCountByCode(code string) (count int64, err error) {
	o := orm.NewOrm()
	count, err = o.QueryTable(new(entity.Camera)).Filter("code", code).Count()
	if err != nil {
		logs.Error("查询出错：%v", err)
		return count, err
	}
	return count, nil
}

func CameraSelectAll() (es []entity.Camera, err error) {
	o := orm.NewOrm()
	// num, err := o.QueryTable(new(entity.Camera)).All(&es)

	qb, _ := orm.NewQueryBuilder("postgres")

	// Construct query object
	qb.Select("*").
		From("camera").
		LeftJoin("camera_share").On("camera_share.camera_id = camera.id").
		Where("camera.code like ?").
		// OrderBy("camera.id").Desc().
		Limit(1000).Offset(0)

	// export raw query string from QueryBuilder object
	sql := qb.String()

	// execute the raw query string
	o.Raw(sql, "%%").QueryRows(&es)

	if err != nil {
		logs.Error("查询出错：%v", err)
		return es, err
	}
	logs.Debug("查询到%d条记录", 20)
	return es, nil
}

func CameraFindCollectionByCondition(condition common.AqCondition) (models []entity.Camera, err error) {
	var sqlStr, params, err_make_sql = dyn_query.MakeSqlByCondition(condition, "Camera")
	if err_make_sql != nil {
		err = fmt.Errorf("make sql error: %v", err_make_sql)
		return
	}
	o := orm.NewOrm()
	// execute the raw query string
	recordNum, err_query := o.Raw(sqlStr, params...).QueryRows(&models)
	o.Raw(sqlStr, params...).Exec()
	if err_query != nil {
		err = fmt.Errorf("query error: %v", err_make_sql)
		return
	}
	logs.Debug("查询到%d条记录", recordNum)
	return
}
