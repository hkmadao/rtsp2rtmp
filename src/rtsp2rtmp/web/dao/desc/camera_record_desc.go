package desc

import (
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
)

func GetCameraRecordDesc() *common.EntityDesc {
	var entityInfo = common.EntityInfo{
		Name:        "CameraRecord",
		DisplayName: "摄像头记录",
		ClassName:   "CameraRecord",
		TableName:   "camera_record",
		BasePath:    "entity::camera_record",
	}
	var idRecordAttributeInfo = &common.AttributeInfo{
		ColumnName:  "id_camera_record",
		Name:        "idCameraRecord",
		DisplayName: "记录id",
		DataType:    "InternalPK",
		ValueType:   "string",
	}
	var createdAttributeInfo = &common.AttributeInfo{
		ColumnName:  "created",
		Name:        "created",
		DisplayName: "创建时间",
		DataType:    "DateTime",
		ValueType:   "DateTime",
	}
	var tempFileNameAttributeInfo = &common.AttributeInfo{
		ColumnName:  "temp_file_name",
		Name:        "tempFileName",
		DisplayName: "临时文件名称",
		DataType:    "String",
		ValueType:   "string",
	}
	var fgTempAttributeInfo = &common.AttributeInfo{
		ColumnName:  "fg_temp",
		Name:        "fgTemp",
		DisplayName: "临时文件标志",
		DataType:    "Boolean",
		ValueType:   "bool",
	}
	var fileNameAttributeInfo = &common.AttributeInfo{
		ColumnName:  "file_name",
		Name:        "fileName",
		DisplayName: "文件名称",
		DataType:    "String",
		ValueType:   "string",
	}
	var fgRemoveAttributeInfo = &common.AttributeInfo{
		ColumnName:  "fg_remove",
		Name:        "fgRemove",
		DisplayName: "文件删除标志",
		DataType:    "Boolean",
		ValueType:   "bool",
	}
	var durationAttributeInfo = &common.AttributeInfo{
		ColumnName:  "duration",
		Name:        "duration",
		DisplayName: "文件时长",
		DataType:    "Integer",
		ValueType:   "number",
	}
	var startTimeAttributeInfo = &common.AttributeInfo{
		ColumnName:  "start_time",
		Name:        "startTime",
		DisplayName: "开始时间",
		DataType:    "DateTime",
		ValueType:   "DateTime",
	}
	var endTimeAttributeInfo = &common.AttributeInfo{
		ColumnName:  "end_time",
		Name:        "endTime",
		DisplayName: "结束时间",
		DataType:    "DateTime",
		ValueType:   "DateTime",
	}
	var idCameraAttributeInfo = &common.AttributeInfo{
		ColumnName:                     "id_camera",
		Name:                           "idCamera",
		DisplayName:                    "摄像头主属性",
		DataType:                       "InternalFK",
		ValueType:                      "string",
		InnerAttributeName:             "camera",
		OutEntityName:                  "Camera",
		OutEntityPkAttributeName:       "id",
		OutEntityReversalAttributeName: "cameraRecords",
	}
	var cameraAttributeInfo = &common.AttributeInfo{
		ColumnName:                     "",
		Name:                           "camera",
		DisplayName:                    "摄像头",
		DataType:                       "InternalRef",
		ValueType:                      "",
		InnerAttributeName:             "idCamera",
		OutEntityName:                  "Camera",
		OutEntityPkAttributeName:       "id",
		OutEntityReversalAttributeName: "cameraRecords",
	}
	var entityDesc = &common.EntityDesc{
		EntityInfo:      entityInfo,
		PkAttributeInfo: idRecordAttributeInfo,
		NormalFkIdAttributeInfos: []*common.AttributeInfo{
			idCameraAttributeInfo,
		},
		NormalFkAttributeInfos: []*common.AttributeInfo{
			cameraAttributeInfo,
		},
		NormalChildren:        []*common.AttributeInfo{},
		NormalOne2OneChildren: []*common.AttributeInfo{},
		AttributeInfoMap: map[string]*common.AttributeInfo{
			"idCameraRecord": idRecordAttributeInfo,
			"created":        createdAttributeInfo,
			"tempFileName":   tempFileNameAttributeInfo,
			"fgTemp":         fgTempAttributeInfo,
			"fileName":       fileNameAttributeInfo,
			"fgRemove":       fgRemoveAttributeInfo,
			"duration":       durationAttributeInfo,
			"startTime":      startTimeAttributeInfo,
			"endTime":        endTimeAttributeInfo,
			"idCamera":       idCameraAttributeInfo,
			"camera":         cameraAttributeInfo,
		},
	}

	return entityDesc
}
