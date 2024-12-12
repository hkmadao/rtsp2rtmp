package base

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
	camera_share_po "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/po/base/camera_share"
	camera_share_vo "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/vo/base/camera_share"
	base_service "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/service/base"
)

func ConvertPOToCameraShare(po camera_share_po.CameraSharePO) (camerashare entity.CameraShare, err error) {
	err = common.POToEntity(po, &camerashare)
	if err != nil {
		logs.Error("convertPOToCameraShare : %v", err)
		err = fmt.Errorf("convertPOToCameraShare : %v", err)
		return
	}
	return
}

func ConvertCameraShareToVO(camerashare entity.CameraShare) (vo camera_share_vo.CameraShareVO, err error) {
	vo = camera_share_vo.CameraShareVO{}
	err = common.EntityToVO(camerashare, &vo)
	if err != nil {
		logs.Error("convertCameraShareToVO : %v", err)
		err = fmt.Errorf("convertCameraShareToVO : %v", err)
		return
	}
	camera, err := base_service.CameraSelectById(vo.CameraId)
	if err != nil {
		logs.Error("convertCameraShareToVO : %v", err)
		err = fmt.Errorf("convertCameraShareToVO : %v", err)
		return
	}
	var cameraVO = camera_share_vo.CameraVO{}
	err = common.EntityToVO(camera, &cameraVO)
	if err != nil {
		logs.Error("convertCameraShareToVO : %v", err)
		err = fmt.Errorf("convertCameraShareToVO : %v", err)
		return
	}
	vo.Camera = cameraVO
	return
}

func ConvertCameraShareToVOList(camerashares []entity.CameraShare) (voList []camera_share_vo.CameraShareVO, err error) {
	voList = make([]camera_share_vo.CameraShareVO, 0)
	for _, camerashare := range camerashares {
		vo, err_convert := ConvertCameraShareToVO(camerashare)
		if err_convert != nil {
			logs.Error("ConvertCameraShareToVOList : %v", err_convert)
			err = fmt.Errorf("ConvertCameraShareToVOList : %v", err_convert)
			return
		}
		voList = append(voList, vo)
	}
	return
}
