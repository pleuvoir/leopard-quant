package config

import (
	"fmt"
	"github.com/pkg/errors"
	"leopard-quant/util"
	"strconv"
)

type LoaderSub interface {
	Load() error
	GetStr(key string) string
}

type Loader interface {
	LoaderSub
	GetStrOrDefault(key string, val string) string
	GetInt(key string) (int, error)
	GetIntOrDefault(key string, val int) int
	GetBool(key string) (bool, error)
	GetBoolOrDefault(key string, val bool) bool
}

func NewConfigLoader(sub LoaderSub) Loader {
	return &absConfigLoader{sub: sub}
}

type absConfigLoader struct {
	sub LoaderSub
}

func (d *absConfigLoader) Load() error {
	err := d.sub.Load()
	if err != nil {
		return err
	}
	return nil
}

func (d *absConfigLoader) GetStr(key string) string {
	return d.sub.GetStr(key)
}

func (d *absConfigLoader) GetStrOrDefault(key string, val string) string {
	v := d.GetStr(key)
	if util.IsBlank(v) {
		return val
	}
	return v
}

func (d *absConfigLoader) GetInt(key string) (int, error) {
	str := d.GetStr(key)
	if util.IsBlank(str) {
		return 0, errors.New(fmt.Sprintf("未获取到值，key=%s", key))
	}
	return strconv.Atoi(str)
}

func (d *absConfigLoader) GetIntOrDefault(key string, val int) int {
	if v, err := d.GetInt(key); err == nil {
		return v
	}
	return val
}

func (d *absConfigLoader) GetBool(key string) (bool, error) {
	str := d.GetStr(key)
	if util.IsBlank(str) {
		return false, errors.New(fmt.Sprintf("未获取到值，key=%s", key))
	}
	return strconv.ParseBool(str)
}

func (d *absConfigLoader) GetBoolOrDefault(key string, val bool) bool {
	if v, err := d.GetBool(key); err == nil {
		return v
	}
	return val
}
