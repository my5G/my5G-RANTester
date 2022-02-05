package templates

import (
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"my5G-RANTester/internal/monitoring"
	"strconv"
	"sync"
	"time"
)

func TestRqsPerTime(numRqs int, interval int) int64 {

	wg := sync.WaitGroup{}

	monitor := monitoring.Monitor{
		Rqs: 0,
	}

	cfg, err := config.GetConfig()
	if err != nil {
		//return nil
		log.Fatal("Error in get configuration")
	}

	ranPort := 800
	for i := 1; i <= numRqs; i++ {
		cfg.GNodeB.PlmnList.GnbId = gnbIdGenerator(i)
		cfg.GNodeB.ControlIF.Port = ranPort
		go gnb.InitGnbForLoad(cfg, &wg, &monitor, time.Duration(interval))
		wg.Add(1)
		ranPort++
	}

	wg.Wait()

	return monitor.GetRqs()
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
