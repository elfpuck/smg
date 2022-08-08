package smg

import (
	"encoding/json"
	"smg/core/logger"
	"smg/core/tools"
	"smg/registry"
	"strings"

	"github.com/urfave/cli/v2"
)

func defaultConsul() *cli.Command {
	reqConsulConfig := &consulConfig{}
	return &cli.Command{
		Name:    "consul",
		Aliases: []string{},
		Usage:   "consul system",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "address",
				Aliases:     []string{"addr"},
				Usage:       "",
				Value:       "http://127.0.0.1:8500",
				Destination: &reqConsulConfig.Address,
			},
			&cli.StringFlag{
				Name:        "token",
				Aliases:     []string{"t"},
				Usage:       "acl token",
				Value:       "",
				Destination: &reqConsulConfig.Token,
			},
			&cli.StringFlag{
				Name:        "prefix",
				Aliases:     []string{"p"},
				Usage:       "prefix",
				Value:       "",
				Destination: &reqConsulConfig.Prefix,
			},
			&cli.StringFlag{
				Name:        "resultpath",
				Aliases:     []string{"rp"},
				Value:       "@pretty",
				Usage:       "json response path using gjson",
				Destination: &reqConsulConfig.ResultPath,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Value:       "",
				Usage:       "output file",
				Destination: &reqConsulConfig.Output,
			},
		},
		Subcommands: []*cli.Command{
			{
				Name:      "get",
				ArgsUsage: "[key]",
				Action: func(ctx *cli.Context) error {
					reqConsulConfig.Method = "GET"
					reqConsulConfig.Key = ctx.Args().First()
					if reqConsulConfig.Key == "" {
						cli.ShowSubcommandHelp(ctx)
						return nil
					}
					return actionToolsConsul(ctx, reqConsulConfig, nil)
				},
			},
			{
				Name:      "set",
				ArgsUsage: "[key, value]",
				Action: func(ctx *cli.Context) error {
					reqConsulConfig.Method = "SET"
					if ctx.Args().Len() != 2 {
						cli.ShowSubcommandHelp(ctx)
						return nil
					}
					reqConsulConfig.Key = ctx.Args().First()
					reqConsulConfig.Value = ctx.Args().Get(1)

					return actionToolsConsul(ctx, reqConsulConfig, nil)
				},
			},
			{
				Name:      "list",
				ArgsUsage: "[filter prefix ]",
				Action: func(ctx *cli.Context) error {
					reqConsulConfig.Method = "LIST"
					reqConsulConfig.Key = ctx.Args().First()
					return actionToolsConsul(ctx, reqConsulConfig, nil)
				},
			},
			{
				Name:      "delete",
				ArgsUsage: "[key]",
				Action: func(ctx *cli.Context) error {
					reqConsulConfig.Method = "DELETE"
					reqConsulConfig.Key = ctx.Args().First()
					if reqConsulConfig.Key == "" {
						cli.ShowSubcommandHelp(ctx)
						return nil
					}
					return actionToolsConsul(ctx, reqConsulConfig, nil)
				},
			},
		},
	}
}

func actionToolsConsul(ctx *cli.Context, cfg *consulConfig, tplData any) error {
	if tplData == nil {
		tplData = ctx.App.Metadata
	}
	var res []byte
	var err error
	rp, err := registry.AddRemoteProvider(&registry.Config{
		Provider: "consul",
		Prefix:   tools.DrawTpl(tplData, cfg.Prefix),
		Address:  tools.DrawTpl(tplData, cfg.Address),
		Token:    tools.DrawTpl(tplData, cfg.Token),
	})
	if err != nil {
		logger.Error("consul rp Err: ", err)
		return err
	}
	cfg.Method = tools.DrawTpl(tplData, cfg.Method)
	cfg.Key = tools.DrawTpl(tplData, cfg.Key)
	cfg.Value = tools.DrawTpl(tplData, cfg.Value)

	switch strings.ToUpper(cfg.Method) {
	case "GET":
		res, err = rp.Get(cfg.Key)
		if err != nil {
			logger.CommonFatal("run consul Error: ", err)
		}
	case "SET":
		err = rp.Set(cfg.Key, []byte(cfg.Value))
	case "LIST":
		tempRes, tempErr := rp.List(cfg.Key)
		if tempErr != nil {
			err = tempErr
		} else {
			res, err = json.Marshal(tempRes)
		}
	case "DELETE":
		err = rp.Delete(cfg.Key)
	default:
		logger.CommonFatal("error method [ GET、SET、LIST、DELETE ]: ", cfg.Method)
	}
	if err != nil {
		logger.Error("run consul Error: ", err)
		return err
	}
	// output
	return outputFn(&cfg.outputConfig, res)
}

type consulConfig struct {
	outputConfig
	Address string `yaml:"address"`
	Prefix  string `yaml:"prefix"`
	Key     string `yaml:"key"`
	Token   string `yaml:"token"`
	Method  string `yaml:"method"`
	Value   string `yaml:"value"`
}
