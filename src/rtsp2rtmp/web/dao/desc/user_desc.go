package desc

import (
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
)

func GetUserDesc() *common.EntityDesc {
    var entityInfo = common.EntityInfo {
        Name: "User",
        DisplayName: "系统用户",
        ClassName: "User",
        TableName: "sys_user",
        BasePath: "entity::user",
    }
    var idUserAttributeInfo = &common.AttributeInfo {
        ColumnName: "id_user",
        Name: "idUser",
        DisplayName: "系统用户id",
        DataType: "InternalPK",
				ValueType: "string",
    };
    var accountAttributeInfo = &common.AttributeInfo {
        ColumnName: "account",
        Name: "account",
        DisplayName: "登录账号 ",
        DataType: "String",
				ValueType: "string",
    };
    var userPwdAttributeInfo = &common.AttributeInfo {
        ColumnName: "user_pwd",
        Name: "userPwd",
        DisplayName: "用户密码 ",
        DataType: "String",
				ValueType: "string",
    };
    var phoneAttributeInfo = &common.AttributeInfo {
        ColumnName: "phone",
        Name: "phone",
        DisplayName: "手机号码",
        DataType: "String",
				ValueType: "string",
    };
    var emailAttributeInfo = &common.AttributeInfo {
        ColumnName: "email",
        Name: "email",
        DisplayName: "邮箱",
        DataType: "String",
				ValueType: "string",
    };
    var nameAttributeInfo = &common.AttributeInfo {
        ColumnName: "name",
        Name: "name",
        DisplayName: "姓名 ",
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
    var genderAttributeInfo = &common.AttributeInfo {
        ColumnName: "gender",
        Name: "gender",
        DisplayName: "性别",
        DataType: "String",
				ValueType: "string",
    };
    var fgActiveAttributeInfo = &common.AttributeInfo {
        ColumnName: "fg_active",
        Name: "fgActive",
        DisplayName: "启用标志",
        DataType: "Boolean",
				ValueType: "bool",
    };
    var entityDesc = &common.EntityDesc {
      EntityInfo: entityInfo,
      PkAttributeInfo: idUserAttributeInfo,
      NormalFkIdAttributeInfos: []*common.AttributeInfo{
			},
      NormalFkAttributeInfos: []*common.AttributeInfo{
			},
      NormalChildren: []*common.AttributeInfo{
			},
      NormalOne2OneChildren: []*common.AttributeInfo{
			},
      AttributeInfoMap: map[string]*common.AttributeInfo{
          "idUser": idUserAttributeInfo,
          "account": accountAttributeInfo,
          "userPwd": userPwdAttributeInfo,
          "phone": phoneAttributeInfo,
          "email": emailAttributeInfo,
          "name": nameAttributeInfo,
          "nickName": nickNameAttributeInfo,
          "gender": genderAttributeInfo,
          "fgActive": fgActiveAttributeInfo,
			},
    }

    return entityDesc
}
