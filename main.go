package main

import (
	"flag"
	"fmt"
	"github.com/sivukhin/cuemon/lib"
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
)

func printUsageAndExit() {
	fmt.Println("cuemon:")
	bootstrap.Usage()
	update.Usage()
	os.Exit(2)
}

func multilineErr(err error, ident string) string {
	result := err.Error()
	lines := strings.Split(result, "\n")
	if len(lines) == 1 {
		return lines[0]
	}
	return strings.Join(append([]string{""}, lines...), "\n"+ident)
}

const errIdent = "  "

func main() {
	if len(os.Args) < 2 {
		printUsageAndExit()
	}

	switch os.Args[1] {
	case "bootstrap":
		if err := bootstrap.Parse(os.Args[2:]); err != nil {
			fmt.Printf("bootstrap error: %v\n", multilineErr(err, errIdent))
			printUsageAndExit()
		}
		if err := lib.Bootstrap(*bootstrapInput, *bootstrapModule, *bootstrapDir, *bootstrapOverwrite); err != nil {
			fmt.Printf("bootstrap error: %v\n", multilineErr(err, errIdent))
			os.Exit(1)
		}
	case "update":
		if err := update.Parse(os.Args[2:]); err != nil {
			fmt.Printf("update error: %v\n", multilineErr(err, errIdent))
			printUsageAndExit()
		}
		if err := lib.Update(*updateInput, *updateDir, *updateOverwrite); err != nil {
			fmt.Printf("update error: %v\n", multilineErr(err, errIdent))
			os.Exit(1)
		}
	default:
		printUsageAndExit()
	}
}
