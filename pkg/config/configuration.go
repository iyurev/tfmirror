package config

import (
	"fmt"
	"github.com/iyurev/tfmirror/pkg/types"
	"github.com/spf13/viper"
	"os"
)

const (
	AllAvailable = "all"
)

var (
	AllAvailableList = []string{AllAvailable}
)

type ProviderConf struct {
	Source    string
	Platforms []types.Platform
	Versions  []string
}

func (p *ProviderConf) HasVersion(version string) bool {
	for _, v := range p.Versions {
		if v == version {
			return true
		}
	}
	return false
}

func (p *ProviderConf) SetDefaults() {
	if p.Versions == nil {
		p.Versions = make([]string, 0)
	}
	if p.Platforms == nil {
		p.Platforms = make([]types.Platform, 0)
	}
}

func (p *ProviderConf) DownloadAllVersions() bool {
	return len(p.Versions) == 0
}
func (p *ProviderConf) DownloadAllPlatforms() bool {
	return len(p.Platforms) == 0
}

type ClientConf struct {
	TimeOut  int    `mapstructure:"timeOut"`
	WorkDir  string `mapstructure:"workDir"`
	LogLevel string `mapstructure:"logLevel"`
}

func (c *ClientConf) SetDefaults() error {
	if c.LogLevel == "" {
		c.LogLevel = "info"
	}
	if c.WorkDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		c.WorkDir = fmt.Sprintf("%s/workdir", cwd)
	}
	if c.TimeOut == 0 {
		c.TimeOut = 5
	}
	return nil
}

type Conf struct {
	Providers map[string]ProviderConf `mapstructure:"providers"`
	Client    *ClientConf             `mapstructure:"client"`
}

func NewConfig() (*Conf, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	conf := &Conf{}
	if err := viper.Unmarshal(conf); err != nil {
		return nil, err
	}
	if err := conf.Client.SetDefaults(); err != nil {
		return nil, err
	}
	for k, _ := range conf.Providers {
		providerConf := conf.Providers[k]
		providerConf.SetDefaults()
		conf.Providers[k] = providerConf
	}

	return conf, nil
}
