package client

import (
	"fmt"
	"reflect"

	"github.com/kassisol/tsa/api/types"
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

func GetDirectory(i interface{}) types.Directory {
	val := GetReflectValue(reflect.Map, i)
	v := val.Interface().(map[string]interface{})

	return types.Directory{
		CAInfo:     v["ca_info"].(string),
		NewApp:     v["new_app"].(string),
		NewAuthz:   v["new_authz"].(string),
		RevokeCert: v["revoke_cert"].(string),
	}
}
