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
	LogPath           = "access.log"
	OutputDir         = "output"
	ConfigPath        = "config.yaml"
	defaultConfigPath = "default.config.yaml"
)

var (
	Config = SmgConfig{
		Logger: &SmgConfigLog{},
		Conf: &SmgConfigConf{
			EchoAddress: "127.0.0.1:11111",
		},
		Translate: &SmgConfigTranslate{
			YD:     &TranslateYD{},
			YDFree: &TranslateYDFree{},
		},
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
	Config.Conf.Version = "0.0.3"

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
	Variables tools.H             `yaml:"variables"`
	Conf      *SmgConfigConf      `yaml:"conf"`
	Registry  *registry.Config    `yaml:"registry"`
	Logger    *SmgConfigLog       `yaml:"logger"`
	Translate *SmgConfigTranslate `yaml:"translate"`
}

type SmgConfigTranslate struct {
	YD     *TranslateYD     `yaml:"yd"`
	YDFree *TranslateYDFree `yaml:"ydFree"`
}

type TranslateYD struct {
	Query     string
	Url       string `yaml:"url"`
	From      string `yaml:"from"`
	To        string `yaml:"to"`
	AppKey    string `yaml:"appKey"`
	AppSecret string `yaml:"appSecret"`
}
type TranslateYDFree struct {
	Query     string
	Url       string `yaml:"url"`
	From      string `yaml:"from"`
	To        string `yaml:"to"`
	Referer   string `yaml:"referer"`
	UserAgent string `yaml:"userAgent"`
	Cookie    string `yaml:"cookie"`
}

type SmgConfigConf struct {
	Name          string `yaml:"-"`
	Usage         string
	Desc          string
	Version       string `yaml:"-"`
	EchoAddress   string `yaml:"echoAddress"`
	TableLen      int    `yaml:"tableLen"`
	SubCommandLen int    `yaml:"subCommandLen"`
}

type SmgConfigLog struct {
	Level string
}
