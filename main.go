package main

import (
	"os"
	"path"
	"smg/config"
	"smg/core/logger"
	"smg/core/tools"
	"smg/smg"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {
	// log file
	f, err := os.OpenFile(path.Join(tools.AbsDir(path.Join(config.ConfigDir, config.LogDir)), config.LogPath), os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	logger.Init(&logger.Config{
		Output: f,
		Level:  logger.ParseLevel(config.Config.Logger.Level),
	})
	logger.Println()
	logger.Info("run : ", "smg ", strings.Join(os.Args[1:], " "))

	// 初始化
	app := cli.App{
		Name:        config.Config.Conf.Name,
		Usage:       config.Config.Conf.Usage,
		Description: config.Config.Conf.Desc,
		Authors:     []*cli.Author{},
		Flags:       []cli.Flag{},
		Commands:    []*cli.Command{},
		Version:     config.Config.Conf.Version,
		Suggest:     true,
	}

	// parse Flags
	mflag, mValue := tools.ModulesFlag(path.Join(config.ConfigDir, config.CacheDir))
	app.Flags = append(app.Flags, &mflag)

	// load User Modules
	loadCmd, err := smg.LoadModules(mValue)
	if err != nil {
		logger.Fatal("load Modules Err: ", err)
	}
	app.Commands = append(app.Commands, loadCmd...)

	app.Commands = append(app.Commands, smg.SystemCommand())
	app.Metadata = config.Config.Variables

	// run
	if err := app.Run(os.Args); err != nil {
		logger.CommonFatal("smg err: ", err)
	}
}
