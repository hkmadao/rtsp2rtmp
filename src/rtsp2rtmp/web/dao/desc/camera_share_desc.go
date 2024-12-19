package desc

import (
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
)

func GetCameraShareDesc() *common.EntityDesc {
	var entityInfo = common.EntityInfo{
		Name:        "CameraShare",
		DisplayName: "摄像头分享",
		ClassName:   "CameraShare",
		TableName:   "camera_share",
		BasePath:    "entity::camera_share",
	}
	var idAttributeInfo = &common.AttributeInfo{
		ColumnName:  "id",
		Name:        "id",
		DisplayName: "摄像头分享主属性",
		DataType:    "InternalPK",
		ValueType:   "string",
	}
	var nameAttributeInfo = &common.AttributeInfo{
		ColumnName:  "name",
		Name:        "name",
		DisplayName: "名称",
		DataType:    "String",
		ValueType:   "string",
	}
	var authCodeAttributeInfo = &common.AttributeInfo{
		ColumnName:  "auth_code",
		Name:        "authCode",
		DisplayName: "权限码",
		DataType:    "String",
		ValueType:   "string",
	}
	var enabledAttributeInfo = &common.AttributeInfo{
		ColumnName:  "enabled",
		Name:        "enabled",
		DisplayName: "启用状态",
		DataType:    "Boolean",
		ValueType:   "bool",
	}
	var createdAttributeInfo = &common.AttributeInfo{
		ColumnName:  "created",
		Name:        "created",
		DisplayName: "创建时间",
		DataType:    "DateTime",
		ValueType:   "DateTime",
	}
	var startTimeAttributeInfo = &common.AttributeInfo{
		ColumnName:  "start_time",
		Name:        "startTime",
		DisplayName: "开始时间",
		DataType:    "DateTime",
		ValueType:   "DateTime",
	}
	var deadlineAttributeInfo = &common.AttributeInfo{
		ColumnName:  "deadline",
		Name:        "deadline",
		DisplayName: "结束时间",
		DataType:    "DateTime",
		ValueType:   "DateTime",
	}
	var cameraIdAttributeInfo = &common.AttributeInfo{
		ColumnName:                     "camera_id",
		Name:                           "cameraId",
		DisplayName:                    "摄像头id",
		DataType:                       "InternalFK",
		ValueType:                      "string",
		InnerAttributeName:             "camera",
		OutEntityName:                  "Camera",
		OutEntityPkAttributeName:       "id",
		OutEntityReversalAttributeName: "cameraShares",
	}
	var cameraAttributeInfo = &common.AttributeInfo{
		ColumnName:                     "camera",
		Name:                           "camera",
		DisplayName:                    "摄像头",
		DataType:                       "InternalRef",
		ValueType:                      "",
		InnerAttributeName:             "cameraId",
		OutEntityName:                  "Camera",
		OutEntityPkAttributeName:       "id",
		OutEntityReversalAttributeName: "cameraShares",
	}
	var entityDesc = &common.EntityDesc{
		EntityInfo:      entityInfo,
		PkAttributeInfo: idAttributeInfo,
		NormalFkIdAttributeInfos: []*common.AttributeInfo{
			cameraIdAttributeInfo,
		},
		NormalFkAttributeInfos: []*common.AttributeInfo{
			cameraAttributeInfo,
		},
		NormalChildren:        []*common.AttributeInfo{},
		NormalOne2OneChildren: []*common.AttributeInfo{},
		AttributeInfoMap: map[string]*common.AttributeInfo{
			"id":        idAttributeInfo,
			"name":      nameAttributeInfo,
			"authCode":  authCodeAttributeInfo,
			"enabled":   enabledAttributeInfo,
			"created":   createdAttributeInfo,
			"startTime": startTimeAttributeInfo,
			"deadline":  deadlineAttributeInfo,
			"cameraId":  cameraIdAttributeInfo,
			"camera":    cameraAttributeInfo,
		},
	}

	return entityDesc
}
