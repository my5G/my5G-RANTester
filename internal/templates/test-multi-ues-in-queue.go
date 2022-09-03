package templates

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"my5G-RANTester/internal/control_test_engine/ue"
	"strconv"
	"sync"
)

func TestMultiUesInQueue(numUes int) {

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

	// wait AMF get active state
	<-synchGnb

	msin := cfg.Ue.Msin

	for i := 1; i <= numUes; i++ {

		// synch UE
		synchUE := make(chan bool, 1)

		imsi := imsiGenerator(i, msin)

		log.Info("[TESTER] TESTING REGISTRATION USING IMSI ", imsi, " UE")

		cfg.Ue.Msin = imsi

		go ue.RegistrationUe(cfg, uint8(i), &wg, synchUE)

		wg.Add(1)

		// wait the UE establish PDU Session
		<-synchUE
	}

	wg.Wait()
}

func imsiGenerator(i int, msin string) string {

	msin_int, err := strconv.Atoi(msin)
	if err != nil {
		log.Fatal("Error in get configuration")
	}
	base := msin_int + (i - 1)

	imsi := fmt.Sprintf("%010d", base)
	return imsi
}
