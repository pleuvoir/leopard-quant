package algorithm

import (
	"errors"
	"fmt"
	"github.com/gookit/color"
	"leopard-quant/util"
	"reflect"
)

var typeRegistry = make(map[string]reflect.Type)

func init() {
	color.Grayln("开始注册算法")
	subs := []TemplateSub{&NoopSub{}}
	for _, sub := range subs {
		k := sub.Name()
		v := util.GetRealType(sub)
		color.Grayln(k)
		if _, ok := typeRegistry[k]; ok {
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
	sub, ok = util.MakeInstance(r).(TemplateSub)
	if !ok {
		return nil, errors.New(fmt.Sprintf("该结构体不是algorithm.TemplateSub的实现，%s", subName))
	}
	if sub == nil {
		return nil, errors.New(fmt.Sprintf("构建该结构体获取到是nil，%s", subName))
	}
	return sub, nil
}

// Exist 是否存在
func Exist(subName string) bool {
	_, ok := typeRegistry[subName]
	return ok
}
