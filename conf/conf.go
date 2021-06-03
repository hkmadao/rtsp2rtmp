package conf

import (
	"errors"

	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

var confs = make(map[interface{}]interface{})

// GetInt get bool value from config file
func GetBool(s string) (bool, error) {
	keys := strings.Split(s, ".")
	tempmap := confs
	for i, key := range keys {
		if i == len(keys)-1 {
			if value, ok := tempmap[key].(bool); ok {
				return value, nil
			}
			return false, errors.New("value type not int")
		}
		tempmap = tempmap[key].(map[interface{}]interface{})
	}
	return false, errors.New("error")
}

// GetInt get int value from config file
func GetInt(s string) (int, error) {
	keys := strings.Split(s, ".")
	tempmap := confs
	for i, key := range keys {
		if i == len(keys)-1 {
			if value, ok := tempmap[key].(int); ok {
				return value, nil
			}
			return 0, errors.New("value type not int")
		}
		tempmap = tempmap[key].(map[interface{}]interface{})
	}
	return 0, errors.New("error")
}

// GetInt get string value from config file
func GetString(s string) (string, error) {
	keys := strings.Split(s, ".")
	tempmap := confs
	for i, key := range keys {
		if i == len(keys)-1 {
			if value, ok := tempmap[key].(string); ok {
				return value, nil
			}
			return "", errors.New("value type not int")
		}
		tempmap = tempmap[key].(map[interface{}]interface{})
	}
	return "", errors.New("error")
}

func init() {
	fileName := "./conf/conf.yml"
	// fileName := "./conf.yml"
	fdata, err := ioutil.ReadFile(fileName)
	log.Printf("\n%s\n", fdata)
	if err != nil {
		log.Fatalf("read config file: %s error %v", fileName, err)
	}
	if err := yaml.Unmarshal(fdata, confs); err != nil {
		log.Fatalf("parse config file: %s error %v", fileName, err)
	}
}
