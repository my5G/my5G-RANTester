package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"my5G-RANTester/internal/templates"
)

func testRegisterMultiUE(c *cli.Context) error {
	if !c.IsSet(argNumUE) {
		log.Info(c.Command.Usage)

		return nil
	}

	const (
		name = "Testing registration of multiple UEs"
	)

	numUE := c.Int(argNumUE)

	testLogCommonInfo(name, numUE)
	templates.TestMultiUesInQueue(numUE)

	return nil
}
