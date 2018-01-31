package main

import (
	"fmt"
	"github.com/moladb/rest"
	"github.com/moladb/rest/example/echo/service/v1"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := cli.NewApp()
	app.Name = "echo"
	//app.Version = fmt.Sprintf("\nversion: %s\nbuild_date: %s\ngo_version: %s",
	//	version.VERSION,
	//	version.BUILDDATE,
	//	version.GOVERSION)
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name: "enable-debug",
		},
		cli.BoolFlag{
			Name: "enable-metrics",
		},
		cli.BoolFlag{
			Name: "enable-discovery",
		},
		cli.StringFlag{
			Name:  "addr",
			Value: "0.0.0.0",
		},
		cli.StringFlag{
			Name:   "port",
			Value:  "32000",
			EnvVar: "PORT0",
		},
	}
	app.Usage = "echo-server"
	app.Action = func(c *cli.Context) error {
		bindAddr := fmt.Sprintf("%s:%s", c.String("addr"), c.String("port"))
		srv := rest.NewServer(rest.Config{
			BindAddr:              bindAddr,
			EnableDebug:           c.Bool("enable-pprof"),
			EnableMetrics:         c.Bool("enable-metrics"),
			EnableDiscovery:       c.Bool("enable-discovery"),
			GraceShutdownTimeoutS: 60,
		})

		// register services
		srv.RegisterServiceGroup("v1", v1.NewEchoService())

		go func() {
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
			<-quit
			srv.Shutdown()
		}()

		if err := srv.Run(); err != nil {
			fmt.Println("err:", err)
			os.Exit(1)
		}
		return nil
	}

	app.Run(os.Args)
}
