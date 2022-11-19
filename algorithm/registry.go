package algorithm

import (
	"errors"
	"fmt"
	"github.com/gookit/color"
	. "leopard-quant/algorithm/impl"
	"reflect"
)

var typeRegistry = make(map[string]reflect.Type)

func init() {
	color.Grayln("开始注册算法")
	subs := []TemplateSub{NoopSub{}}
	for _, sub := range subs {
		k := sub.Name()
		v := reflect.TypeOf(sub)
		color.Grayln(k)
		if _, ok := typeRegistry[k]; ok {
			color.Redf(fmt.Sprintf("错误，发现重复的算法名称，请检查。%s", k))
			panic(fmt.Sprintf("错误，发现重复的算法名称，请检查。%s", k))
		}
		typeRegistry[k] = v
	}
	color.Grayln("算法注册完成")
}

func MakeInstance(subName string) (sub TemplateSub, err error) {
	r, ok := typeRegistry[subName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("注册表未找到该算法，%s", subName))
	}
	v := reflect.New(r).Elem()
	sub, ok = v.Interface().(TemplateSub)
	if !ok {
		return nil, errors.New(fmt.Sprintf("该结构体不是algorithm.TemplateSub的实现，%s", subName))
	}
	if sub == nil {
		return nil, errors.New(fmt.Sprintf("构建该结构体获取到是nil，%s", subName))
	}
	return sub, nil
}
