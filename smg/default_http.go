package smg

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"smg/config"
	"smg/core/logger"
	"smg/core/tools"
	"strings"

	"github.com/urfave/cli/v2"
)

func defaultHttp() *cli.Command {
	cfg := &httpConfig{
		Header: map[string][]string{},
	}
	return &cli.Command{
		Name:      "http",
		Aliases:   []string{},
		Usage:     "http request",
		ArgsUsage: "[request url]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "method",
				Aliases:     []string{"X"},
				Value:       "GET",
				Destination: &cfg.Method,
			},
			&cli.StringFlag{
				Name:        "url",
				Aliases:     []string{},
				Value:       "",
				Destination: &cfg.URL,
			},
			&cli.StringFlag{
				Name:        "user-agent",
				Aliases:     []string{"A"},
				Value:       fmt.Sprintf("%s/%s", config.Config.Conf.Name, config.Config.Conf.Version),
				Destination: &cfg.UserAgent,
			},
			&cli.StringSliceFlag{
				Name:        "header",
				Aliases:     []string{"H"},
				Usage:       `example: -H "Accept-Language: en-US"`,
				Destination: &cfg.FlagHeader,
			},
			&cli.StringSliceFlag{
				Name:        "cookie",
				Aliases:     []string{"b"},
				Usage:       `example: -b "token1=v1; token2=v2"`,
				Destination: &cfg.FlagCookie,
			},
			&cli.StringFlag{
				Name:        "content-type",
				Aliases:     []string{"ct"},
				Value:       "",
				Destination: &cfg.ContentType,
			},
			&cli.StringFlag{
				Name:        "referer",
				Aliases:     []string{"e"},
				Value:       "",
				Destination: &cfg.Referer,
			},
			&cli.StringSliceFlag{
				Name:        "form-data",
				Aliases:     []string{"form"},
				Usage:       `example: --form 'file[0]=@"./smg.yaml"'`,
				Destination: &cfg.FlagFormData,
			},
			&cli.StringFlag{
				Name:        "data-raw",
				Aliases:     []string{"dr"},
				Usage:       `example: --data-raw '{"key": "value"}'`,
				Destination: &cfg.DR,
			},
			&cli.StringSliceFlag{
				Name:        "data-urlencode",
				Aliases:     []string{"du"},
				Usage:       `example: --data-urlencode 'key=value'`,
				Destination: &cfg.FlagDU,
			},
			&cli.StringFlag{
				Name:        "resultpath",
				Aliases:     []string{"rp"},
				Value:       "@pretty",
				Usage:       "json response path using gjson",
				Destination: &cfg.ResultPath,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Value:       "",
				Usage:       "output file",
				Destination: &cfg.Output,
			},
		},
		Action: func(ctx *cli.Context) error {
			if cfg.URL == "" {
				cfg.URL = strings.Join(ctx.Args().Slice(), "")
			}
			if cfg.URL == "" {
				cli.ShowSubcommandHelp(ctx)
				return nil
			}
			return actionToolsHttp(ctx, cfg, nil)
		},
	}
}

func actionToolsHttp(ctx *cli.Context, cfg *httpConfig, tplData any) error {
	if tplData == nil {
		tplData = ctx.App.Metadata
	}

	//data append
	cfg.Cookie = append(cfg.Cookie, cfg.FlagCookie.Value()...)
	// cfg.Header = append(cfg.Header, cfg.FlagHeader.Value()...)
	cfg.DU = append(cfg.DU, cfg.FlagDU.Value()...)
	cfg.FormData = append(cfg.FormData, cfg.FlagFormData.Value()...)

	// header merge
	for _, v := range cfg.FlagHeader.Value() {
		v = tools.DrawTpl(tplData, v)
		hv := strings.Split(v, ":")
		if len(hv) != 2 {
			continue
		}
		cfg.Header[hv[0]] = append(cfg.Header[hv[0]], hv[1])
	}

	var payload io.Reader

	// data-raw
	if cfg.DR != "" {
		payload = strings.NewReader(tools.DrawTpl(tplData, cfg.DR))
	} else if len(cfg.DU) != 0 {
		// data-urlencode
		if cfg.ContentType == "" {
			cfg.ContentType = "application/x-www-form-urlencoded"
		}
		tempBody := strings.Join(cfg.DU, `&`)

		logger.Info("data-urlencode: ", tempBody)
		payload = strings.NewReader(tools.DrawTpl(tplData, tempBody))
	} else if len(cfg.DUM) != 0 {
		payload = strings.NewReader(tools.DrawTpl(tplData, cfg.DUM.Encode()))
	} else if len(cfg.FormData) != 0 {
		// form-data
		tempPayload := &bytes.Buffer{}
		writer := multipart.NewWriter(tempPayload)
		for _, v := range cfg.FormData {
			v = tools.DrawTpl(tplData, v)
			fv := strings.Split(v, "@")
			if len(fv) != 2 {
				continue
			}
			filePath := tools.AbsDir(strings.TrimRight(strings.TrimLeft(fv[1], `"`), `"`))
			file, err := os.Open(filePath)
			if err != nil {
				logger.CommonFatal("open form file err: ", err)
			}
			defer file.Close()
			part, err := writer.CreateFormFile(fv[0], filepath.Base(filePath))
			if err != nil {
				logger.CommonFatal("createFormFile err: ", err)
			}
			_, err = io.Copy(part, file)
			if err != nil {
				logger.CommonFatal("io.copy err: ", err)
			}
		}
		err := writer.Close()
		if err != nil {
			logger.CommonFatal("close writer err: ", err)
		}
		payload = tempPayload
		if cfg.ContentType == "" {
			cfg.ContentType = writer.FormDataContentType()
		}
	}

	// request
	req, err := http.NewRequest(strings.ToUpper(tools.DrawTpl(tplData, cfg.Method)), tools.DrawTpl(tplData, cfg.URL), payload)
	if err != nil {
		logger.Fatal(err)
	}

	// header
	for k, v := range cfg.Header {
		switch strings.ToLower(k) {
		case "user-agent":
			if cfg.UserAgent == "" {
				cfg.UserAgent = v[0]
			}
		case "content-type":
			if cfg.ContentType == "" {
				cfg.ContentType = v[0]
			}
		case "referer":
			if cfg.Referer == "" {
				cfg.Referer = v[0]
			}
		case "cookie":
			if cfg.Cookie == nil || len(cfg.Cookie) == 0 {
				cfg.Cookie = v
			}
		default:
			for _, value := range v {
				req.Header.Add(strings.TrimSpace(k), strings.TrimSpace(value))
			}
		}
	}

	// user-agent
	if cfg.UserAgent == "" {
		cfg.UserAgent = fmt.Sprintf("%s/%s", config.Config.Conf.Name, config.Config.Conf.Version)
	}
	req.Header.Add("User-Agent", fmt.Sprintf("%s/%s", config.Config.Conf.Name, config.Config.Conf.Version))

	// Content-Type
	if cfg.ContentType != "" {
		req.Header.Add("Content-Type", tools.DrawTpl(tplData, cfg.ContentType))
	}
	// cookie
	if cfg.Cookie != nil && len(cfg.Cookie) != 0 {
		req.Header.Add("Cookie", tools.DrawTpl(tplData, strings.Join(cfg.Cookie, "; ")))
	}
	// referer
	if cfg.Referer != "" {
		req.Header.Add("Referer", tools.DrawTpl(tplData, strings.TrimSpace(cfg.Referer)))
	}

	//logger
	logger.Debug("req url: ", req.URL)
	logger.Debug("req header: ", req.Header.Clone())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	logger.Info("res status:", resp.Status)
	logger.Debug("res header: ", resp.Header)

	// output
	return outputFn(&cfg.outputConfig, body)
}

type httpConfig struct {
	outputConfig `yaml:",inline"`
	Method       string              `yaml:"method"`
	FlagHeader   cli.StringSlice     `yaml:"-"`
	Header       map[string][]string `yaml:"header"`
	FlagCookie   cli.StringSlice     `yaml:"-"`
	Cookie       []string            `yaml:"-"`
	Referer      string              `yaml:"-"`
	UserAgent    string              `yaml:"-"`
	DR           string              `yaml:"data-raw"`
	FlagDU       cli.StringSlice     `yaml:"-"`
	DU           []string            `yaml:"data-urlencode"`
	DUM          url.Values          `yaml:"data-urlencode-map"`
	FlagFormData cli.StringSlice     `yaml:"-"`
	FormData     []string            `yaml:"form-data"`
	ContentType  string              `yaml:"-"`
	URL          string              `yaml:"url"`
}