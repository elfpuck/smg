package smg

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"smg/config"
	"smg/core/logger"
	"strings"
	"syscall"

	"github.com/urfave/cli/v2"
)

const (
	echoServer = "echoServer.pid"
)

func defaultEchoServer() *cli.Command {
	return &cli.Command{
		Name:      "echoServer",
		Aliases:   []string{"es"},
		Usage:     "response request data",
		ArgsUsage: "[addr]",
		Action:    actionEchoServer(),
	}
}

func actionEchoServer() func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		address := ctx.Args().First()
		if address == "" {
			address = config.Config.Conf.EchoAddress
		}
		if pid, err := config.LocalRp.Get(echoServer); err == nil {
			logger.CommonFatal("echoServer started! pid: " + string(pid))
		}

		engine := Engine{}
		server := &http.Server{Addr: address, Handler: &engine}
		sch := make(chan os.Signal, 1)
		signal.Notify(sch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
		logger.CommonInfo("echoServer Listen: ", address)
		logger.Info("echoServer Pid: ", os.Getpid())
		if err := config.LocalRp.Set(echoServer, []byte(fmt.Sprintf("%v", os.Getpid()))); err != nil {
			logger.Fatal(err)
		}
		go func() {
			for {
				sig := <-sch
				logger.Info("echoServer signal cause stop: ", sig)
				signal.Stop(sch)
				config.LocalRp.Delete(echoServer)
				server.Shutdown(context.TODO())
			}
		}()
		return server.ListenAndServe()
	}
}

type Engine struct {
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	reqBodyByte, _ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	var body interface{}
	if strings.Contains(strings.ToLower(req.Header.Get("Content-Type")), "json") {
		json.Unmarshal(reqBodyByte, &body)
	} else {
		body = string(reqBodyByte)
	}
	json.Unmarshal(reqBodyByte, &body)

	obj := map[string]interface{}{
		"Host":   req.Host,
		"Method": req.Method,
		"Header": req.Header,
		"Router": req.URL.Path,
		"Params": req.URL.Query(),
		"Body":   body,
	}
	encoder := json.NewEncoder(w)

	if err := encoder.Encode(obj); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
