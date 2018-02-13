package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
)

var (
	// Version is set at compile time, passing -ldflags "-X main.Version=<build_version>"
	Version string

	logInfo  *log.Logger
	logError *log.Logger
	logDebug *log.Logger

	// Globals
	c *cli.Context // CLI context
)

func init() {
	logInfo = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	logError = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	logDebug = log.New(os.Stderr, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	app := cli.NewApp()
	app.Name = "smoketest"
	app.Author = "Marcin Ciszak <https://github.com/marcinc>"
	app.Version = Version
	app.Usage = "simple kuberang smoke test wrapper"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug output",
			EnvVar: "DEBUG",
		},
		cli.BoolFlag{
			Name:   "push-metrics",
			Usage:  "push metrics to sysdig",
			EnvVar: "PUSH_METRICS",
		},
		cli.StringFlag{
			Name:   "namespace, n",
			Usage:  "kubernetes namespace `namespace`",
			EnvVar: "KUBE_NAMESPACE",
			Value:  "smoke-test",
		},
		cli.StringFlag{
			Name:   "registry-url, r",
			Usage:  "smoke test Docker registry override",
			EnvVar: "REGISTRY_URL",
			Value:  "https://index.docker.io/v1/",
		},
		cli.DurationFlag{
			Name:   "interval",
			Usage:  "smoke test check interval `INTERVAL`",
			EnvVar: "INTERVAL",
			Value:  time.Duration(300000) * time.Millisecond,
		},
	}

	app.Action = func(cx *cli.Context) error {
		c = cx

		if err := run(); err != nil {
			logError.Print(err)
			return cli.NewExitError("", 1)
		}
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		logError.Fatal(err)
	}
}

func run() error {
	for {
		err := runCmd("kuberang --namespace " + c.String("namespace") + " --registry-url " + c.String("registry-url"), kuberangOutputHandler)
		if err != nil {
			cleanupServices()
			cleanupDeployments()
		} else {
			networkErrorRateCheck()
		}
		logInfo.Print("-------")
		time.Sleep(c.Duration("interval"))
	}
	return nil
}
