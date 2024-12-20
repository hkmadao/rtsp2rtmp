package base

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
	user_po "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/po/base/user"
	user_vo "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/vo/base/user"
)

func ConvertPOToUser(po user_po.UserPO) (user entity.User, err error) {
	err = common.POToEntity(po, &user)
	if err != nil {
		logs.Error("convertPOToUser : %v", err)
		err = fmt.Errorf("convertPOToUser : %v", err)
		return
	}
	return
}

func ConvertPOListToUser(poes []user_po.UserPO) ([]entity.User, error) {
	users := make([]entity.User, len(poes))
	for i, po := range poes {
		user, err_convert := ConvertPOToUser(po)
		if err_convert != nil {
			logs.Error("ConvertPOListToUser : %v", err_convert)
			err := fmt.Errorf("ConvertPOListToUser : %v", err_convert)
			return nil, err
		}
		users[i] = user
	}
	return users, nil
}

func ConvertUserToVO(user entity.User) (vo user_vo.UserVO, err error) {
	vo = user_vo.UserVO{}
	err = common.EntityToVO(user, &vo)
	if err != nil {
		logs.Error("convertUserToVO : %v", err)
		err = fmt.Errorf("convertUserToVO : %v", err)
		return
	}

	return
}

func ConvertUserToVOList(users []entity.User) (voList []user_vo.UserVO, err error) {
	voList = make([]user_vo.UserVO, 0)
	for _, user := range users {
		vo, err_convert := ConvertUserToVO(user)
		if err_convert != nil {
			logs.Error("convertUserToVO : %v", err_convert)
			err = fmt.Errorf("ConvertUserToVOList : %v", err_convert)
			return
		}
		voList = append(voList, vo)
	}
	return
}
