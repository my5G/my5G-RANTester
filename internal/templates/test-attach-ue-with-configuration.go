package templates

import (
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"my5G-RANTester/internal/control_test_engine/ue"
	"sync"
	"time"
)

func TestAttachUeWithConfiguration(cfg config.Config) error {
	wg := sync.WaitGroup{}

	go gnb.InitGnb(cfg)

	wg.Add(1)

	time.Sleep(time.Second)

	go ue.RegistrationUe(cfg, 1)

	wg.Add(1)

	wg.Wait()

	return nil
}
