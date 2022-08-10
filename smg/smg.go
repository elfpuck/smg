package smg

import (
	"path"
	"smg/config"
	"smg/core/logger"
	"smg/core/tools"
	"smg/registry"
	"strings"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func LoadModules(filename string) ([]*cli.Command, error) {
	logger.Debug("load modules: ", filename)
	result := []*cli.Command{}
	rp := config.LocalRp
	if filename != path.Join(config.ConfigDir, config.CacheDir) {
		localRp, err := registry.AddRemoteProvider(&registry.Config{
			Provider: "local",
			Prefix:   filename,
		})
		if err != nil {
			return result, err
		}
		rp = localRp
	}

	kvpairs, err := rp.List("/")
	if err != nil {
		return result, err
	}
	for _, v := range kvpairs {
		if path.Ext(v.Key) != ".yaml" {
			continue
		}
		smg := Smg{
			Path: v.Key,
		}
		if err := yaml.Unmarshal(v.Value, &smg); err != nil {
			return nil, err
		}
		if smg.Name == "" {
			smg.Name = smg.getId()
		}

		command := RenderCommand(smg.Command, &smg, "")

		result = append(result, command...)
	}
	return result, nil
}

func (smg *Smg) cmdShow(n int) string {
	cmds := []string{}
	for k, v := range smg.Command {
		temp := k + ":("
		tempArr := []string{}
		for k1 := range v.Subcommand {
			if len(tempArr) >= n {
				tempArr = append(tempArr, "...")
				break
			}
			tempArr = append(tempArr, k1)
		}
		temp += strings.Join(tempArr, ",")
		temp += ")"
		cmds = append(cmds, temp)
	}
	return strings.Join(cmds, "\n")
}

func (smg *Smg) getId() string {
	if smg.Id != "" {
		return smg.Id
	}
	return tools.Md5(smg.Name, smg.Version, smg.Path)
}

func UnmarshalSmg(b []byte, path string) (Smg, error) {
	smg := Smg{
		Path: path,
	}
	if err := yaml.Unmarshal(b, &smg); err != nil {
		return smg, err
	}
	return smg, nil
}

type Smg struct {
	Variables tools.H `yaml:"variables"`
	Id        string  `yaml:"-"`
	Name      string
	Desc      string
	Path      string `yaml:"-"`
	Version   string
	Command   map[string]*Command `yaml:"command"`
}
