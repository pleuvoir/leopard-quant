package base

import (
	"fmt"
	"github.com/pkg/errors"
	"leopard-quant/util"
	"strconv"
)

//type FileConfigLoader struct {
//	*AbsConfigLoader
//}

type FileConfigLoader struct {
	AbsConfigLoader
	properties map[string]string
}

func (f *FileConfigLoader) load() error {
	f.properties = make(map[string]string)
	f.properties["name"] = "f"
	f.properties["int"] = "2"
	return nil
}

func (f *FileConfigLoader) getStr(key string) string {
	return f.properties[key]
}

type configLoader interface {
	load() error
	getStr(key string) string
	getStrOrDefault(key string, val string) string
	getInt(key string) (int, error)
	getIntOrDefault(key string, val int) int
	getBool(key string) (bool, error)
	getBoolOrDefault(key string, val bool) bool
}

type AbsConfigLoader struct {
	configLoader
}

func (d *AbsConfigLoader) getStrOrDefault(key string, val string) string {
	v := d.getStr(key)
	if util.IsBlank(v) {
		return val
	}
	return v
}

func (d *AbsConfigLoader) getInt(key string) (int, error) {
	str := d.getStr(key)
	if util.IsBlank(str) {
		return 0, errors.New(fmt.Sprintf("未获取到值，key=%s", key))
	}
	return strconv.Atoi(str)
}

func (d *AbsConfigLoader) getIntOrDefault(key string, val int) int {
	if v, err := d.getInt(key); err == nil {
		return v
	}
	return val
}

func (d *AbsConfigLoader) getBool(key string) (bool, error) {
	str := d.getStr(key)
	if util.IsBlank(str) {
		return false, errors.New(fmt.Sprintf("未获取到值，key=%s", key))
	}
	return strconv.ParseBool(str)
}

func (d *AbsConfigLoader) getBoolOrDefault(key string, val bool) bool {
	if v, err := d.getBool(key); err == nil {
		return v
	}
	return val
}
