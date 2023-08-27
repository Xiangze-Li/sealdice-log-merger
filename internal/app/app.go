package app

import (
	"fmt"
	"xiangzeli/logmerger/internal/errcode"
	"xiangzeli/logmerger/internal/logmerger"
	"xiangzeli/logmerger/internal/preprocessor"

	"github.com/urfave/cli/v2"
)

const version = "2.0.0"

// App is the main application.
var App = &cli.App{
	Name:  "logmerger",
	Usage: "合并多个群的Log",

	Commands: []*cli.Command{
		preprocessor.CmdPreProcess,
		logmerger.CmdMerge,
		CmdVersion,
		CmdHelp,
	},

	HideHelp:    true,
	HideVersion: true,

	Version: version,
	Authors: []*cli.Author{{Name: "Xiangze Li", Email: "lee_johnson@qq.com"}},
	ExitErrHandler: func(cCtx *cli.Context, err error) {
		if err == nil {
			return
		}
		fmt.Fprintf(cCtx.App.Writer, "Error: %v\n\n", err)
		if exitCoder, ok := err.(cli.ExitCoder); ok {
			cli.ShowAppHelpAndExit(cCtx, exitCoder.ExitCode())
		} else {
			cli.ShowAppHelpAndExit(cCtx, errcode.ErrUnknown)
		}
	},
}
