package common

import (
	"fmt"
	"reflect"
	"strings"
)

func POToEntity(po interface{}, entity interface{}) (err error) {
	entityVal := reflect.ValueOf(entity).Elem()
	if entityVal.Kind() != reflect.Struct {
		err = fmt.Errorf("entity not struct")
		return
	}

	poType := reflect.TypeOf(po)
	poVal := reflect.ValueOf(po)
	if poVal.Kind() != reflect.Struct {
		err = fmt.Errorf("po not struct")
		return
	}

	fieldNum := poVal.NumField()

	for i := 0; i < fieldNum; i++ {
		poField := poType.Field(i)
		tagVal := poField.Tag.Get("po")
		if strings.Contains(tagVal, "ignore") {
			continue
		}
		poFieldVal := poVal.FieldByName(poField.Name)
		entityVal.FieldByName(poField.Name).Set(poFieldVal)
	}

	return
}

func EntityToVO(entity interface{}, vo interface{}) (err error) {
	voVal := reflect.ValueOf(vo).Elem()
	voType := reflect.TypeOf(vo).Elem()
	if voVal.Kind() != reflect.Struct {
		err = fmt.Errorf("vo not struct")
		return
	}

	entityVal := reflect.ValueOf(entity)
	if entityVal.Kind() != reflect.Struct {
		err = fmt.Errorf("entity not struct")
		return
	}

	fieldNum := voVal.NumField()

	for i := 0; i < fieldNum; i++ {
		voField := voType.Field(i)
		tagVal := voField.Tag.Get("vo")
		if strings.Contains(tagVal, "ignore") {
			continue
		}
		entityFieldVal := entityVal.FieldByName(voField.Name)
		voVal.FieldByName(voField.Name).Set(entityFieldVal)
	}

	return
}
