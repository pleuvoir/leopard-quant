package util

import "reflect"

// GetRealType 解指针获取真实类型
func GetRealType(any interface{}) reflect.Type {
	v := reflect.ValueOf(any)
	if v.Type().Kind() == reflect.Ptr {
		return reflect.Indirect(v).Type()
	}
	return v.Type()
}

// MakeInstance 构建实例
func MakeInstance(p reflect.Type) any {
	return reflect.New(p).Interface()
}
