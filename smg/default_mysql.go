package smg

import (
	"database/sql"
	"smg/core/logger"
	"smg/core/sqltocsv"
	"smg/core/tools"
	"strings"

	"github.com/go-sql-driver/mysql"

	"github.com/urfave/cli/v2"
)

func defaultMysql() *cli.Command {
	reqMysqlConfig := &mysqlConfig{}
	return &cli.Command{
		Name:      "mysql",
		Aliases:   []string{},
		Usage:     "mysql request",
		ArgsUsage: "[sql...]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "dsn",
				Aliases:     []string{},
				Usage:       "dataSourceName",
				Value:       "",
				Destination: &reqMysqlConfig.Dsn,
			},
			&cli.StringFlag{
				Name:        "net",
				Aliases:     []string{},
				Usage:       "net set [tcp、unix]",
				Value:       "tcp",
				Destination: &reqMysqlConfig.Net,
			},
			&cli.StringFlag{
				Name:        "user",
				Aliases:     []string{"u"},
				Usage:       "user set",
				Value:       "",
				Destination: &reqMysqlConfig.User,
			},
			&cli.StringFlag{
				Name:        "passwd",
				Aliases:     []string{"p"},
				Usage:       "passwd set",
				Value:       "",
				Destination: &reqMysqlConfig.Passwd,
			},
			&cli.StringFlag{
				Name:        "addr",
				Aliases:     []string{"a"},
				Usage:       "addr",
				Value:       "127.0.0.1:3306",
				Destination: &reqMysqlConfig.Addr,
			},
			&cli.StringFlag{
				Name:        "dbname",
				Aliases:     []string{"db"},
				Usage:       "dbname set",
				Value:       "",
				Destination: &reqMysqlConfig.DBName,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Value:       "",
				Usage:       "output file",
				Destination: &reqMysqlConfig.Output,
			},
			&cli.BoolFlag{
				Name:        "allownativepasswords",
				Aliases:     []string{"anp"},
				Usage:       "",
				Value:       true,
				Destination: &reqMysqlConfig.AllowNativePasswords,
			},
		},
		Action: func(ctx *cli.Context) error {
			return actionToolsMysql(ctx, reqMysqlConfig, nil)
		},
	}
}

func actionToolsMysql(ctx *cli.Context, cfg *mysqlConfig, tplData any) error {
	if tplData == nil {
		tplData = ctx.App.Metadata
	}
	if cfg.Query == "" {
		cfg.Query = strings.Join(ctx.Args().Slice(), " ")
	}
	cfg.Query = tools.DrawTpl(tplData, cfg.Query)
	if len(cfg.Query) < 10 {
		cli.ShowSubcommandHelpAndExit(ctx, 0)
	}

	dsn := tools.DrawTpl(tplData, cfg.Dsn)
	if dsn == "" {
		cfg.Net = tools.DrawTpl(tplData, cfg.Net)
		cfg.Addr = tools.DrawTpl(tplData, cfg.Addr)
		cfg.User = tools.DrawTpl(tplData, cfg.User)
		cfg.Passwd = tools.DrawTpl(tplData, cfg.Passwd)
		cfg.DBName = tools.DrawTpl(tplData, cfg.DBName)
		dsn = cfg.FormatDSN()
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		return err
	}
	if strings.ToLower(cfg.Query[:6]) == "select" {
		rows, err := db.Query(cfg.Query)
		if err != nil {
			logger.Fatal(err)
		}
		// output
		res, err := sqltocsv.WriteBytes(rows)
		if err != nil {
			logger.Error(err)
			return err
		}
		return outputFn(&cfg.outputConfig, res)
	} else {
		res, err := db.Exec(cfg.Query)
		if err != nil {
			return err
		}
		if insertId, err := res.LastInsertId(); err == nil {
			logger.CommonInfo("output insertId: ", insertId)
		}
		if affected, err := res.RowsAffected(); err == nil {
			logger.CommonInfo("output affected: ", affected)
		}
	}
	return nil
}

type mysqlConfig struct {
	outputConfig `yaml:",inline"`
	mysql.Config `yaml:",inline"`
	Dsn          string `yaml:"dsn"`
	Query        string `yaml:"query"`
}
