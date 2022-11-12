package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gookit/color"
	"gopkg.in/yaml.v3"
	"leopard-quant/util"
	"os"
	"path/filepath"
)

type LoadState int

const (
	ApplicationEnvVar                  = "application.path"    //环境变量中配置文件路径
	ApplicationEnvProfileVar           = `application.profile` //环境变量中配置文件profile
	Unload                   LoadState = -1                    // 未加载
	Loading                            = 0                     // 加载中
	Loaded                             = 1                     // 已加载
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
	var err error
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

	//获取profile
	var profile string
	profile, err = c.fileProfileFromEnvVar()
	if err != nil {
		return err
	}
	c.Environment = profile

	var basePath string

	//从当前目录下读取，可执行目录
	if basePath, err = c.currentPath(); err != nil {
		return err
	}

	if exist, loadErr := c.tryLoad(basePath); exist && loadErr == nil {
		return nil
	}

	//从项目根目录下读取
	if basePath, err = c.rootPath(); err != nil {
		return err
	}

	if exist, loadErr := c.tryLoad(basePath); exist && loadErr == nil {
		return nil
	}

	err = errors.New("configuration file not found")
	return nil
}

// 尝试加载文件，存在返回 true 不存在返回 false
func (c *ApplicationConf) tryLoad(basePath string) (isLoad bool, loadErr error) {
	configFile := filepath.Join(basePath, "build", fmt.Sprintf("application-%s.yml", c.Environment))
	if util.IsExists(configFile) {
		if loadErr = c.loadFromYml(configFile); loadErr != nil {
			return true, loadErr
		}
		return true, nil
	}
	configFile = filepath.Join(basePath, "build", fmt.Sprintf("application-%s.json", c.Environment))
	if util.IsExists(configFile) {
		if loadErr = c.loadFromJson(configFile); loadErr != nil {
			return true, loadErr
		}
		return true, nil
	}

	return false, errors.New(fmt.Sprintf("configuration file not found, basePath is %s", basePath))
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

// 从环境变量中获取配置文件profile，必须提供
func (c *ApplicationConf) fileProfileFromEnvVar() (string, error) {
	configPathProfile := os.Getenv(ApplicationEnvProfileVar)
	if configPathProfile != "" {
		return configPathProfile, nil
	}
	return "", errors.New(fmt.Sprintf("not found env %s ", ApplicationEnvProfileVar))
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
	conf := &EnvironmentConf{}
	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return err
	}
	color.Printf("<light_green>global config:</> %+v\n", c.EnvironmentConf)
	return nil
}

func (c *ApplicationConf) loadFromJson(path string) error {
	color.Println("<light_green>load config from:</>", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	conf := &EnvironmentConf{}
	err = json.Unmarshal(data, conf)
	if err != nil {
		return err
	}
	color.Printf("<light_green>global config:</> %+v\n", c.EnvironmentConf)
	return nil
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
	Database Database `json:"database" yaml:"database"`
	Server   Server   `json:"server" yaml:"server"`
	Logger   Logger   `json:"logger" yaml:"logger"`
}
