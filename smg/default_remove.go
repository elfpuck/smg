package smg

import (
	"smg/config"

	"github.com/urfave/cli/v2"
)

func defaultRemove() *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Aliases:   []string{"rm"},
		Usage:     "remove file from cache local",
		ArgsUsage: "[remove file path...]",
		Action:    actionCacheRemove(),
	}
}

func actionCacheRemove() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		args := ctx.Args().Slice()
		if len(args) == 0 {
			cli.ShowSubcommandHelp(ctx)
			return nil
		}
		for _, v := range args {
			if err := config.CacheRp.Delete(v); err != nil {
				return err
			}
		}
		return nil
	}
}
