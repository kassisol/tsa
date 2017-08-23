package client

import (
	"fmt"
	"reflect"
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
