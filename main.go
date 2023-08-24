package main

import (
	"flag"
	"fmt"
	"github.com/sivukhin/cuemon/lib"
	"github.com/sivukhin/cuemon/lib/auth"
	"os"
	"strings"
)

var (
	bootstrap          = flag.NewFlagSet("bootstrap", flag.ExitOnError)
	bootstrapInput     = bootstrap.String("input", "", "input file with Grafana dashboard JSON (stdin if not provided)")
	bootstrapModule    = bootstrap.String("module", "", "CUE module name")
	bootstrapDir       = bootstrap.String("dir", "", "target directory where cuemon will be initialized")
	bootstrapOverwrite = bootstrap.Bool("overwrite", false, "enable unsafe mode which can overwrite files")

	update          = flag.NewFlagSet("update", flag.ExitOnError)
	updateInput     = update.String("input", "", "input file with Grafana panel JSON (stdin if not provided)")
	updateDir       = update.String("dir", "", "target directory with cuemon setup")
	updateOverwrite = update.Bool("overwrite", false, "enable unsafe mode which can overwrite files")

	push          = flag.NewFlagSet("push", flag.ExitOnError)
	pushDashboard = push.String("dashboard", "", "target CUE file with cuemon dashboard setup")
	pushMessage   = push.String("message", "", "message describing dashboard updates")
	pushGrafana   = push.String("grafana", "", "url to Grafana instance")
	pushTemp      = push.String("temp", "", "temp dashboard name which will be used instead of original dashboard id")
)

func printUsage() {
	fmt.Println("cuemon:")
	bootstrap.Usage()
	update.Usage()
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
	authorization, err := auth.AnalyzeSubjectAuthorization(os.Environ())
	if err != nil {
		return fmt.Errorf("failed to analyze authorization methods: %w", err)
	}
	switch args[0] {
	case "bootstrap":
		if err := bootstrap.Parse(args[1:]); err != nil {
			return fmt.Errorf("bootstrap error: %v", multilineErr(err, errIdent))
		}
		if err := lib.Bootstrap(*bootstrapInput, *bootstrapModule, *bootstrapDir, *bootstrapOverwrite); err != nil {
			return fmt.Errorf("bootstrap error: %v", multilineErr(err, errIdent))
		}
	case "update":
		if err := update.Parse(args[1:]); err != nil {
			return fmt.Errorf("update error: %v", multilineErr(err, errIdent))
		}
		if err := lib.Update(*updateInput, *updateDir, *updateOverwrite); err != nil {
			return fmt.Errorf("update error: %v", multilineErr(err, errIdent))
		}
	case "push":
		if err := push.Parse(args[1:]); err != nil {
			return fmt.Errorf("push error: %v", multilineErr(err, errIdent))
		}
		if err := lib.Push(strings.TrimRight(*pushGrafana, "/"), authorization[authorizationSubject], *pushDashboard, *pushMessage, *pushTemp); err != nil {
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
