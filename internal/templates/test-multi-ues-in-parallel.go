package templates

import (
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"my5G-RANTester/internal/control_test_engine/ue"
	"sync"
	"time"
	
	log_time "my5G-RANTester/internal/analytics/log_time"
)

func TestMultiUesInParallel(numUes int, delayUes int, delayStart int, delayDsc int, showAnalytics bool) {

	wg := sync.WaitGroup{}

	cfg, err := config.GetConfig()
	if err != nil {
		//return nil
		log.Fatal("Error in get configuration")
	}

	log_time.ChangeAnalyticsState(showAnalytics)

	go gnb.InitGnb(cfg, &wg)

	wg.Add(1)

	time.Sleep(time.Duration(delayStart) * time.Second)
    msin :=  cfg.Ue.Msin

	for i := 1; i <= numUes; i++ {
		go registerSingleUe(cfg, wg, msin, i, delayDsc)
		time.Sleep(time.Duration(delayUes) * time.Millisecond)
	}

	wg.Wait()
}

func registerSingleUe(cfg config.Config, wg sync.WaitGroup, msin string, i int, delayDsc int) {
	imsi := imsiGenerator(i, msin)
	log.Info("[TESTER] TESTING REGISTRATION USING IMSI ", imsi, " UE")
	cfg.Ue.Msin = imsi
	log_time.LogUeTime(0, imsi, "StartRegistration")
	go ue.RegistrationUe2(cfg, int64(i), &wg, delayDsc)
	//wg.Add(1)
}
