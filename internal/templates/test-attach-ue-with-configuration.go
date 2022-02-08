package templates

import (
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"my5G-RANTester/internal/control_test_engine/ue"
	"my5G-RANTester/internal/monitoring"
	"sync"
	"time"
)

func TestAttachUeWithConfiguration() {

	wg := sync.WaitGroup{}

	cfg, err := config.GetConfig()
	if err != nil {
		//return nil
		log.Fatal("Error in get configuration")
	}

	go gnb.InitGnb(cfg, &wg)

	wg.Add(1)

	time.Sleep(1 * time.Second)

	go ue.RegistrationUe(cfg, 1, &wg)

	wg.Add(1)

	wg.Wait()
}

// testa a latência do ue no registro
func TestUeLatency() int64 {

	wg := sync.WaitGroup{}

	monitor := monitoring.Monitor{
		LtRegisterLocal: 0,
	}

	cfg, err := config.GetConfig()
	if err != nil {
		//return nil
		log.Fatal("Error in get configuration")
	}

	go gnb.InitGnb(cfg, &wg)

	wg.Add(1)

	time.Sleep(1 * time.Second)

	go ue.RegistrationUeMonitor(cfg, 1, &monitor, &wg)

	wg.Add(1)

	wg.Wait()

	wg.Done()

	return monitor.LtRegisterLocal
}
