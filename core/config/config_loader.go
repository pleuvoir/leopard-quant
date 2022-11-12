package config

import (
	"encoding/json"
	"errors"
	"github.com/gookit/color"
	"gopkg.in/yaml.v3"
	"leopard-quant/util"
	"os"
	"path/filepath"
)

type LoadState int

const (
	ApplicationEnvVar               = "application" //环境变量中配置文件路径
	globalSettingFileName           = `application.yml`
	Unload                LoadState = -1 // 未加载
	Loading                         = 0  // 加载中
	Loaded                          = 1  // 已加载
)

type ApplicationConf struct {
	Environment string
	EnvironmentConf
	load LoadState //配置加载状态
}

func NewApplicationConf() *ApplicationConf {
	return &ApplicationConf{load: Unload}
}

// Load 加载配置
func (c *ApplicationConf) Load() error {
	var err error = nil
	if c.load == Loading {
		return errors.New("正在加载应用配置，请不要重复调用Load")
	}
	if c.load == Loaded {
		return errors.New("应用配置已加载，请不要重复调用Load")
	}
	c.load = Loading

	//返回前，根据err是否为空修改加载状态
	defer func() {
		if err != nil {
			c.load = Unload
		} else {
			c.load = Loaded
		}
	}()

	//首先，尝试从环境变量中读取配置文件
	var isLoad bool
	if isLoad, err = c.loadFromEnvVar(); err != nil {
		return err
	}

	if isLoad {
		return nil
	}
	var basePath string
	//从当前目录下读取，可执行目录
	if basePath, err = c.currentPath(); err != nil {
		return err
	}

	if err = c.tryLoad(basePath); err != nil {
		return err
	}

	//从项目根目录下读取
	if basePath, err = c.rootPath(); err != nil {
		return err
	}
	if err = c.tryLoad(basePath); err != nil {
		return err
	}

	err = errors.New("configuration file not found")
	return err
}

func (c *ApplicationConf) tryLoad(basePath string) error {
	var err error
	configFile := filepath.Join(basePath, "build", "application.yml")
	if util.IsExists(configFile) {
		if err = c.loadFromYml(configFile); err != nil {
			return err
		}
		return nil
	}
	configFile = filepath.Join(basePath, "build", "application.yaml")
	if util.IsExists(configFile) {
		if err = c.loadFromYml(configFile); err != nil {
			return err
		}
		return nil
	}
	configFile = filepath.Join(basePath, "build", "application.json")
	if util.IsExists(configFile) {
		if err = c.loadFromJson(configFile); err != nil {
			return err
		}
		return nil
	}
	return err
}

func (c *ApplicationConf) loadFromEnvVar() (bool, error) {
	configFile := c.fileFromEnvVar()
	if configFile == "" {
		return false, nil
	}
	color.Println("<light_green>load config from environment variable:</>", ApplicationEnvVar, configFile)
	extName := filepath.Ext(configFile)
	if extName == ".yaml" || extName == ".yml" {
		return true, c.loadFromYml(configFile)
	} else if extName == ".json" {
		return true, c.loadFromJson(configFile)
	}
	return false, nil
}

// 从环境变量中获取配置文件路径
func (c *ApplicationConf) fileFromEnvVar() string {
	configPath := os.Getenv(ApplicationEnvVar)
	if configPath != "" {
		if util.IsExists(configPath) {
			return configPath
		}
	}
	return ""
}

func (c *ApplicationConf) rootPath() (string, error) {
	dir, err := filepath.Abs("../")
	if err != nil {
		return "", err
	}
	return filepath.Dir(dir), nil
}

func (c *ApplicationConf) currentPath() (string, error) {
	dir, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(dir), nil
}

func (c *ApplicationConf) loadFromYml(path string) error {
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

func (c *ApplicationConf) loadFromJson(path string) error {
	color.Println("<light_green>load config from:</>", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	content := &ConfContent{}
	err = json.Unmarshal(data, content)
	if err != nil {
		return err
	}
	c.Environment = content.Environment
	c.findEnvironmentConfig(content, content.Environment)
	return nil
}

func (c *ApplicationConf) findEnvironmentConfig(content *ConfContent, env string) {
	for _, config := range content.Configurations {
		if config.Profile == env {
			c.EnvironmentConf = config
			color.Printf("<light_green>global config:</> %+v\n", c.EnvironmentConf)
			return
		}
	}
}

func (c *ApplicationConf) IsProd() bool {
	return c.Environment == "prod"
}

func (c *ApplicationConf) IsTest() bool {
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

type EnvironmentConf struct {
	Profile  string   `json:"profile" yaml:"profile"`
	Database Database `json:"database" yaml:"database"`
	Server   Server   `json:"server" yaml:"server"`
	Logger   Logger   `json:"logger" yaml:"logger"`
}

type ConfContent struct {
	Environment    string            `json:"environment" yaml:"environment"`
	Configurations []EnvironmentConf `json:"configurations" yaml:"configurations"`
}
