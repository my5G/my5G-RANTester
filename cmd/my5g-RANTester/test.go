package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func testLogCommonInfo(name string, numUE int) {
	const (
		fmtStartTest   = "Starting test function: %v"
		fmtNumUE       = "Number of UEs: %v"
		fmtIPPort      = "%v/%v"
		fmtControlInfo = "Control interface IP/Port: " + fmtIPPort
		fmtDataInfo    = "Data interface IP/Port: " + fmtIPPort
		fmtAMFInfo     = "AMF IP/Port: " + fmtIPPort
	)

	log.Info(logSep)

	msgStartTest := fmt.Sprintf(fmtStartTest, name)
	log.Infof(fmtLog, msgStartTest)

	msgNumUE := fmt.Sprintf(fmtNumUE, numUE)
	log.Infof(fmtLogUE, msgNumUE)

	msgControlInfo := fmt.Sprintf(fmtControlInfo, cfg.GNodeB.ControlIF.Ip, cfg.GNodeB.ControlIF.Port)
	log.Infof(fmtLogGNB, msgControlInfo)

	msgDataInfo := fmt.Sprintf(fmtDataInfo, cfg.GNodeB.DataIF.Ip, cfg.GNodeB.DataIF.Port)
	log.Infof(fmtLogGNB, msgDataInfo)

	msgAMFInfo := fmt.Sprintf(fmtAMFInfo, cfg.AMF.Ip, cfg.AMF.Port)
	log.Infof(fmtLogAMF, msgAMFInfo)

	log.Info(logSep)
}
