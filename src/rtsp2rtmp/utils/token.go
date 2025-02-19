package utils

import (
	"fmt"
	"strings"
	"time"

	"math/rand"

	"github.com/beego/beego/v2/core/logs"
	"github.com/google/uuid"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

// start with date token
func NextToken() (string, error) {
	timestamp := time.Now().Format("20060102150405")
	id, err := uuid.NewRandom()
	if err != nil {
		logs.Error("Random error : %v", err)
		return "", err
	}
	idstring := id.String()
	idstring = strings.ReplaceAll(idstring, "-", "")
	return timestamp + "-" + idstring, nil
}

// validate token
func TokenTimeOut(token string, duration time.Duration) bool {
	tokenTimeString := token[0:14]
	if len(tokenTimeString) != 14 {
		return true
	}
	tokenTime, err := time.Parse("20060102150405", tokenTimeString)
	if err != nil {
		return true
	}
	return time.Now().After(tokenTime.Add(duration))
}

func GenerateId() (string, error) {
	idstring, err := gonanoid.New()
	if err != nil {
		logs.Error("GenerateId error : %v", err)
		return "", err
	}
	return idstring, nil
}

func GenarateRandStr(pwdLeng int) (string, error) {
	if pwdLeng == 0 {
		err := fmt.Errorf("length is zero")
		return "", err
	}
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()"
	chars := []rune(str)
	pwdStrArr := make([]string, pwdLeng)
	for i := 0; i < pwdLeng; {
		num := rand.Intn(len(chars))
		pwdStrArr = append(pwdStrArr, string(chars[num]))
		i++
	}
	return strings.Join(pwdStrArr, ""), nil
}

func GenarateRandName() string {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	pwdLeng := rand.Intn(11) + 10
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	chars := []rune(str)
	pwdStrArr := make([]string, pwdLeng)
	for i := 0; i < pwdLeng; {

		num := rand.Intn(len(chars))
		pwdStrArr = append(pwdStrArr, string(chars[num]))
		i++
	}
	return strings.Join(pwdStrArr, "")
}

func GenaratePwd() (string, error) {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	pwdLeng := rand.Intn(11) + 10
	return GenarateRandStr(pwdLeng)
}
