package templates

import (
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"my5G-RANTester/internal/control_test_engine/ue"
	"strconv"
	"sync"
	"time"
	"fmt"
	
	log_time "my5G-RANTester/internal/analytics/log_time"
)

func TestMultiUesInParallel(numUes int, delayUes int, delayStart int) {

	wg := sync.WaitGroup{}

	cfg, err := config.GetConfig()
	if err != nil {
		//return nil
		log.Fatal("Error in get configuration")
	}

	go gnb.InitGnb(cfg, &wg)

	wg.Add(1)

	time.Sleep(time.Duration(delayStart) * time.Second)
    msin :=  cfg.Ue.Msin

	for i := 1; i <= numUes; i++ {
		go registerSingleUe2(cfg, wg, msin, i)
		time.Sleep(time.Duration(delayUes) * time.Millisecond)
	}

	wg.Wait()
}

var lock sync.Mutex

func registerSingleUe(cfg config.Config, wg sync.WaitGroup, msin string, i int) {
	imsi := imsiGenerator2(i, msin)
	log.Info("[TESTER] TESTING REGISTRATION USING IMSI ", imsi, " UE")
	cfg.Ue.Msin = imsi
	go ue.RegistrationUe(cfg, uint8(i), &wg, i)
	log_time.LogUeTime(i, 0)
	//wg.Add(1)
}

func imsiGenerator2(i int, msin string) string {
	lock.Lock()
	msin_int, err := strconv.Atoi(msin)
	if err != nil {
		defer lock.Unlock()
		log.Fatal("Error in get configuration")
	}
	base := msin_int + (i -1)

	imsi := fmt.Sprintf("%010d", base)
	defer lock.Unlock()
	return imsi
}