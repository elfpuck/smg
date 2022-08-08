package smg

import (
	"smg/config"
	"smg/registry"

	"github.com/urfave/cli/v2"
)

func defaultPull() *cli.Command {
	return &cli.Command{
		Name:      "pull",
		Aliases:   []string{"pull"},
		Usage:     "pull smg.yaml from remote",
		ArgsUsage: "[pull file ...]",
		Action:    actionRegistryPull(),
	}
}

func actionRegistryPull() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		rp, err := registry.AddRemoteProvider(config.Config.Registry)
		if err != nil {
			return err
		}
		args := ctx.Args().Slice()
		if len(args) == 0 {
			cli.ShowSubcommandHelp(ctx)
			return nil
		}
		for _, v := range args {
			b, err := rp.Get(v)
			if err != nil {
				return err
			}
			if err := config.CacheRp.Set(v, b); err != nil {
				return err
			}
		}

		return nil
	}
}
