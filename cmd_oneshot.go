package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/errwrap"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vbatts/quay-scheduled-build/quay"
	"github.com/vbatts/quay-scheduled-build/types"
)

var cmdOneshot = cli.Command{
	Name:        "oneshot",
	Aliases:     []string{"o"},
	Usage:       "trigger the builds right meow",
	Description: "trigger your builds on quay.io right meow.\n\tIf there are multiple builds in your config, they are triggered in serial.",
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
	Action: cmdOneshotAction,
}

func cmdOneshotAction(c *cli.Context) error {
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

	for i, bldinfo := range cfg.Builds {
		logrus.Infof("requesting imediate build of %q", bldinfo.QuayRepo)
		resp, err := quay.BuildRequest(bldinfo)
		if err != nil {
			return cli.NewExitError(errwrap.Wrapf(fmt.Sprintf("[%d] BuildRequest error: {{err}}", i), err), 1)
		}
		buf, err := json.Marshal(resp)
		if err != nil {
			logrus.Infof("%v", resp)
			continue
		}
		logrus.Info(string(buf))
	}
	return nil
}
