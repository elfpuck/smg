package smg

import (
	"errors"
	"path"
	"smg/config"
	"smg/registry"

	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

func defaultImages() *cli.Command {
	return &cli.Command{
		Name:    "images",
		Aliases: []string{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "prefix",
				Aliases: []string{"p"},
				Usage:   "images prefix",
				Value:   "/",
			},
		},
		Usage:  "list remote smg.yaml files",
		Action: actionRegistryList(),
	}
}

func actionRegistryList() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		registryCfg := config.Config.Registry
		if registryCfg == nil {
			return errors.New("lost .yaml file path args")
		}
		rp, err := registry.AddRemoteProvider(config.Config.Registry)
		if err != nil {
			return err
		}
		kvs, err := rp.List(ctx.String("prefix"))
		if err != nil {
			return err
		}
		t := tableRender([]string{"name", "path", "version", "commands", "desc"})
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
			if smg.Command == nil {
				continue
			}
			t.AppendRow([]any{smg.Name, smg.Path, smg.Version, smg.cmdShow(config.Config.Conf.SubCommandLen), smg.Desc})
			t.AppendSeparator()
		}
		t.SetAllowedRowLength(config.Config.Conf.TableLen)
		t.Render()
		return nil
	}
}
