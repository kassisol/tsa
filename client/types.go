package client

import (
	"fmt"
	"reflect"

	"github.com/kassisol/tsa/api"
)

func GetReflectValue(k reflect.Kind, i interface{}) reflect.Value {
	val := reflect.ValueOf(i)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != k {
		fmt.Printf("%v type can't have attributes inspected\n", val.Kind())
		return reflect.Value{}
	}

	return val
}

func GetReflectStringValue(i interface{}) string {
	val := GetReflectValue(reflect.String, i)

	return val.String()
}

func GetDirectory(i interface{}) api.Directory {
	val := GetReflectValue(reflect.Map, i)
	v := val.Interface().(map[string]interface{})

	return api.Directory{
		CAInfo:     v["CAInfo"].(string),
		NewApp:     v["NewApp"].(string),
		NewAuthz:   v["NewAuthz"].(string),
		RevokeCert: v["RevokeCert"].(string),
	}
}
