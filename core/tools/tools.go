package tools

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v2"
)

type H map[string]any

func JoinSlash(s ...string) string {
	if len(s) > 0 && s[0] == "" {
		s = s[1:]
	}
	return strings.Join(s, "/")
}
func SplitSlash(s string) []string {
	return strings.Split(s, "/")
}

func DmlExec(script string) error {
	cmd := exec.Command("bash", "-c", script)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func AbsDir(fp string) string {
	if path.IsAbs(fp) {
		return fp
	}
	homeFlag := "~"
	if strings.HasPrefix(fp, homeFlag) {
		homeUrl, _ := os.UserHomeDir()
		return path.Join(homeUrl, strings.TrimLeft(fp, homeFlag))
	} else {
		absDp, _ := filepath.Abs(fp)
		return absDp
	}
}

// 相对路径请基于程序 os.GetWd()
func EnsureDir(fp string) (absDp string, err error) {
	absDp = AbsDir(fp)
	if _, err := os.Stat(absDp); errors.Is(err, fs.ErrNotExist) {
		err := os.MkdirAll(absDp, os.ModePerm)
		if err != nil {
			return absDp, err
		}
	}
	return absDp, nil
}

func ReadOrCreateFile(f string, b []byte, update bool) ([]byte, error) {
	dir := path.Dir(f)
	absDp, err := EnsureDir(dir)
	if err != nil {
		return nil, err
	}
	absFp := path.Join(absDp, path.Base(f))
	fpByte, fpErr := ioutil.ReadFile(absFp)
	if fpErr != nil && !errors.Is(fpErr, fs.ErrNotExist) {
		return nil, err
	}
	file, err := os.OpenFile(absFp, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if update || errors.Is(fpErr, fs.ErrNotExist) {
		if _, err := file.Write(b); err != nil {
			return nil, err
		}
		return b, err
	}
	return fpByte, nil
}

func DrawTpl(data any, tpl string) (resData string) {
	if tpl == "" {
		return tpl
	}
	defer func() {
		if err := recover(); err != nil {
			resData = tpl
			fmt.Print(err)
		}
	}()
	t := template.Must(template.New("").Funcs(
		template.FuncMap{
			"unescaped": func(str string) template.HTML { return template.HTML(str) },
			"join": func(args []string, sep string) string {
				return strings.Join(args, sep)
			},
		}).Parse(tpl))
	res := new(bytes.Buffer)
	if err := t.Execute(res, data); err != nil {
		return tpl
	}
	return res.String()
}

func DrawTplMulti(data any, tpl ...string) []string {
	res := []string{}
	for _, v := range tpl {
		res = append(res, DrawTpl(data, v))
	}
	return res
}

func MergeH(args ...H) H {
	res := H{}
	for _, h := range args {
		for k, v := range h {
			res[k] = v
		}
	}
	return res
}

func Md5(args ...string) string {
	str := fmt.Sprint(args)
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func ModulesFlag(cacheDir string) (cli.StringFlag, string) {
	mFlag := cli.StringFlag{
		Name:    "modules",
		Aliases: []string{"m"},
		Usage:   "load modules path",
		Value:   cacheDir,
	}
	var mFlagValue string
	flag.StringVar(&mFlagValue, mFlag.Name, mFlag.Value, "")
	for _, v := range mFlag.Aliases {
		flag.StringVar(&mFlagValue, v, mFlag.Value, "")
	}
	flag.Parse()
	return mFlag, mFlagValue
}
