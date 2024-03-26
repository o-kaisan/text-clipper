package common

import "fmt"

// 起動オプション
const (
	USAGE = `usage: text-clipper [option]
options:
-h,--help             show this usage
-v,--version          display the version`
)

type (
	Options struct {
		Help    bool
		Version bool
		Debug   bool
	}
)

func newOptions() *Options {
	return &Options{
		Help:    false,
		Version: false,
		Debug:   false,
	}
}

type Args []string

func ParseOptions(args Args) *Options {
	opts := newOptions()
	for i := 1; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "-h", "--help":
			opts.Help = true
		case "-v", "--version":
			opts.Version = true
		case "-D", "--debug":
			opts.Debug = true
		default:
			panic(fmt.Sprintf("unrecognized option %s", arg))
		}
	}
	return opts
}
