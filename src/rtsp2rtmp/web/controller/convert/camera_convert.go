package base

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
	camera_po "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/po/base/camera"
	camera_vo "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/vo/base/camera"
)

func ConvertPOToCamera(po camera_po.CameraPO) (camera entity.Camera, err error) {
	err = common.POToEntity(po, &camera)
	if err != nil {
		logs.Error("convertPOToCamera : %v", err)
		err = fmt.Errorf("convertPOToCamera : %v", err)
		return
	}
	return
}

func ConvertCameraToVO(camera entity.Camera) (vo camera_vo.CameraVO, err error) {
	vo = camera_vo.CameraVO{}
	err = common.EntityToVO(camera, &vo)
	if err != nil {
		logs.Error("convertCameraToVO : %v", err)
		err = fmt.Errorf("convertCameraToVO : %v", err)
		return
	}
	// condition := common.GetEqualCondition("cameraId", vo.Id)
	// var cameraShareVOList = make([]camera_vo.CameraShareVO, 0)
	// var cameraShares = make([]entity.CameraShare, 0)
	// cameraShares, err = base_service.CameraShareFindCollectionByCondition(condition)
	// if err != nil {
	// 	logs.Error("convertCameraToVO : %v", err)
	// 	err = fmt.Errorf("convertCameraToVO : %v", err)
	// 	return
	// }
	// for _, cameraShare := range cameraShares {
	// 	var cameraShareVO = camera_vo.CameraShareVO{}
	// 	err = common.EntityToVO(cameraShare, &cameraShareVO)
	// 	if err != nil {
	// 		logs.Error("convertCameraToVO : %v", err)
	// 		err = fmt.Errorf("convertCameraToVO : %v", err)
	// 		return
	// 	}
	// 	cameraShareVOList = append(cameraShareVOList, cameraShareVO)
	// }
	// vo.CameraShares = cameraShareVOList

	return
}

func ConvertCameraToVOList(cameras []entity.Camera) (voList []camera_vo.CameraVO, err error) {
	voList = make([]camera_vo.CameraVO, 0)
	for _, camera := range cameras {
		vo, err_convert := ConvertCameraToVO(camera)
		if err_convert != nil {
			logs.Error("convertCameraToVO : %v", err_convert)
			err = fmt.Errorf("ConvertCameraToVOList : %v", err_convert)
			return
		}
		voList = append(voList, vo)
	}
	return
}
