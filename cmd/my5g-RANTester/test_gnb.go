package main

import (
	"github.com/urfave/cli/v2"

	"my5G-RANTester/internal/templates"
)

func testGNB(_ *cli.Context) error {
	const (
		name = "Testing an gnb attached with configuration"
	)

	testLogCommonInfo(name, argNumUEDefault)

	templates.TestAttachGnbWithConfiguration()

	return nil
}
