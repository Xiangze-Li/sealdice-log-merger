package logmerger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"xiangzeli/logmerger/internal/errcode"

	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
)

var CmdMerge = &cli.Command{
	Name:      "merge",
	Usage:     "合并多个文本Log",
	ArgsUsage: "( InputFile [InputFile ...] | InputDir )",
	Flags:     flags,
	Action:    action,
}

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "main",
		Aliases: []string{"m"},
		Usage:   "主要Log`文件`, 默认为第一个输入文件; 如果使用目录输入, 必须指定",
	},
	&cli.StringFlag{
		Name:    "output",
		Aliases: []string{"o"},
		Usage:   "输出`文件`",
		Value:   "output.txt",
	},
	&cli.BoolFlag{
		Name:    "quiet",
		Aliases: []string{"q"},
		Usage:   "不输出统计信息",
		Value:   false,
	},
}

func action(ctx *cli.Context) error {
	var (
		err    error
		inputs = ctx.Args().Slice()
		mainIn = ctx.String("main")
		output = ctx.String("output")
		quiet  = ctx.Bool("quiet")
	)

	if len(inputs) == 0 {
		return errcode.ErrNoInput
	}

	if len(inputs) == 1 {
		s, err := os.Stat(inputs[0])
		if err != nil {
			return errcode.ErrInputFileError(err)
		}
		if s.IsDir() {
			if len(mainIn) == 0 {
				return errcode.ErrDirInputNoMain
			}
			inputs, _ = filepath.Glob(filepath.Join(inputs[0], "*"))
		}
	}

	output, err = filepath.Abs(output)
	if err != nil {
		return errcode.ErrOutputFileError(err)
	}
	fdOutput, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errcode.ErrOutputFileError(err)
	}
	defer fdOutput.Close()

	var mainPath string
	var restPaths []string
	for _, input := range inputs {
		var path string
		path, err = filepath.Abs(input)
		if err != nil {
			return errcode.ErrInputFileError(err)
		}
		restPaths = append(restPaths, path)
	}

	if len(mainIn) == 0 {
		mainPath = restPaths[0]
	} else {
		mainPath, err = filepath.Abs(mainIn)
		if err != nil {
			return errcode.ErrInputFileError(err)
		}
	}

	mainIdx := slices.Index(restPaths, mainPath)
	if mainIdx == -1 {
		return errcode.ErrMainInputNotListed
	}
	restPaths = append(restPaths[:mainIdx], restPaths[mainIdx+1:]...)

	var merged []logItem
	var stat statistics
	merged, stat, err = mergeLogs(mainPath, restPaths)
	if err != nil {
		if exitCoder, ok := err.(cli.ExitCoder); ok {
			return exitCoder
		}
		return errcode.ErrMergeError(err)
	}

	for _, v := range merged {
		if v.Source == mainPath {
			fmt.Fprintf(fdOutput,
				"%s(%s) %s\n%s\n",
				v.Nickname, v.UserID, v.Time.Format("2006/01/02 15:04:05"), v.Content,
			)
		} else {
			fn := filepath.Base(v.Source)
			lastDot := strings.LastIndex(fn, ".")
			if lastDot != -1 {
				fn = fn[:lastDot]
			}

			fmt.Fprintf(fdOutput,
				"%s(%s) %s\n[%s]\n%s\n",
				v.Nickname, v.UserID, v.Time.Format("2006/01/02 15:04:05"),
				fn, v.Content,
			)
		}
	}

	if !quiet {
		fmt.Fprintf(ctx.App.Writer, statToTable(stat))
	}

	return nil
}

func statToTable(st statistics) string {
	sb := strings.Builder{}
	table := tablewriter.NewWriter(&sb)

	table.SetHeader([]string{"File", "#Log"})
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoFormatHeaders(false)

	table.SetBorders(tablewriter.Border{Left: false, Top: true, Right: false, Bottom: true})
	table.SetCenterSeparator("-")
	table.SetColumnSeparator(" ")

	for _, s := range st.PerFile {
		table.Append([]string{s.FileName, fmt.Sprintf("%d", s.Count)})
	}
	table.SetAlignment(tablewriter.ALIGN_RIGHT)

	table.SetFooter([]string{"TOTAL", fmt.Sprintf("%d", st.Total)})
	table.SetFooterAlignment(tablewriter.ALIGN_RIGHT)

	table.Render()
	return sb.String()
}
