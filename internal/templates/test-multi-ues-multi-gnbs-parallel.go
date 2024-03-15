package templates

import (
	"math/rand"

	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"strconv"

	"my5G-RANTester/internal/control_test_engine/ue"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestMultiUesMultiGNBsParallel(numUes int, numGNBs int) {

	log.Info("Num of UEs = ", numUes)
	log.Info("Num of GNBs = ", numGNBs)

	wg := sync.WaitGroup{}

	cfg, err := config.GetConfig()
	if err != nil {
		//return nil
		log.Fatal("Error in get configuration")
	}

	gnbControlPort := cfg.GNodeB.ControlIF.Port
	gnbID, err := strconv.Atoi(string(cfg.GNodeB.PlmnList.GnbId))
	if err != nil {
		log.Error("Failed to extract gnbID")
	}

	baseGnbID := gnbID

	for i := 0; i < numGNBs; i++ {
		//newGnbID := fmt.Sprintf("%d", gnbID)
		newGnbID := constructGnbID(gnbID)

		cfg.GNodeB.PlmnList.GnbId = newGnbID
		cfg.GNodeB.ControlIF.Port = gnbControlPort + i
		log.Info("Initializing gnb with GnbId = ", cfg.GNodeB.PlmnList.GnbId)
		log.Info("Initializing gnb with gnbControlPort = ", cfg.GNodeB.ControlIF.Port)

		go gnb.InitGnb(cfg, &wg)
		wg.Add(1)
		time.Sleep(1 * time.Second)

		gnbID++
	}

	// time.Sleep(1 * time.Second)

	ueRegistrationSignal := make(chan int, numUes)

	msin := cfg.Ue.Msin
	var ueSessionIdOffset uint8 = 0
	var ueSessionId uint8
	// startTime := time.Now()
	for i := 1; i <= numUes; i++ {

		offset := rand.Intn(numGNBs)
		cfg.GNodeB.ControlIF.Port = gnbControlPort + offset
		cfg.GNodeB.PlmnList.GnbId = constructGnbID(baseGnbID + offset)

		log.Info("Registering ue with gnbControlPort = ", cfg.GNodeB.ControlIF.Port)

		imsi := imsiGenerator(i, msin)
		log.Info("[TESTER] TESTING REGISTRATION USING IMSI ", imsi, " UE")
		cfg.Ue.Msin = imsi

		if (i+int(ueSessionIdOffset))%256 == 0 {
			ueSessionIdOffset++
		}
		ueSessionId = uint8(i) + ueSessionIdOffset

		go ue.RegistrationUe(cfg, ueSessionId, &wg, ueRegistrationSignal)
		wg.Add(1)

		select {
		case <-ueRegistrationSignal:
		default:
		}

		sleepTime := 200 * time.Millisecond
		time.Sleep(sleepTime)
	}

	wg.Wait()
	// endTime := time.Now()
	// executionTime := endTime.Sub(startTime)
	// log.Info("Total Registeration Time =", executionTime)
}
