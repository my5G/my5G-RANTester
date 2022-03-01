package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	version       = "0.1"
	fmtMsgVersion = "my5G-RANTester version %v"
)

func main() {
	app := &cli.App{Before: runAllActions}

	initLogger()
	log.Infof(fmtMsgVersion, version)
	setupOptions(app)
	setupCommands(app)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initLogger() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	log.SetLevel(0)

	spew.Config.Indent = "\t"
}

func runAllActions(c *cli.Context) (err error) {
	for _, fn := range []func(c *cli.Context) error{
		showUsageString,
		setLogLevel,
		setConfigFile,
	} {
		err = fn(c)
		if err != nil {
			break
		}
	}

	return err
}
