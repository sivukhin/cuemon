package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/google/shlex"
	"github.com/sivukhin/cuemon/lib"
	"log"
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

const errIdent = "  "

type RunContext struct {
	cookie      string
	interactive bool
}

func (c *RunContext) run(args []string) error {
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
		if c.cookie == "" && c.interactive {
			cookie, err := InitializeCookie("https://moj-monitoring.sharechat.com")
			if err != nil {
				return fmt.Errorf("unable to initialize run context: %v", err)
			}
			c.cookie = cookie
		}
		if err := push.Parse(args[1:]); err != nil {
			return fmt.Errorf("push error: %v", multilineErr(err, errIdent))
		}
		if err := lib.Push(strings.TrimRight(*pushGrafana, "/"), c.cookie, *pushDashboard, *pushMessage, *pushTemp); err != nil {
			return fmt.Errorf("push error: %v", multilineErr(err, errIdent))
		}
	case "help":
		printUsage()
	default:
		return fmt.Errorf("unknown command: %v", args[0])
	}
	return nil
}

func InitializeCookie(grafanaUrl string) (string, error) {
	dir, err := os.MkdirTemp("", "chromedp-example")
	if err != nil {
		log.Fatal(err)
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[3:], chromedp.DisableGPU, chromedp.UserDataDir(dir))
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	var cookie string
	err = chromedp.Run(taskCtx,
		chromedp.Navigate(grafanaUrl),
		chromedp.WaitVisible(".dashboard-container", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetCookies().Do(ctx)
			if err != nil {
				return err
			}
			cookieValues := make([]string, 0)
			for _, cookie := range cookies {
				cookieValues = append(cookieValues, fmt.Sprintf("%v=%v", cookie.Name, cookie.Value))
			}
			cookie = strings.Join(cookieValues, "; ")
			cancel()
			return nil
		}),
	)
	if err != nil {
		return cookie, fmt.Errorf("unable to initialize context: %w", err)
	}
	<-taskCtx.Done()
	return cookie, nil
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}
	runContext := RunContext{cookie: os.Getenv("GRAFANA_COOKIE")}
	if os.Args[1] == "interactive" {
		runContext.interactive = true
		fmt.Printf("started interactive mode\n$> ")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			args, err := shlex.Split(scanner.Text())
			if err != nil {
				fmt.Printf("invalid command format: %v\n", err)
			} else {
				err = runContext.run(args)
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			}
			fmt.Printf("$> ")
		}
	} else {
		err := runContext.run(os.Args[1:])
		if err != nil {
			fmt.Printf("%v\n", err)
			printUsage()
			os.Exit(2)
		}
	}
}
