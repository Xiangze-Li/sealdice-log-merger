package preprocessor

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"xiangzeli/logmerger/internal/errcode"

	"github.com/urfave/cli/v2"
)

var CmdPreProcess = &cli.Command{
	Name:      "preprocess",
	Usage:     "从海豹Log压缩包中提取文本Log",
	ArgsUsage: "( InputFile [InputFile ...] | InputDir )",

	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "输出文件/目录",
			Value:   ".",
		},
	},

	// HideHelp: true,
	HideHelpCommand: true,

	Action: action,
}

var regexpZipFn = regexp.MustCompile(
	`^(.*)_(.*)\.[0-9]+\.zip$`,
)

func action(ctx *cli.Context) error {
	var (
		err    error
		inputs = ctx.Args().Slice()
		output = ctx.String("output")

		outputDir  string
		outputFile string
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
			inputs, _ = filepath.Glob(filepath.Join(inputs[0], "*"))
		}
	}

	{
		outputDir, outputFile = filepath.Split(output)

		if len(inputs) > 1 {
			switch {
			case len(outputDir) == 0 && len(outputFile) > 0:
				outputDir, outputFile = outputFile, ""
			case len(outputDir) > 0 && len(outputFile) > 0:
				return errcode.ErrOutputFileError(fmt.Errorf("多个输入文件不可指定输出文件"))
			default:
				// no-op
			}
		}

		if len(outputDir) == 0 {
			outputDir = "."
		}

		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			return errcode.ErrOutputFileError(err)
		}
	}

	for _, input := range inputs {
		inputBase := filepath.Base(input)
		matches := regexpZipFn.FindStringSubmatch(inputBase)
		if len(matches) != 3 {
			fmt.Fprintf(ctx.App.Writer, "文件 %q 看起来并非海豹Log\n", input)
			continue
		} else {
			fmt.Fprintf(ctx.App.Writer, "正在处理 %s (%s)\n", matches[2], matches[1])
		}

		var (
			grpName string = matches[2]
			path    string
			zipR    *zip.ReadCloser
			zipFd   fs.File
			outFd   *os.File
			outFn   string = filepath.Join(outputDir, grpName+".txt")
		)
		if len(outputFile) > 0 {
			outFn = filepath.Join(outputDir, outputFile)
		}

		path, err = filepath.Abs(input)
		if err != nil {
			return errcode.ErrInputFileError(err)
		}

		zipR, err = zip.OpenReader(path)
		if err != nil {
			return errcode.ErrInputFileError(err)
		}
		defer zipR.Close()

		zipFd, err = zipR.Open("文本log.txt")
		if err != nil {
			return errcode.ErrInputFileError(err)
		}
		defer zipFd.Close()

		outFd, err = os.OpenFile(outFn, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			return errcode.ErrOutputFileError(err)
		}
		defer outFd.Close()

		_, err = io.Copy(outFd, zipFd)
		if err != nil {
			return errcode.ErrOutputFileError(err)
		}
	}

	fmt.Fprintf(ctx.App.Writer, "处理完成\n")
	return nil
}
