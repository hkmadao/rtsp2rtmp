package desc

import (
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
)

func GetTokenDesc() *common.EntityDesc {
    var entityInfo = common.EntityInfo {
        Name: "Token",
        DisplayName: "令牌",
        ClassName: "Token",
        TableName: "sys_token",
        BasePath: "entity::token",
    }
    var idTokenAttributeInfo = &common.AttributeInfo {
        ColumnName: "id_sys_token",
        Name: "idToken",
        DisplayName: "令牌主属性",
        DataType: "InternalPK",
				ValueType: "string",
    };
    var usernameAttributeInfo = &common.AttributeInfo {
        ColumnName: "username",
        Name: "username",
        DisplayName: "用户名称",
        DataType: "String",
				ValueType: "string",
    };
    var nickNameAttributeInfo = &common.AttributeInfo {
        ColumnName: "nick_name",
        Name: "nickName",
        DisplayName: "昵称",
        DataType: "String",
				ValueType: "string",
    };
    var createTimeAttributeInfo = &common.AttributeInfo {
        ColumnName: "create_time",
        Name: "createTime",
        DisplayName: "创建时间",
        DataType: "DateTime",
				ValueType: "DateTime",
    };
    var tokenAttributeInfo = &common.AttributeInfo {
        ColumnName: "token",
        Name: "token",
        DisplayName: "令牌",
        DataType: "String",
				ValueType: "string",
    };
    var expiredTimeAttributeInfo = &common.AttributeInfo {
        ColumnName: "expired_time",
        Name: "expiredTime",
        DisplayName: "过期时间",
        DataType: "DateTime",
				ValueType: "DateTime",
    };
    var userInfoStringAttributeInfo = &common.AttributeInfo {
        ColumnName: "user_info_string",
        Name: "userInfoString",
        DisplayName: "用户信息序列化",
        DataType: "String",
				ValueType: "string",
    };
    var entityDesc = &common.EntityDesc {
      EntityInfo: entityInfo,
      PkAttributeInfo: idTokenAttributeInfo,
      NormalFkIdAttributeInfos: []*common.AttributeInfo{
			},
      NormalFkAttributeInfos: []*common.AttributeInfo{
			},
      NormalChildren: []*common.AttributeInfo{
			},
      NormalOne2OneChildren: []*common.AttributeInfo{
			},
      AttributeInfoMap: map[string]*common.AttributeInfo{
          "idToken": idTokenAttributeInfo,
          "username": usernameAttributeInfo,
          "nickName": nickNameAttributeInfo,
          "createTime": createTimeAttributeInfo,
          "token": tokenAttributeInfo,
          "expiredTime": expiredTimeAttributeInfo,
          "userInfoString": userInfoStringAttributeInfo,
			},
    }

    return entityDesc
}
