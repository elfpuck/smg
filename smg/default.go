package smg

import (
	"os"
	"smg/core/logger"
	"smg/core/tools"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v2"
)

func SystemCommand() *cli.Command {
	command := &cli.Command{
		Name:        "system",
		Aliases:     []string{"sys"},
		Usage:       "smg system command",
		Description: "smg self command, contains config, module download etc",
		Subcommands: []*cli.Command{
			defaultMysql(),
			defaultHttp(),
			defaultConsul(),
			defaultExec(),
			defaultList(),
			defaultPull(),
			defaultPush(),
			defaultImages(),
			defaultRemove(),
			defaultUpgrade(),
			defaultEchoServer(),
		},
	}
	return command
}

func tableRender(header []string) table.Writer {
	h := table.Row{}
	for _, v := range header {
		h = append(h, text.FgCyan.Sprint(v))
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(h)
	t.SetStyle(table.StyleBold)
	return t
}

func outputFn(cfg *outputConfig, res []byte) (err error) {
	logger.Debug("output: ", cfg.Output)
	logger.Debug("output path: ", cfg.ResultPath)
	if cfg.ResultPath == "" {
		cfg.ResultPath = "@pretty"
	}
	if len(res) == 0 {
		return
	}
	if cfg.Output != "" {
		fileName := tools.AbsDir(cfg.Output)
		if _, err := tools.ReadOrCreateFile(fileName, res, true); err != nil {
			return err
		}
		logger.CommonInfo("output file: ", fileName)
	} else {
		var outputBody string
		if gjson.ValidBytes(res) {
			logger.Info("json")
			outputBody = gjson.GetBytes(res, cfg.ResultPath).Raw
		}
		if outputBody == "" {
			outputBody = string(res)
		}
		logger.CommonInfo("res body: ", "\n", outputBody)
	}
	return
}

type outputConfig struct {
	Output     string `yaml:"output"`
	ResultPath string `yaml:"resultPath"`
}
