package base

import (
	"fmt"

	"github.com/beego/beego/v2/core/logs"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/common"
	"github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dao/entity"
	token_po "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/po/base/token"
	token_vo "github.com/hkmadao/rtsp2rtmp/src/rtsp2rtmp/web/dto/vo/base/token"
)

func ConvertPOToToken(po token_po.TokenPO) (token entity.Token, err error) {
	err = common.POToEntity(po, &token)
	if err != nil {
		logs.Error("convertPOToToken : %v", err)
		err = fmt.Errorf("convertPOToToken : %v", err)
		return
	}
	return
}

func ConvertPOListToToken(poes []token_po.TokenPO) ([]entity.Token, error) {
	tokens := make([]entity.Token, len(poes))
	for i, po := range poes {
		token, err_convert := ConvertPOToToken(po)
		if err_convert != nil {
			logs.Error("ConvertPOListToToken : %v", err_convert)
			err := fmt.Errorf("ConvertPOListToToken : %v", err_convert)
			return nil, err
		}
		tokens[i] = token
	}
	return tokens, nil
}

func ConvertTokenToVO(token entity.Token) (vo token_vo.TokenVO, err error) {
	vo = token_vo.TokenVO{}
	err = common.EntityToVO(token, &vo)
	if err != nil {
		logs.Error("convertTokenToVO : %v", err)
		err = fmt.Errorf("convertTokenToVO : %v", err)
		return
	}

	return
}

func ConvertTokenToVOList(tokens []entity.Token) (voList []token_vo.TokenVO, err error) {
	voList = make([]token_vo.TokenVO, 0)
	for _, token := range tokens {
		vo, err_convert := ConvertTokenToVO(token)
		if err_convert != nil {
			logs.Error("convertTokenToVO : %v", err_convert)
			err = fmt.Errorf("ConvertTokenToVOList : %v", err_convert)
			return
		}
		voList = append(voList, vo)
	}
	return
}
