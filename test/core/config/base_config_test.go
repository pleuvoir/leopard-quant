package config

import (
	"encoding/json"
	"fmt"
	"leopard-quant/core/config"
	"os"
	"path"
	"testing"
)

type AdaptLoadFunc func() error
type AdaptGetStrFunc func(key string) string

func (funcW AdaptLoadFunc) Load() error {
	return funcW()
}
func (funcW AdaptGetStrFunc) GetStr(key string) string {
	return funcW(key)
}

type Impl struct {
	AdaptLoadFunc
	AdaptGetStrFunc
}

// 测试函数形式实现接口
func TestFunc(t *testing.T) {
	impl := Impl{
		func() error {
			t.Log("load ..")
			return nil
		},
		func(string) string {
			t.Log("getStr ..")
			return ""
		}}

	loader := config.NewConfigLoader(&impl)
	loader.Load()

	loader.GetStr("")

	configLoader := loader.(config.Loader)
	configLoader.Load()
}

// 正常的实现，使用一个结构体
type FileConfigLoader struct {
	properties map[string]string
}

func (f *FileConfigLoader) Load() error {
	f.properties = make(map[string]string)
	f.properties["name"] = "f"
	f.properties["int"] = "2"
	return nil
}

func (f *FileConfigLoader) GetStr(key string) string {
	return f.properties[key]
}

func TestGet(t *testing.T) {

	loader := config.NewConfigLoader(&FileConfigLoader{})
	if err := loader.Load(); err != nil {
		panic(err)
	}

	getInt, err := loader.GetInt("int")

	t.Log(getInt)
	t.Log(err)

}

type SimpleJsonConfigLoader struct {
	path string
	c    map[string]any
}

func NewJsonConfigLoader(path string) config.LoaderSub {
	return &SimpleJsonConfigLoader{path: path, c: make(map[string]any)}
}

func (j *SimpleJsonConfigLoader) Load() error {
	data, err := os.ReadFile(j.path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &j.c)
	if err != nil {
		return err
	}

	return nil
}

func (j *SimpleJsonConfigLoader) GetStr(key string) string {
	return fmt.Sprint(j.c[key])
}

func TestJson(t *testing.T) {

	m := map[string]any{
		"key": "nest",
	}
	bytes, err := json.Marshal(m)

	dir, _ := os.Getwd()
	filePath := path.Join(dir, "TestJson.json")
	err = os.WriteFile(filePath, bytes, 0777)
	defer func() {
		_ = os.Remove(filePath)
	}()
	if err != nil {
		panic(err)
	}

	loader := config.NewConfigLoader(NewJsonConfigLoader(filePath))
	if err = loader.Load(); err != nil {
		panic(err)
	}

	str := loader.GetStr("key")
	t.Log(str)

}
