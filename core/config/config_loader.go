package config

import (
	"github.com/gookit/color"
	"gopkg.in/yaml.v3"
	"os"
)

type LoadState int

const (
	globalSettingFileName           = `global_setting.json`
	Unload                LoadState = -1 // 未加载
	Loading                         = 0  // 加载中
	Loaded                          = 1  // 已加载
)

type GlobalConfig struct {
	Environment string
	EnvironmentConfig
	//配置加载状态
	load LoadState
}

func (c *GlobalConfig) loadFromYml(path string) error {
	color.Println("<light_green>load config from:</>", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := &ConfContent{}
	err = yaml.Unmarshal(data, content)
	if err != nil {
		return err
	}
	c.Environment = content.Environment
	c.findEnvironmentConfig(content, content.Environment)
	return nil
}

func (c *GlobalConfig) findEnvironmentConfig(content *ConfContent, env string) {
	for _, config := range content.Configurations {
		if config.Profile == env {
			c.EnvironmentConfig = config
			color.Printf("<light_green>global config:</> %+v\n", c.EnvironmentConfig)
			return
		}
	}
}

func (c *GlobalConfig) IsProd() bool {
	return c.Environment == "prod"
}

func (c *GlobalConfig) IsTest() bool {
	return c.Environment == "test"
}

type Database struct {
	Name string `json:"name" yaml:"name"`
	Url  string `json:"url" yaml:"url"`
}

type Server struct {
	Port string `json:"port" yaml:"port"`
}

type Logger struct {
	Level        string `json:"level" yaml:"level"`
	Path         string `json:"path" yaml:"path"`
	Filename     string `json:"filename" yaml:"filename"`
	MaxAge       string `json:"maxAge" yaml:"maxAge"`
	RotationTime string `json:"rotationTime" yaml:"rotationTime"`
}

type EnvironmentConfig struct {
	Profile  string   `json:"profile" yaml:"profile"`
	Database Database `json:"database" yaml:"database"`
	Server   Server   `json:"server" yaml:"server"`
	Logger   Logger   `json:"logger" yaml:"logger"`
}

type ConfContent struct {
	Environment    string              `json:"environment" yaml:"environment"`
	Configurations []EnvironmentConfig `json:"configurations" yaml:"configurations"`
}
