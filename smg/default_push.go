package smg

import (
	"io/ioutil"
	"path"
	"smg/config"
	"smg/core/tools"
	"smg/registry"

	"github.com/urfave/cli/v2"
)

func defaultPush() *cli.Command {
	return &cli.Command{
		Name:      "push",
		Aliases:   []string{"push"},
		Usage:     "push smg.yaml to remote",
		ArgsUsage: "[push file dir]",
		Action:    actionRegistryPush(),
	}
}

func actionRegistryPush() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		rp, err := registry.AddRemoteProvider(config.Config.Registry)
		if err != nil {
			return err
		}
		args0 := ctx.Args().First()
		if args0 == "" {
			cli.ShowSubcommandHelp(ctx)
			return nil
		}
		b, err := ioutil.ReadFile(tools.AbsDir(args0))
		if err != nil {
			return err
		}
		smg, err := UnmarshalSmg(b, args0)
		if err != nil {
			return err
		}
		err = rp.Set(path.Join("public", path.Base(smg.Path)), b)
		if err != nil {
			return err
		}
		return nil
	}
}
