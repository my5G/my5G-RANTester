package main

import (
	"github.com/urfave/cli/v2"
)

const (
	argNumUE        = "number-of-ues"
	argNumUEDefault = 1

	cmdUeName  = "ue"
	cmdUeUsage = "Testing an UE attached with configuration"

	cmdGnbName  = "gnb"
	cmdGnbUsage = "Testing a GNB attached with configuration"

	cmdLoadTestName  = "load-test"
	cmdLoadTestUsage = `Load endurance stress tests.
	Example for testing multiple UEs: load-test -n 5
	`
)

func setupCommands(a *cli.App) {
	var commands []*cli.Command

	loadTestFlags := []cli.Flag{
		&cli.IntFlag{Name: argNumUE, Value: argNumUEDefault, Aliases: []string{"n"}},
	}

	for _, cmd := range []struct {
		name, usage string
		fn          func(c *cli.Context) error
		flags       []cli.Flag
	}{
		{cmdUeName, cmdUeUsage, testUE, nil},
		{cmdGnbName, cmdGnbUsage, testGNB, nil},
		{cmdLoadTestName, cmdLoadTestUsage, testRegisterMultiUE, loadTestFlags},
	} {
		commands = append(commands, &cli.Command{
			Name:   cmd.name,
			Usage:  cmd.usage,
			Action: cmd.fn,
			Flags:  cmd.flags,
		})
	}

	a.Commands = commands
}
