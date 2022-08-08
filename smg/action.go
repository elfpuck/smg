package smg

import (
	"smg/core/logger"
	"smg/core/tools"

	"github.com/urfave/cli/v2"
)

func Action(ca *CommandAction, smg *Smg) cli.ActionFunc {
	if ca == nil {
		return nil
	}
	return func(ctx *cli.Context) error {
		logger.Debug("action:")
		logger.Debug("app path: ", smg.Path)
		logger.Debug("app name: ", ctx.App.Name)
		logger.Debug("app usage: ", ctx.App.Usage)
		logger.Debug("command name: ", ctx.Command.Name)
		logger.Debug("command usage: ", ctx.Command.Usage)
		logger.Info("action type: ", ca.Type)

		logger.Info("args: ", ctx.Args().Slice())
		ca.Variables = tools.MergeH(ctx.App.Metadata, smg.Variables, tools.H{
			"args":  ctx.Args().Slice(),
			"flags": getFlags(ctx),
		})

		switch ca.Type {
		case "http":
			if ca.Http == nil {
				logger.CommonFatal("yaml Lost action.http config")
			}
			return actionToolsHttp(ctx, ca.Http, ca.Variables)
		case "exec":
			if ca.Exec == nil {
				logger.CommonFatal("yaml Lost action.exec config")
			}
			return actionToolsExec(ctx, ca.Exec, ca.Variables)
		case "mysql":
			if ca.Mysql == nil {
				logger.CommonFatal("yaml Lost action.mysql config")
			}
			return actionToolsMysql(ctx, ca.Mysql, ca.Variables)
		case "consul":
			if ca.Exec == nil {
				logger.CommonFatal("yaml Lost action.exec config")
			}
			return actionToolsConsul(ctx, ca.Consul, ca.Variables)
		default:
			cli.ShowAppHelp(ctx)
			return nil
		}
	}
}

func getFlags(ctx *cli.Context) tools.H {
	flags := tools.H{}
	flagKey := []string{}

	// command flags
	for _, v := range ctx.Command.Flags {
		flagKey = append(flagKey, v.Names()...)
	}

	// app flags
	for _, c := range ctx.Lineage() {
		if c.Command == nil {
			continue
		}
		for _, v := range c.App.Flags {
			flagKey = append(flagKey, v.Names()...)
		}
	}

	for _, k := range flagKey {
		flags[k] = ctx.Value(k)
	}

	return flags
}

type CommandAction struct {
	Type      string
	Variables tools.H
	Output    string
	Exec      *execConfig   `yaml:"exec"`
	Http      *httpConfig   `yaml:"http"`
	Mysql     *mysqlConfig  `yaml:"mysql"`
	Consul    *consulConfig `yaml:"consul"`
}
