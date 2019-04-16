package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/errwrap"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vbatts/quay-scheduled-build/quay"
	"github.com/vbatts/quay-scheduled-build/types"
)

func main() {
	app := cli.NewApp()
	app.Name = "quay-scheduled-build"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		{
			Name:        "serve",
			Aliases:     []string{"s"},
			Usage:       "the build scheduler",
			Description: "run the scheduler for your builds on quay.io",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config",
					Usage: "build config to manage",
					Value: "quay-build.json",
				},
			},
			Action: func(c *cli.Context) error {
				fh, err := os.Open(c.String("config"))
				if err != nil {
					return cli.NewExitError(errwrap.Wrapf("config file error: {{err}}", err), 1)
				}
				defer fh.Close()
				dec := json.NewDecoder(fh)
				cfg := types.Config{}
				err = dec.Decode(&cfg)
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
			},
		},
		{
			Name:        "oneshot",
			Aliases:     []string{"o"},
			Usage:       "trigger the builds right meow",
			Description: "trigger your builds on quay.io right meow.\n\tIf there are multiple builds in your config, they are triggered in serial.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config",
					Usage: "build config to manage",
					Value: "quay-build.json",
				},
			},
			Action: func(c *cli.Context) error {
				fh, err := os.Open(c.String("config"))
				if err != nil {
					return cli.NewExitError(errwrap.Wrapf("config file error: {{err}}", err), 1)
				}
				defer fh.Close()
				logrus.Infof("reading config from %q", c.String("config"))
				dec := json.NewDecoder(fh)
				cfg := types.Config{}
				err = dec.Decode(&cfg)
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
			},
		},
		{
			Name:        "generate",
			Aliases:     []string{"gen", "g"},
			Usage:       "generate output",
			Description: "helper to produce config output",
			Flags:       []cli.Flag{},
			Subcommands: []cli.Command{
				{
					Name:    "config",
					Aliases: []string{"c"},
					Usage:   "helper to get a new buildref configuration",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:   "token",
							Usage:  "quay.io oauth token for the build",
							EnvVar: "BUILD_TOKEN",
						},
						cli.StringFlag{
							Name:   "repo",
							Usage:  "quay.io container repo for the build",
							EnvVar: "BUILD_REPO",
						},
						cli.StringFlag{
							Name:   "schedule",
							Usage:  "cron style schedule to trigger this containter build (more info https://godoc.org/github.com/robfig/cron#hdr-CRON_Expression_Format)",
							EnvVar: "BUILD_SCHEDULE",
							Value:  "@weekly",
						},
						cli.StringSliceFlag{
							Name:   "tags",
							Usage:  "container name tags to apply to this build",
							EnvVar: "BUILD_TAG",
						},
						cli.StringFlag{
							Name:   "robot",
							Usage:  "quay.io robot account username for the build",
							EnvVar: "BUILD_ROBOT",
						},
						cli.StringFlag{
							Name:   "archive-url",
							Usage:  "URL to the source of the build (which includes the Dockerfile)",
							EnvVar: "BUILD_ARCHIVE_URL",
						},
						cli.StringFlag{
							Name:   "dockerfile-path",
							Usage:  "path (within the source archive) to the Dockerfile",
							EnvVar: "BUILD_DOCKERFILE_PATH",
						},
						cli.StringFlag{
							Name:   "subdirectory",
							Usage:  "path (within the source archive) to the root of the build directory",
							EnvVar: "BUILD_SUBDIRECTORY",
						},
					},
					Action: func(c *cli.Context) error {
						cfg := types.Config{
							Builds: []types.Build{
								types.Build{
									Token:    c.String("token"),
									QuayRepo: c.String("repo"),
									Schedule: c.String("schedule"),
									BuildRef: types.BuildRef{
										Tags:           c.StringSlice("tags"),
										PullRobot:      c.String("robot"),
										ArchiveURL:     c.String("archive-url"),
										DockerfilePath: c.String("dockerfile-path"),
										Subdirectory:   c.String("subdirectory"),
									},
								},
							},
						}

						buf, err := json.MarshalIndent(cfg, "", "  ")
						if err != nil {
							return cli.NewExitError(err, 1)
						}
						fmt.Println(string(buf))

						return nil
					},
				},
			},
			Action: func(c *cli.Context) error {
				return cli.ShowSubcommandHelp(c)
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelpAndExit(c, 1)
		return nil
	}

	_ = app.Run(os.Args)
}
