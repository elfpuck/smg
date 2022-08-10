package smg

import (
	"smg/core/logger"
	"smg/core/tools"
	"sync"

	"github.com/urfave/cli/v2"
)

var RegistryCommands sync.Map

func RenderCommand(subCommand map[string]*Command, smg *Smg, root string) []*cli.Command {
	result := []*cli.Command{}
	for k, v := range subCommand {
		if v.Name == "" {
			v.Name = k
		}
		if root == "" {
			if actual, loaded := RegistryCommands.LoadOrStore(v.Name, smg.Path); loaded {
				logger.CommonFatal("\ncommand:       %s is loaded, please check your modules file!\nloaded file:   %s\nconflict file: %s\n", v.Name, actual, smg.Path)
			}
		}
		command := cli.Command{
			Name:        v.Name,
			Aliases:     v.Aliases,
			Usage:       v.Usage,
			Description: v.Desc,
			ArgsUsage:   v.ArgsUsage,
			Category:    v.Category,
			Flags:       v.ParseFlags(),
			Subcommands: RenderCommand(v.Subcommand, smg, tools.JoinSlash(root, v.Name)),
			Action:      Action(v, smg),
		}

		result = append(result, &command)
	}
	return result
}

func (cmd *Command) ParseFlags() (res []cli.Flag) {
	for k, v := range cmd.Flag {
		res = append(res, &cli.StringFlag{
			Name:       k,
			Aliases:    v.Aliases,
			Usage:      v.Usage,
			Required:   v.Required,
			Value:      v.Value,
			HasBeenSet: v.HasBeenSet,
			EnvVars:    v.EnvVars,
		})
	}
	return res
}

type Command struct {
	Name        string `yaml:"-"`
	Aliases     []string
	Usage       string
	Flag        map[string]CommandFlag
	Desc string	
	ArgsUsage   string `yaml:"argsUsage"`
	ArgsMin     int    `yaml:"argsMin"`
	Category    string `yaml:"category"`
	Action      *CommandAction
	// 自定义参数
	Subcommand map[string]*Command `yaml:"subCommand"`
}

type CommandFlag struct {
	Usage      string
	Aliases    []string
	Required   bool
	Value      string
	HasBeenSet bool     `yaml:"hasBeenSet"`
	EnvVars    []string `yaml:"envVars"`
}
