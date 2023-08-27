package app

import "github.com/urfave/cli/v2"

var CmdVersion = &cli.Command{
	Name:  "version",
	Usage: "显示版本信息",
	Action: func(ctx *cli.Context) error {
		cli.ShowVersion(ctx)
		return nil
	},
	HideHelp:  true,
	UsageText: `logmerger version`,
}
