package gateway

import (
	"fmt"
	"github.com/gookit/color"
	"gopkg.in/yaml.v3"
	"os"
)

type ApiOptions struct {
	Addr          string `toml:"addr"`
	Proxy         string `toml:"proxy"`
	ApiKey        string `toml:"api-key"`
	SecretKey     string `toml:"secret-key"`
	AutoReconnect bool   `toml:"auto-reconnect"`
	DebugMode     bool   `toml:"debug-mode"`
}

type apiOptions = func(o *ApiOptions)

func NewApiOptions(opts ...apiOptions) *ApiOptions {
	options := &ApiOptions{}
	for _, opt := range opts {
		opt(options)
	}
	return options
}

func WithConfig(filepath string) apiOptions {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return func(o *ApiOptions) {
			color.Redln(fmt.Sprintf("api options WithConfig fail. %s", err))
		}
	}
	c := &ApiOptions{}
	if err = yaml.Unmarshal(data, c); err != nil {
		color.Redln(fmt.Sprintf("api options WithConfig fail. %s", err))
	}
	return func(o *ApiOptions) {
		o.Addr = c.Addr
		o.Proxy = c.Proxy
		o.ApiKey = c.ApiKey
		o.SecretKey = c.SecretKey
		o.DebugMode = c.DebugMode
		o.AutoReconnect = c.AutoReconnect
	}
}

func WithAddr(addr string) apiOptions {
	return func(o *ApiOptions) {
		o.Addr = addr
	}
}

func WithProxy(proxy string) apiOptions {
	return func(o *ApiOptions) {
		o.Proxy = proxy
	}
}

func WithApiKey(apiKey string) apiOptions {
	return func(o *ApiOptions) {
		o.ApiKey = apiKey
	}
}

func WithSecretKey(secretKey string) apiOptions {
	return func(o *ApiOptions) {
		o.SecretKey = secretKey
	}
}

func WithAutoReconnect(autoReconnect bool) apiOptions {
	return func(o *ApiOptions) {
		o.AutoReconnect = autoReconnect
	}
}

func WithDebugMode(debugMode bool) apiOptions {
	return func(o *ApiOptions) {
		o.DebugMode = debugMode
	}
}
