package main

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"my5G-RANTester/config"
)

const (
	optLogLevel        = "loglevel"
	optLogLevelAlias1  = "l"
	optLogLevelUsage   = "set the log level of the app, (6 is most verbose)"
	optLogLevelDefault = log.InfoLevel

	optConfig        = "config"
	optConfigAlias   = "c"
	optConfigUsage   = "set the config file to use"
	optConfigDefault = ""
)

// TODO: remove this ugly global variable *hisssss*
var cfg *config.Config

func setupOptions(a *cli.App) {
	a.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:    optLogLevel,
			Aliases: []string{optLogLevelAlias1},
			Usage:   optLogLevelUsage,
			Value:   int(optLogLevelDefault),
		},
		&cli.StringFlag{
			Name:    optConfig,
			Aliases: []string{optConfigAlias},
			Usage:   optConfigUsage,
			Value:   optConfigDefault,
		},
	}
}

func showUsageString(c *cli.Context) error {
	if c.Args().Len() == 0 {
		_ = cli.ShowAppHelp(c)

		return errors.New("no commands specified, stopping")
	}

	return nil
}

func setLogLevel(c *cli.Context) error {
	ll := c.Int(optLogLevel)

	log.SetLevel(log.Level(ll))

	if ll != int(optLogLevelDefault) {
		log.Debugf("setting log level to %v", log.Level(ll))
	}

	return nil
}

func setConfigFile(c *cli.Context) error {
	path := c.String(optConfig)

	loaded, err := config.Load(path)
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("config loaded successfully")
	cfg = loaded

	return nil
}
