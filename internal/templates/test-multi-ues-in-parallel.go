package templates

import (
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"my5G-RANTester/internal/control_test_engine/ue"
	"sync"
	"time"
	
	log_time "my5G-RANTester/internal/analytics/log_time"
	"github.com/gookit/event"
	"my5G-RANTester/internal/control_test_engine/ue/context"
)

func TestMultiUesInParallel(numUes int, delayUes int, delayStart int, showAnalytics bool) {

	wg := sync.WaitGroup{}

	cfg, err := config.GetConfig()
	if err != nil {
		//return nil
		log.Fatal("Error in get configuration")
	}

	log_time.ChangeAnalyticsState(showAnalytics)

	go gnb.InitGnb(cfg, &wg)

	wg.Add(1)

	event.On("DataPlaneReady", event.ListenerFunc(onDataPlaneReady))

	time.Sleep(time.Duration(delayStart) * time.Second)
    msin :=  cfg.Ue.Msin

	for i := 1; i <= numUes; i++ {
		go registerSingleUe(cfg, wg, msin, i)
		time.Sleep(time.Duration(delayUes) * time.Millisecond)
	}

	wg.Wait()
}

func registerSingleUe(cfg config.Config, wg sync.WaitGroup, msin string, i int) {
	imsi := imsiGenerator(i, msin)
	log.Info("[TESTER] TESTING REGISTRATION USING IMSI ", imsi, " UE")
	cfg.Ue.Msin = imsi
	log_time.LogUeTime(0, imsi, "StartRegistration")
	go ue.RegistrationUe(cfg, int64(i), &wg)
	//wg.Add(1)
}

func onDataPlaneReady(e event.Event) error {
	ue := e.Get("ue").(*context.UEContext)
	go deregisterSingleUe(ue)
	return nil
}

func deregisterSingleUe(ue *context.UEContext) {
	time.Sleep(time.Duration(50) * time.Millisecond)
	log_time.LogUeTime(0, ue.GetMsin(), "StartDeregistration")
	ue.Terminate()
}
