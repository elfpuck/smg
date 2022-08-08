package smg

import (
	"path"
	"smg/config"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func defaultList() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "list cached file",
		Action:  actionCacheList(),
	}
}

func actionCacheList() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		kvs, err := config.CacheRp.List("/")
		if err != nil {
			return err
		}
		t := tableRender([]string{"name", "path", "version", "desc", "commands"})
		for _, v := range kvs {
			if path.Ext(v.Key) != ".yaml" {
				continue
			}
			smg := Smg{
				Path: v.Key,
			}
			if err := yaml.Unmarshal(v.Value, &smg); err != nil {
				continue
			}
			t.AppendRow([]any{smg.Name, smg.Path, smg.Version, smg.Desc, smg.cmdShow(5)})
			t.AppendSeparator()
		}
		t.Render()
		return nil
	}
}
