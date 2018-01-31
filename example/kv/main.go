/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/moladb/rest"
	"github.com/moladb/rest/example/kv/kv-service/v0"
	"github.com/moladb/rest/example/kv/kv-service/v1"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "kv-server"
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
			Name:  "bind-addr",
			Value: "0.0.0.0:8500",
		},
	}
	app.Usage = "kv-server"
	app.Action = func(c *cli.Context) error {
		srv := rest.NewServer(rest.Config{
			BindAddr:              c.String("bind-addr"),
			EnableDebug:           c.Bool("enable-pprof"),
			EnableMetrics:         c.Bool("enable-metrics"),
			EnableDiscovery:       c.Bool("enable-discovery"),
			GraceShutdownTimeoutS: 60,
		})

		// register services
		srv.RegisterServiceGroup("v0", v0.NewKVService())
		srv.RegisterServiceGroup("v1", v1.NewKVService())

		go func() {
			quit := make(chan os.Signal, 1)
			signal.Notify(quit, os.Interrupt)
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
