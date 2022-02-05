package templates

import (
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"strconv"
	"sync"
)

func TestRqsPerSecond(numRqs int) {

	wg := sync.WaitGroup{}

	cfg, err := config.GetConfig()
	if err != nil {
		//return nil
		log.Fatal("Error in get configuration")
	}

	for i := 1; i <= numRqs; i++ {
		cfg.GNodeB.PlmnList.GnbId = gnbIdGenerator(i)
		go gnb.InitGnb(cfg, &wg)
		wg.Add(1)
	}

	wg.Wait()
}

func gnbIdGenerator(i int) string {

	var base string
	switch true {
	case i < 10:
		base = "00000"
	case i < 100:
		base = "0000"
	case i >= 100:
		base = "000"
	}

	gnbId := base + strconv.Itoa(i)
	return gnbId
}
