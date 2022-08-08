package smg

import (
	"smg/core/tools"

	"github.com/urfave/cli/v2"
)

func defaultUpgrade() *cli.Command {
	return &cli.Command{
		Name:    "upgrade",
		Aliases: []string{"up"},
		Usage:   "upgrade smg version",
		Action: func(ctx *cli.Context) error {
			if err := tools.DmlExec("go install ./"); err != nil {
				return err
			}
			return nil
		},
	}
}
