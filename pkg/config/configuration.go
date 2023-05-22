package config

import (
	"github.com/iyurev/tfmirror/pkg/types"
	"github.com/spf13/viper"
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
	TimeOut int
	WorkDir string `mapstructure:"workDir"`
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

	for k, _ := range conf.Providers {
		providerConf := conf.Providers[k]
		providerConf.SetDefaults()
		conf.Providers[k] = providerConf
	}

	return conf, nil
}
