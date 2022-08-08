package smg

import (
	"os"
	"smg/core/logger"
	"smg/core/tools"
	"strings"

	"github.com/urfave/cli/v2"
)

func defaultExec() *cli.Command {
	reqExecConfig := &execConfig{}
	return &cli.Command{
		Name:      "exec",
		Aliases:   []string{},
		Usage:     "exec command",
		ArgsUsage: "[command, args...]",
		Flags:     []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			args := ctx.Args().Slice()
			if len(os.Args) == 0 {
				cli.ShowSubcommandHelp(ctx)
				return nil
			}
			reqExecConfig.Script = append(reqExecConfig.Script, strings.Join(args, " "))
			return actionToolsExec(ctx, reqExecConfig, nil)
		},
	}
}

func actionToolsExec(ctx *cli.Context, cfg *execConfig, tplData any) error {
	if tplData == nil {
		tplData = ctx.App.Metadata
	}
	for _, v := range tools.DrawTplMulti(tplData, cfg.Script...) {
		logger.Info("run exec: ", v)
		if err := tools.DmlExec(v); err != nil {
			return err
		}
	}
	return nil
}

type execConfig struct {
	Script []string `yaml:"script"`
}
