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
	initLogger()

	log.Infof(fmtMsgVersion, version)

	app := &cli.App{}

	setupCommands(app)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func initLogger() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)

	spew.Config.Indent = "\t"
}
