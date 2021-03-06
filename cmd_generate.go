package main

import (
	"encoding/json"
	"fmt"

	"github.com/urfave/cli"
	"github.com/vbatts/quay-scheduled-build/types"
)

var cmdGenerate = cli.Command{
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
}
