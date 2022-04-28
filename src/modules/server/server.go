package main

import (
	"Distributed-System-Awareness-Platform/src/models"
	"Distributed-System-Awareness-Platform/src/modules/server/cloudsync"
	"Distributed-System-Awareness-Platform/src/modules/server/config"
	"Distributed-System-Awareness-Platform/src/modules/server/memoryindex"
	"Distributed-System-Awareness-Platform/src/modules/server/metric"
	"Distributed-System-Awareness-Platform/src/modules/server/rpc"
	"Distributed-System-Awareness-Platform/src/modules/server/statistics"
	"Distributed-System-Awareness-Platform/src/modules/server/web"
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
	//初始化mysql
	models.InitMySQL(sConfig.MysqlS)
	level.Info(logger).Log("msg", "load.mysql.success", "db.num", len(models.DB))
	//初始化倒排索引模块
	memoryindex.Init(logger, sConfig.IndexModules)
	// 注册stree相关的metrics
	metric.NewMetrics()
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
		// rpc server
		g.Add(func() error {
			errChan := make(chan error, 1)
			go func() {
				errChan <- rpc.Start(":8080", logger)
			}()
			select {
			case err := <-errChan:
				level.Error(logger).Log("msg", "rpc server error", "err", err)
				return err
			case <-ctxAll.Done():
				level.Info(logger).Log("msg", "receive_quit_signal_rpc_server_exit")
				return nil
			}

		}, func(err error) {
			cancelAll()
		},
		)
	}
	{

		g.Add(func() error {
			errChan := make(chan error, 1)
			go func() {
				errChan <- web.StartGin(sConfig.HttpAddr, logger)
			}()
			select {
			case err := <-errChan:
				level.Error(logger).Log("msg", "web server error", "err", err)
				return err
			case <-ctxAll.Done():
				level.Info(logger).Log("msg", "receive_quit_signal_web_server_exit")
				return nil
			}

		}, func(err error) {
			cancelAll()
		},
		)
	}
	{
		// 公有云同步
		if sConfig.PCC.Enable {
			cloudsync.Init(logger)
			g.Add(func() error {
				err := cloudsync.CloudSyncManager(ctxAll, logger)
				if err != nil {
					level.Error(logger).Log("msg", "cloudsync.CloudSyncManager.error", "err", err)

				}
				return err

			}, func(err error) {
				cancelAll()
			},
			)
		}
	}
	{
		// 刷新倒排索引
		g.Add(func() error {
			err := memoryindex.RevertedIndexSyncManager(ctxAll, logger)
			if err != nil {
				level.Error(logger).Log("msg", "mem_index.RevertedIndexSyncManager.error", "err", err)

			}
			return err

		}, func(err error) {
			cancelAll()
		},
		)
	}
	{
		// 统计资源分布

		g.Add(func() error {
			err := statistics.TreeNodeStatisticsManager(ctxAll, logger)
			if err != nil {
				level.Error(logger).Log("msg", "statistics.TreeNodeStatisticsManager.error", "err", err)

			}
			return err

		}, func(err error) {
			cancelAll()
		},
		)
	}
	g.Run()
}
