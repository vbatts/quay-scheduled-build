package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/errwrap"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vbatts/quay-scheduled-build/quay"
	"github.com/vbatts/quay-scheduled-build/types"
)

var cmdServe = cli.Command{
	Name:        "serve",
	Aliases:     []string{"s"},
	Usage:       "the build scheduler",
	Description: "run the scheduler for your builds on quay.io",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:   "config",
			Usage:  "configuration for your builds",
			Value:  "quay-build.json",
			EnvVar: "BUILD_CONFIG_FILE",
		},
		cli.StringFlag{
			Name:   "config-json",
			Usage:  "full JSON blob of your build configuration (if provided this takes precedence)",
			EnvVar: "BUILD_CONFIG",
		},
	},
	Action: cmdServeAction,
}

func cmdServeAction(c *cli.Context) error {
	var rdr io.Reader

	if c.String("config-json") != "" {
		rdr = bytes.NewBufferString(c.String("config-json"))
		logrus.Infof("reading config from BUILD_CONFIG")
	} else {
		fh, err := os.Open(c.String("config"))
		if err != nil {
			return cli.NewExitError(errwrap.Wrapf("config file error: {{err}}", err), 1)
		}
		defer fh.Close()

		rdr = fh
		logrus.Infof("reading config from %q", c.String("config"))
	}

	dec := json.NewDecoder(rdr)
	cfg := types.Config{}
	err := dec.Decode(&cfg)
	if err != nil {
		return cli.NewExitError(errwrap.Wrapf("config parse error: {{err}}", err), 1)
	}

	sched := cron.New()
	logrus.Info("readying the scheduler ...")
	for i, bldinfo := range cfg.Builds {
		if bldinfo.Schedule != "" {
			_, err := cron.Parse(bldinfo.Schedule)
			if err != nil {
				logrus.Errorf("failed to parse schedule %q", bldinfo.Schedule)
				continue
			}
			logrus.Infof("queuing build of %q for %q", bldinfo.QuayRepo, bldinfo.Schedule)
			err = sched.AddFunc(bldinfo.Schedule, func() {
				resp, err := quay.BuildRequest(bldinfo)
				if err != nil {
					logrus.Fatal(errwrap.Wrapf(fmt.Sprintf("[%d] BuildRequest error: {{err}}", i), err))
				}
				buf, err := json.Marshal(resp)
				if err != nil {
					logrus.Infof("%v", resp)
					return
				}
				logrus.Info(string(buf))
			})
			if err != nil {
				return cli.NewExitError(errwrap.Wrapf("Scheduling error: {{err}}", err), 1)
			}
		}
	}
	if len(sched.Entries()) <= 0 {
		return cli.NewExitError("error: no builds scheduled", 1)
	}
	logrus.Info("running the build scheduler ...")
	sched.Run()

	return nil
}
