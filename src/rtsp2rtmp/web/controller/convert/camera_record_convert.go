package base

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
	camera_record_po "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/po/base/camera_record"
	camera_record_vo "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/vo/base/camera_record"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

func ConvertPOToCameraRecord(po camera_record_po.CameraRecordPO) (cameraRecord entity.CameraRecord, err error) {
	err = common.POToEntity(po, &cameraRecord)
	if err != nil {
		logs.Error("convertPOToCameraRecord : %v", err)
		err = fmt.Errorf("convertPOToCameraRecord : %v", err)
		return
	}
	return
}

func ConvertPOListToCameraRecord(poes []camera_record_po.CameraRecordPO) ([]entity.CameraRecord, error) {
	cameraRecords := make([]entity.CameraRecord, len(poes))
	for i, po := range poes {
		cameraRecord, err_convert := ConvertPOToCameraRecord(po)
		if err_convert != nil {
			logs.Error("ConvertPOListToCameraRecord : %v", err_convert)
			err := fmt.Errorf("ConvertPOListToCameraRecord : %v", err_convert)
			return nil, err
		}
		cameraRecords[i] = cameraRecord
	}
	return cameraRecords, nil
}

func ConvertCameraRecordToVO(cameraRecord entity.CameraRecord) (vo camera_record_vo.CameraRecordVO, err error) {
	vo = camera_record_vo.CameraRecordVO{}
	err = common.EntityToVO(cameraRecord, &vo)
	if err != nil {
		logs.Error("convertCameraRecordToVO : %v", err)
		err = fmt.Errorf("convertCameraRecordToVO : %v", err)
		return
	}
camera, err := base_service.CameraSelectById(vo.IdCamera)
	if err != nil {
		logs.Error("convertCameraRecordToVO : %v", err)
		err = fmt.Errorf("convertCameraRecordToVO : %v", err)
		return
	}
	var cameraVO = camera_record_vo.CameraVO{}
	err = common.EntityToVO(camera, &cameraVO)
	if err != nil {
		logs.Error("convertCameraRecordToVO : %v", err)
		err = fmt.Errorf("convertCameraRecordToVO : %v", err)
		return
	}
	vo.Camera = cameraVO
	
	return
}

func ConvertCameraRecordToVOList(cameraRecords []entity.CameraRecord) (voList []camera_record_vo.CameraRecordVO, err error) {
	voList = make([]camera_record_vo.CameraRecordVO, 0)
	for _, cameraRecord := range cameraRecords {
		vo, err_convert := ConvertCameraRecordToVO(cameraRecord)
		if err_convert != nil {
			logs.Error("convertCameraRecordToVO : %v", err_convert)
			err = fmt.Errorf("ConvertCameraRecordToVOList : %v", err_convert)
			return
		}
		voList = append(voList, vo)
	}
	return
}
