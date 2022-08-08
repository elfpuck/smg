package config

import (
	"embed"
	"path"
	"smg/core/tools"
	"smg/registry"

	"gopkg.in/yaml.v3"
)

//go:embed default.config.yaml
var f embed.FS

var LocalRp registry.RemoteProvider
var CacheRp registry.RemoteProvider

const (
	ConfigDir         = "~/.smg"
	CacheDir          = "cache"
	LogDir            = "log"
	OutputDir         = "output"
	ConfigPath        = "config.yaml"
	defaultConfigPath = "default.config.yaml"
)

var (
	Config = SmgConfig{
		Logger: &SmgConfigLog{},
		Conf:   &SmgConfigConf{},
	}
)

func init() {
	if _, err := tools.EnsureDir(path.Join(ConfigDir, LogDir)); err != nil {
		panic(err)
	}
	if _, err := tools.EnsureDir(path.Join(ConfigDir, CacheDir)); err != nil {
		panic(err)
	}
	if _, err := tools.EnsureDir(path.Join(ConfigDir, OutputDir)); err != nil {
		panic(err)
	}
	baseCfgByte, err := f.ReadFile(defaultConfigPath)
	if err != nil {
		panic(err)
	}

	b, err := tools.ReadOrCreateFile(path.Join(ConfigDir, ConfigPath), baseCfgByte, false)
	if err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal(b, &Config); err != nil {
		panic(err)
	}

	Config.Conf.Name = "smg"
	Config.Conf.Version = "0.0.1"

	rp, err := registry.AddRemoteProvider(&registry.Config{
		Provider: "local",
		Prefix:   ConfigDir,
	})
	if err != nil {
		panic(err)
	}
	LocalRp = rp
	cacheRp, err := registry.AddRemoteProvider(&registry.Config{
		Provider: "local",
		Prefix:   path.Join(ConfigDir, CacheDir),
	})
	if err != nil {
		panic(err)
	}
	CacheRp = cacheRp
}

type SmgConfig struct {
	Variables tools.H          `yaml:"variables"`
	Conf      *SmgConfigConf   `yaml:"conf"`
	Registry  *registry.Config `yaml:"registry"`
	Logger    *SmgConfigLog    `yaml:"logger"`
}

type SmgConfigConf struct {
	Name        string
	Usage       string
	Description string
	Version     string
}

type SmgConfigLog struct {
	Level string
}
