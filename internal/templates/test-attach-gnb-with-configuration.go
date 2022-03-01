package templates

import (
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"sync"
)

func TestAttachGnbWithConfiguration(cfg config.Config) error {

	wg := sync.WaitGroup{}

	// wrong messages:
	// cfg.GNodeB.PlmnList.Mcc = "891"
	// cfg.GNodeB.PlmnList.Mnc = "23"
	// cfg.GNodeB.PlmnList.Tac = "000002"
	// cfg.GNodeB.SliceSupportList.St = "10"
	// cfg.GNodeB.SliceSupportList.Sst = "010239"

	go gnb.InitGnb(cfg, &wg)

	wg.Add(1)

	wg.Wait()

	return nil
}
