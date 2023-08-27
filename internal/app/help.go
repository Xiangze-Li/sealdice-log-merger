package app

import (
	"github.com/urfave/cli/v2"
)

var CmdHelp = &cli.Command{
	Name:  "help",
	Usage: "显示帮助信息",
	Action: func(ctx *cli.Context) error {

		if ctx.Args().Len() == 1 {
			return cli.ShowCommandHelp(ctx, ctx.Args().First())
		}

		cli.ShowAppHelp(ctx)
		return nil
	},
	HideHelp:  true,
	ArgsUsage: "[command]",
}
