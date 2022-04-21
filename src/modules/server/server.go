package main

import (
	"Distributed-System-Awareness-Platform/src/models"
	"Distributed-System-Awareness-Platform/src/modules/server/config"
	"context"
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	_ "github.com/go-sql-driver/mysql"
	"github.com/oklog/run"
	"github.com/prometheus/common/promlog"
	prometheuslogflag "github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var (
	// 命令行解析
	app = kingpin.New(filepath.Base(os.Args[0]), "The open-devops-server")
	// 指定配置文件
	configFile = app.Flag("config.file", "open-devops-server configuration file path").Short('c').Default("open-devops-server.yml").String()
)

func main() {

	app.Version(version.Print("open-devops-server"))

	app.HelpFlag.Short('h')

	prometheusLogConfig := promlog.Config{}

	prometheuslogflag.AddFlags(app, &prometheusLogConfig)
	// 强制解析
	kingpin.MustParse(app.Parse(os.Args[1:]))
	fmt.Println(*configFile)

	// 设置logger
	var logger log.Logger
	logger = func(config *promlog.Config) log.Logger {
		var (
			l  log.Logger
			le level.Option
		)
		if config.Format.String() == "logfmt" {
			l = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
		} else {
			l = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
		}

		switch config.Level.String() {
		case "debug":
			le = level.AllowDebug()
		case "info":
			le = level.AllowInfo()
		case "warn":
			le = level.AllowWarn()
		case "error":
			le = level.AllowError()
		}
		l = level.NewFilter(l, le)
		l = log.With(l, "ts", log.TimestampFormat(
			func() time.Time { return time.Now().Local() },
			"2006-01-02T15:04:05.000Z07:00",
		), "caller", log.DefaultCaller)
		return l
	}(&prometheusLogConfig)
	level.Info(logger).Log("msg", "using config.file", "file.path", *configFile)

	//读取配置文件
	sConfig, err := config.LoadFile(*configFile)
	if err != nil {
		level.Error(logger).Log("msg", "config.LoadFile Error,Exiting ...", "error", err)
		return
	}
	level.Info(logger).Log("msg", "load.config.success", "file.path", *configFile, "content.mysql.num", len(sConfig.MysqlS))

	models.InitMySQL(sConfig.MysqlS)
	level.Info(logger).Log("msg", "load.mysql.success", "db.num", len(models.DB))
	// TODO 本地测试node的添加，后续需要删掉
	//models.StreePathAddTest(logger)
	//models.StreePathQueryTest3(logger)

	// 编排开始
	var g run.Group
	ctxAll, cancelAll := context.WithCancel(context.Background())
	fmt.Println(ctxAll)
	{

		// 处理信号退出的handler
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		cancelC := make(chan struct{})
		g.Add(
			func() error {
				//通过select监听channel 如果有数据就说明发送了ctrl c
				select {
				case <-term:
					level.Warn(logger).Log("msg", "Receive SIGTERM ,exiting gracefully....")
					cancelAll()
					return nil
				case <-cancelC:
					level.Warn(logger).Log("msg", "other cancel exiting")
					return nil
				}
			},
			func(err error) {
				close(cancelC)
			},
		)
	}
	{

		g.Add(func() error {
			for {
				ticker := time.NewTicker(5 * time.Second)
				select {
				case <-ctxAll.Done():
					level.Warn(logger).Log("msg", "我是模块01退出了，接收到了cancelall")
					return nil
				case <-ticker.C:
					level.Warn(logger).Log("msg", "我是模块01")
				}

			}

		}, func(err error) {

		},
		)
	}
	g.Run()
}
