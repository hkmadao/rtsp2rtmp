package desc

import (
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
)

func GetCameraDesc() *common.EntityDesc {
    var entityInfo = common.EntityInfo {
        Name: "Camera",
        DisplayName: "摄像头",
        ClassName: "Camera",
        TableName: "camera",
        BasePath: "entity::camera",
    }
    var idAttributeInfo = &common.AttributeInfo {
        ColumnName: "id",
        Name: "id",
        DisplayName: "摄像头主属性",
        DataType: "InternalPK",
				ValueType: "string",
    };
    var codeAttributeInfo = &common.AttributeInfo {
        ColumnName: "code",
        Name: "code",
        DisplayName: "编号",
        DataType: "String",
				ValueType: "string",
    };
    var rtspUrlAttributeInfo = &common.AttributeInfo {
        ColumnName: "rtsp_url",
        Name: "rtspUrl",
        DisplayName: "rtsp地址",
        DataType: "String",
				ValueType: "string",
    };
    var rtmpUrlAttributeInfo = &common.AttributeInfo {
        ColumnName: "rtmp_url",
        Name: "rtmpUrl",
        DisplayName: "rtmp推送地址",
        DataType: "String",
				ValueType: "string",
    };
    var playAuthCodeAttributeInfo = &common.AttributeInfo {
        ColumnName: "play_auth_code",
        Name: "playAuthCode",
        DisplayName: "播放权限码",
        DataType: "String",
				ValueType: "string",
    };
    var onlineStatusAttributeInfo = &common.AttributeInfo {
        ColumnName: "online_status",
        Name: "onlineStatus",
        DisplayName: "在线状态",
        DataType: "Boolean",
				ValueType: "bool",
    };
    var enabledAttributeInfo = &common.AttributeInfo {
        ColumnName: "enabled",
        Name: "enabled",
        DisplayName: "启用状态",
        DataType: "Boolean",
				ValueType: "bool",
    };
    var rtmpPushStatusAttributeInfo = &common.AttributeInfo {
        ColumnName: "rtmp_push_status",
        Name: "rtmpPushStatus",
        DisplayName: "rtmp推送状态",
        DataType: "Boolean",
				ValueType: "bool",
    };
    var saveVideoAttributeInfo = &common.AttributeInfo {
        ColumnName: "save_video",
        Name: "saveVideo",
        DisplayName: "保存录像状态",
        DataType: "Boolean",
				ValueType: "bool",
    };
    var liveAttributeInfo = &common.AttributeInfo {
        ColumnName: "live",
        Name: "live",
        DisplayName: "直播状态",
        DataType: "Boolean",
				ValueType: "bool",
    };
    var createdAttributeInfo = &common.AttributeInfo {
        ColumnName: "created",
        Name: "created",
        DisplayName: "创建时间",
        DataType: "DateTime",
				ValueType: "DateTime",
    };
    var fgSecretAttributeInfo = &common.AttributeInfo {
        ColumnName: "fg_secret",
        Name: "fgSecret",
        DisplayName: "加密标志",
        DataType: "Boolean",
				ValueType: "bool",
    };
    var secretAttributeInfo = &common.AttributeInfo {
        ColumnName: "secret",
        Name: "secret",
        DisplayName: "密钥",
        DataType: "String",
				ValueType: "string",
    };
    var fgPassiveAttributeInfo = &common.AttributeInfo {
        ColumnName: "fg_passive",
        Name: "fgPassive",
        DisplayName: "被动推送rtmp标志",
        DataType: "Boolean",
				ValueType: "bool",
    };
    var rtmpAuthCodeAttributeInfo = &common.AttributeInfo {
        ColumnName: "rtmp_auth_code",
        Name: "rtmpAuthCode",
        DisplayName: "rtmp识别码",
        DataType: "String",
				ValueType: "string",
    };
    var cameraTypeAttributeInfo = &common.AttributeInfo {
        ColumnName: "camera_type",
        Name: "cameraType",
        DisplayName: "摄像头类型",
        DataType: "String",
				ValueType: "string",
    };
    var cameraSharesAttributeInfo = &common.AttributeInfo {
        ColumnName: "",
        Name: "cameraShares",
        DisplayName: "摄像头分享",
        DataType: "InternalArray",
				ValueType: "",
        OutEntityName: "CameraShare",
        OutEntityPkAttributeName: "id",
        OutEntityReversalAttributeName: "camera",
        OutEntityIdReversalAttributeName: "cameraId",
    };
    var entityDesc = &common.EntityDesc {
      EntityInfo: entityInfo,
      PkAttributeInfo: idAttributeInfo,
      NormalFkIdAttributeInfos: []*common.AttributeInfo{
			},
      NormalFkAttributeInfos: []*common.AttributeInfo{
			},
      NormalChildren: []*common.AttributeInfo{
          cameraSharesAttributeInfo,
			},
      NormalOne2OneChildren: []*common.AttributeInfo{
			},
      AttributeInfoMap: map[string]*common.AttributeInfo{
          "id": idAttributeInfo,
          "code": codeAttributeInfo,
          "rtspUrl": rtspUrlAttributeInfo,
          "rtmpUrl": rtmpUrlAttributeInfo,
          "playAuthCode": playAuthCodeAttributeInfo,
          "onlineStatus": onlineStatusAttributeInfo,
          "enabled": enabledAttributeInfo,
          "rtmpPushStatus": rtmpPushStatusAttributeInfo,
          "saveVideo": saveVideoAttributeInfo,
          "live": liveAttributeInfo,
          "created": createdAttributeInfo,
          "fgSecret": fgSecretAttributeInfo,
          "secret": secretAttributeInfo,
          "fgPassive": fgPassiveAttributeInfo,
          "rtmpAuthCode": rtmpAuthCodeAttributeInfo,
          "cameraType": cameraTypeAttributeInfo,
          "cameraShares": cameraSharesAttributeInfo,
			},
    }

    return entityDesc
}
