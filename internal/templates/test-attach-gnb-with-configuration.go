package templates

import (
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"sync"
)

func TestAttachGnbWithConfiguration() {

	wg := sync.WaitGroup{}

	cfg, err := config.GetConfig()
	if err != nil {
		//return nil
		log.Fatal("Error in get configuration")
	}

	// synch GNB
	synchGnb := make(chan bool, 1)

	go gnb.InitGnb(cfg, &wg, synchGnb)

	wg.Add(1)

	wg.Wait()
}
