package templates

import (
	"math/rand"

	"fmt"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"strconv"

	"my5G-RANTester/internal/control_test_engine/ue"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestMultiUesMultiGNBs(numUes int, numGNBs int) {

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

	msin := cfg.Ue.Msin
	// startTime := time.Now()
	for i := 1; i <= numUes; i++ {

		offset := rand.Intn(numGNBs)
		cfg.GNodeB.ControlIF.Port = gnbControlPort + offset
		cfg.GNodeB.PlmnList.GnbId = constructGnbID(baseGnbID + offset)

		log.Info("Registering ue with gnbControlPort = ", cfg.GNodeB.ControlIF.Port)

		imsi := imsiGenerator(i, msin)
		log.Info("[TESTER] TESTING REGISTRATION USING IMSI ", imsi, " UE")
		cfg.Ue.Msin = imsi
		go ue.RegistrationUe(cfg, uint8(i), &wg)
		wg.Add(1)

		// sleepTime := time.Duration(rand.Intn(100)+1) * time.Millisecond
		sleepTime := 500 * time.Millisecond
		time.Sleep(sleepTime)
	}

	wg.Wait()
	// endTime := time.Now()
	// executionTime := endTime.Sub(startTime)
	// log.Info("Total Registeration Time =", executionTime)
}

func constructGnbID(gnbID int) string {
	var newGnbID string

	if gnbID <= 9 {
		newGnbID = fmt.Sprintf("00000%d", gnbID)
	} else if gnbID > 9 && gnbID <= 99 {
		newGnbID = fmt.Sprintf("0000%d", gnbID)
	} else if gnbID > 99 && gnbID <= 999 {
		newGnbID = fmt.Sprintf("000%d", gnbID)
	}

	return newGnbID
}
