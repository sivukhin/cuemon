package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sivukhin/cuemon/lib"
	"github.com/sivukhin/cuemon/lib/auth"
)

type ArrayFlags []string

func (f *ArrayFlags) String() string { return strings.Join(*f, ",") }
func (f *ArrayFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

var (
	bootstrap          = flag.NewFlagSet("bootstrap", flag.ExitOnError)
	bootstrapInput     = bootstrap.String("input", "", "input file with Grafana dashboard JSON (stdin if not provided)")
	bootstrapModule    = bootstrap.String("module", "", "CUE module name")
	bootstrapDir       = bootstrap.String("dir", "", "target directory where cuemon will be initialized")
	bootstrapOverwrite = bootstrap.Bool("overwrite", false, "enable unsafe mode which can overwrite files")

	export = flag.NewFlagSet("export", flag.ExitOnError)

	push           = flag.NewFlagSet("push", flag.ExitOnError)
	pushMessage    = push.String("message", "", "message describing dashboard updates")
	pushGrafana    = push.String("grafana", "", "url to Grafana instance")
	pushPlayground = push.String("playground", "", "playground dashboard name which will be updated instead of original dashboard id")
)

func printUsage() {
	fmt.Println("cuemon:")
	bootstrap.Usage()
	export.Usage()
	push.Usage()
}

func multilineErr(err error, ident string) string {
	result := err.Error()
	lines := strings.Split(result, "\n")
	if len(lines) == 1 {
		return lines[0]
	}
	return strings.Join(append([]string{""}, lines...), "\n"+ident)
}

const (
	errIdent             = "  "
	authorizationSubject = "GRAFANA"
)

func run(args []string) error {
	var tags []string
	tagSet := func(s string) error {
		tags = append(tags, s)
		return nil
	}
	bootstrap.Func("t", "tags for cue export", tagSet)
	export.Func("t", "tags for cue export", tagSet)

	authorization, err := auth.AnalyzeSubjectAuthorization(os.Environ())
	if err != nil {
		return fmt.Errorf("failed to analyze authorization methods: %w", err)
	}
	switch args[0] {
	case "export":
		if err := export.Parse(args[1:]); err != nil {
			return fmt.Errorf("export error: %v", multilineErr(err, errIdent))
		}
		result, err := lib.Export(export.Args(), tags)
		if err != nil {
			return fmt.Errorf("export error: %v", multilineErr(err, errIdent))
		}
		fmt.Printf("%v", result)
	case "bootstrap":
		if err := bootstrap.Parse(args[1:]); err != nil {
			return fmt.Errorf("bootstrap error: %v", multilineErr(err, errIdent))
		}
		if err := lib.Bootstrap(*bootstrapInput, *bootstrapModule, *bootstrapDir, *bootstrapOverwrite); err != nil {
			return fmt.Errorf("bootstrap error: %v", multilineErr(err, errIdent))
		}
	case "push":
		if err := push.Parse(args[1:]); err != nil {
			return fmt.Errorf("push error: %v", multilineErr(err, errIdent))
		}
		if err := lib.Push(strings.TrimRight(*pushGrafana, "/"), authorization[authorizationSubject], *pushMessage, *pushPlayground, tags, push.Args()); err != nil {
			return fmt.Errorf("push error: %v", multilineErr(err, errIdent))
		}
	case "help":
		printUsage()
	default:
		return fmt.Errorf("unknown command: %v", args[0])
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}
	err := run(os.Args[1:])
	if err != nil {
		fmt.Printf("%v\n", err)
		printUsage()
		os.Exit(2)
	}
}
