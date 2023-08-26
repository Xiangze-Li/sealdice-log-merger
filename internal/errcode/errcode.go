package errcode

import (
	"github.com/urfave/cli/v2"
)

const (
	_ = -iota
	ErrUnknown
	ErrCodeNoInput
	ErrCodeMainInputNotListed
	ErrCodeDirInputNoMain
	ErrCodeOutputFileError
	ErrCodeInputFileError
	ErrCodeMergeError
	ErrCodeReadError
)

var (
	ErrNoInput            = cli.Exit("没有指定输入文件", ErrCodeNoInput)
	ErrMainInputNotListed = cli.Exit("主要输入文件没有在输入文件列表中", ErrCodeMainInputNotListed)
	ErrDirInputNoMain     = cli.Exit("使用目录输入但未指定主要Log", ErrCodeDirInputNoMain)
	ErrOutputFileError    = func(err error) cli.ExitCoder { return cli.Exit(err, ErrCodeOutputFileError) }
	ErrInputFileError     = func(err error) cli.ExitCoder { return cli.Exit(err, ErrCodeInputFileError) }
	ErrMergeError         = func(err error) cli.ExitCoder { return cli.Exit(err, ErrCodeMergeError) }
	ErrReadError          = func(err error) cli.ExitCoder { return cli.Exit(err, ErrCodeReadError) }
)
