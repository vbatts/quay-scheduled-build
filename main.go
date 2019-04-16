package main

import (
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "quay-scheduled-build"
	app.Usage = "schedule container builds on quay.io"
	app.Description = `this utility is a helper to trigger container builds on quay.io
	More information at https://github.com/vbatts/quay-scheduled-build

	Environment variable BUILD_COMMAND=serve is will function the same as 'quay-scheduled-build serve',
	and BUILD_COMMAND=oneshot is will function the same as 'quay-scheduled-build oneshot'.
	This is particularly useful when using the quay.io/vbatts/quay-scheduled-build container image.
`
	app.Authors = []cli.Author{
		{
			Name:  "Vincent Batts",
			Email: "vbatts@hashbangbash.com",
		},
	}
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		cmdServe,
		cmdOneshot,
		cmdGenerate,
	}
	app.Action = cmdMainAction

	_ = app.Run(os.Args)
}

func cmdMainAction(c *cli.Context) error {
	switch os.Getenv("BUILD_COMMAND") {
	case "serve":
		return cmdServe.Run(c)
	case "oneshot":
		return cmdOneshot.Run(c)
	}
	cli.ShowAppHelpAndExit(c, 1)
	return nil
}
