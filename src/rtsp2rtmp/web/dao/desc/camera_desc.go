package desc

import (
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
)

func GetCameraDesc() *common.EntityDesc {
	var entityInfo = common.EntityInfo{
		Name:        "Camera",
		DisplayName: "摄像头",
		ClassName:   "Camera",
		TableName:   "camera",
		BasePath:    "entity::camera",
	}
	var idAttributeInfo = &common.AttributeInfo{
		ColumnName:  "id",
		Name:        "id",
		DisplayName: "摄像头主属性",
		DataType:    "InternalPK",
		ValueType:   "string",
	}
	var codeAttributeInfo = &common.AttributeInfo{
		ColumnName:  "code",
		Name:        "code",
		DisplayName: "code",
		DataType:    "String",
		ValueType:   "string",
	}
	var rtspUrlAttributeInfo = &common.AttributeInfo{
		ColumnName:  "rtsp_url",
		Name:        "rtspUrl",
		DisplayName: "rtsp地址",
		DataType:    "String",
		ValueType:   "string",
	}
	var rtmpUrlAttributeInfo = &common.AttributeInfo{
		ColumnName:  "rtmp_url",
		Name:        "rtmpUrl",
		DisplayName: "rtmp推送地址",
		DataType:    "String",
		ValueType:   "string",
	}
	var playAuthCodeAttributeInfo = &common.AttributeInfo{
		ColumnName:  "play_auth_code",
		Name:        "playAuthCode",
		DisplayName: "播放权限码",
		DataType:    "String",
		ValueType:   "string",
	}
	var onlineStatusAttributeInfo = &common.AttributeInfo{
		ColumnName:  "online_status",
		Name:        "onlineStatus",
		DisplayName: "在线状态",
		DataType:    "Boolean",
		ValueType:   "bool",
	}
	var enabledAttributeInfo = &common.AttributeInfo{
		ColumnName:  "enabled",
		Name:        "enabled",
		DisplayName: "启用状态",
		DataType:    "Boolean",
		ValueType:   "bool",
	}
	var rtmpPushStatusAttributeInfo = &common.AttributeInfo{
		ColumnName:  "rtmp_push_status",
		Name:        "rtmpPushStatus",
		DisplayName: "rtmp推送状态",
		DataType:    "Boolean",
		ValueType:   "bool",
	}
	var saveVideoAttributeInfo = &common.AttributeInfo{
		ColumnName:  "save_video",
		Name:        "saveVideo",
		DisplayName: "保存录像状态",
		DataType:    "Boolean",
		ValueType:   "bool",
	}
	var liveAttributeInfo = &common.AttributeInfo{
		ColumnName:  "live",
		Name:        "live",
		DisplayName: "直播状态",
		DataType:    "Boolean",
		ValueType:   "bool",
	}
	var createdAttributeInfo = &common.AttributeInfo{
		ColumnName:  "created",
		Name:        "created",
		DisplayName: "创建时间",
		DataType:    "DateTime",
		ValueType:   "string",
	}
	var cameraSharesAttributeInfo = &common.AttributeInfo{
		ColumnName:                       "camera_shares",
		Name:                             "cameraShares",
		DisplayName:                      "摄像头分享",
		DataType:                         "InternalArray",
		ValueType:                        "",
		OutEntityName:                    "CameraShare",
		OutEntityPkAttributeName:         "id",
		OutEntityReversalAttributeName:   "camera",
		OutEntityIdReversalAttributeName: "cameraId",
	}
	var entityDesc = &common.EntityDesc{
		EntityInfo:               entityInfo,
		PkAttributeInfo:          idAttributeInfo,
		NormalFkIdAttributeInfos: []*common.AttributeInfo{},
		NormalFkAttributeInfos:   []*common.AttributeInfo{},
		NormalChildren: []*common.AttributeInfo{
			cameraSharesAttributeInfo,
		},
		NormalOne2OneChildren: []*common.AttributeInfo{},
		AttributeInfoMap: map[string]*common.AttributeInfo{
			"id":             idAttributeInfo,
			"code":           codeAttributeInfo,
			"rtspUrl":        rtspUrlAttributeInfo,
			"rtmpUrl":        rtmpUrlAttributeInfo,
			"playAuthCode":   playAuthCodeAttributeInfo,
			"onlineStatus":   onlineStatusAttributeInfo,
			"enabled":        enabledAttributeInfo,
			"rtmpPushStatus": rtmpPushStatusAttributeInfo,
			"saveVideo":      saveVideoAttributeInfo,
			"live":           liveAttributeInfo,
			"created":        createdAttributeInfo,
			"cameraShares":   cameraSharesAttributeInfo,
		},
	}

	return entityDesc
}
