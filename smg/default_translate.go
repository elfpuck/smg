package smg

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"smg/config"
	"smg/core/logger"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

func defaultTranslate() *cli.Command {
	reqTranslateConfig := translateConfig{}
	return &cli.Command{
		Name:    "translate",
		Aliases: []string{"fy"},
		Usage:   "翻译",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "from",
				Aliases:     []string{},
				Usage:       "翻译源类型",
				Value:       "AUTO",
				Destination: &reqTranslateConfig.From,
			},
			&cli.StringFlag{
				Name:        "to",
				Aliases:     []string{},
				Usage:       "翻译目标类型",
				Value:       "AUTO",
				Destination: &reqTranslateConfig.To,
			},
			&cli.StringFlag{
				Name:        "resultpath",
				Aliases:     []string{"rp"},
				Value:       "@pretty",
				Usage:       "json response path using gjson",
				Destination: &reqTranslateConfig.ResultPath,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Value:       "",
				Usage:       "output file",
				Destination: &reqTranslateConfig.Output,
			},
		},
		Action: func(ctx *cli.Context) error {
			args := ctx.Args().Slice()
			if len(args) == 0 {
				cli.ShowSubcommandHelp(ctx)
				return nil
			}

			ydConfig := config.Config.Translate.YD
			//有道OpenAPI
			if ydConfig.AppKey != "" && ydConfig.AppSecret != "" {
				ydConfig.From = reqTranslateConfig.From
				ydConfig.To = reqTranslateConfig.To
				ydConfig.Query = strings.Join(args, " ")

				res, err := YdOpenApi(ydConfig)
				if err != nil {
					logger.Fatal(err)
				}
				if !ctx.IsSet("resultpath") {
					reqTranslateConfig.outputConfig.ResultPath = "web|@pretty"
				}
				return outputFn(&reqTranslateConfig.outputConfig, res)
			} else {
				// 有道 FreeAPI
				ydFreeConfig := config.Config.Translate.YDFree

				ydFreeConfig.From = reqTranslateConfig.From
				ydFreeConfig.To = reqTranslateConfig.To
				ydFreeConfig.Query = strings.Join(args, " ")

				res, err := YdFreeApi(ydFreeConfig)
				if err != nil {
					logger.Fatal(err)
				}
				if !ctx.IsSet("resultpath") {
					reqTranslateConfig.outputConfig.ResultPath = "translateResult.0.0.tgt"
				}
				return outputFn(&reqTranslateConfig.outputConfig, res)
			}
		},
	}
}

func YdFreeApi(cfg *config.TranslateYDFree) ([]byte, error) {
	if cfg.Url == "" {
		cfg.Url = "https://fanyi.youdao.com/translate_o?smartresult=dict&smartresult=rule"
	}
	if cfg.Referer == "" {
		cfg.Referer = "https://fanyi.youdao.com/"
	}
	if cfg.Cookie == "" {
		cfg.Cookie = "OUTFOX_SEARCH_USER_ID=-2138848423@10.110.96.154"
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36"
	}
	if cfg.From == "" {
		cfg.From = "auto"
	}
	if cfg.To == "" {
		cfg.To = "auto"
	}
	salt, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	query := url.Values{
		"from":        []string{cfg.From},
		"to":          []string{cfg.To},
		"salt":        []string{salt.String()},
		"smartresult": []string{"dict"},
		"client":      []string{"fanyideskweb"},
		"version":     []string{"2.1"},
		"doctype":     []string{"json"},
		"keyfrom":     []string{"fanyi.web"},
		"action":      []string{"FY_BY_CLICKBUTTION"},
		"i":           []string{cfg.Query},
	}
	query.Add("sign", fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%s%s]BjuETDhU)zqSxf-=B#7m", query.Get("client"), query.Get("i"), query.Get("salt"))))))

	req, err := http.NewRequest("POST", cfg.Url, strings.NewReader(query.Encode()))
	if err != nil {
		return nil, err
	}
	req.Method = "POST"
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Referer", cfg.Referer)
	req.Header.Add("User-Agent", cfg.UserAgent)
	req.Header.Add("Cookie", cfg.Cookie)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

func YdOpenApi(cfg *config.TranslateYD) ([]byte, error) {
	if cfg.Url == "" {
		cfg.Url = "https://openapi.youdao.com/api"
	}
	if cfg.From == "" {
		cfg.From = "auto"
	}
	if cfg.To == "" {
		cfg.To = "auto"
	}
	salt, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	query := url.Values{
		"from":     []string{cfg.From},
		"to":       []string{cfg.To},
		"salt":     []string{salt.String()},
		"curtime":  []string{fmt.Sprintf("%d", time.Now().Unix())},
		"q":        []string{cfg.Query},
		"appKey":   []string{cfg.AppKey},
		"signType": []string{"v3"},
	}
	signInput := ""
	if len(query.Get("q")) <= 20 {
		signInput = query.Get("q")
	} else {
		signInput = fmt.Sprintf("%s%d%s", query.Get("q")[:10], len(query.Get("q")), query.Get("q")[len(query.Get("q"))-10:])
	}
	query["sign"] = []string{ydSign(fmt.Sprintf("%s%s%s%s%s", query.Get("appKey"), signInput, query.Get("salt"), query.Get("curtime"), cfg.AppSecret))}

	req, err := http.NewRequest("POST", cfg.Url, strings.NewReader(query.Encode()))
	if err != nil {
		return nil, err
	}
	req.Method = "POST"
	req.Header.Add("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}

type translateConfig struct {
	outputConfig `yaml:",inline"`
	From         string
	To           string
}

func ydSign(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}
