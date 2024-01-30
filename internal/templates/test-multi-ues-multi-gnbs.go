package templates

import (
	// "math/rand"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	// "my5G-RANTester/internal/control_test_engine/ue"
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

	// Modify gNB port
	// gnbControlPort := cfg.GNodeB.ControlIF.Port
	// gnbDataPort := cfg.GNodeB.DataIF.Port

	// log.Info("Initial gnbControlIP = ", gnbControlPort)
	// log.Info("Initial gnbDataIP = ", gnbDataPort)
	// log.Info("Initial gnbControlPort = ", gnbControlPort)
	// log.Info("Initial gnbDataPort = ", gnbDataPort)

	// log.Info("Initializing gnb with gnbControlIP = ", cfg.GNodeB.ControlIF.Ip)
	// log.Info("Initializing gnb with gnbDataIP = ", cfg.GNodeB.DataIF.Ip)
	// log.Info("Initializing gnb with gnbControlPort = ", cfg.GNodeB.ControlIF.Port)
	// log.Info("Initializing gnb with gnbDataPort = ", cfg.GNodeB.DataIF.Port)

	cfg.GNodeB.ControlIF.Port = 9488
	cfg.GNodeB.DataIF.Port = 2153

	go gnb.InitGnb(cfg, &wg)
	wg.Add(1)
	time.Sleep(1 * time.Second)

	// cfg.GNodeB.ControlIF.Ip = "172.23.0.250"
	// cfg.GNodeB.DataIF.Ip = "172.23.0.250"
	// cfg.GNodeB.ControlIF.Port = 9487
	// cfg.GNodeB.DataIF.Port = 2152
	// log.Info("Initializing gnb with gnbControlIP = ", cfg.GNodeB.ControlIF.Ip)
	// log.Info("Initializing gnb with gnbDataIP = ", cfg.GNodeB.DataIF.Ip)
	// log.Info("Initializing gnb with gnbControlPort = ", cfg.GNodeB.ControlIF.Port)
	// log.Info("Initializing gnb with gnbDataPort = ", cfg.GNodeB.DataIF.Port)

	// go gnb.InitGnb(cfg, &wg)
	// wg.Add(1)
	// time.Sleep(1 * time.Second)

	//time.Sleep(1 * time.Second)

	// msin := cfg.Ue.Msin
	//startTime := time.Now()
	// for i := 1; i <= numUes; i++ {

	// 	portOffset := rand.Intn(numGNBs)
	// 	cfg.GNodeB.ControlIF.Port = gnbControlPort + portOffset
	// 	cfg.GNodeB.DataIF.Port = gnbDataPort + portOffset

	// 	log.Info("Registering ue with gnbControlPort = ", cfg.GNodeB.ControlIF.Port)
	// 	log.Info("Registering ue with gnbDataPort = ", cfg.GNodeB.DataIF.Port)

	// 	imsi := imsiGenerator(i, msin)
	// 	log.Info("[TESTER] TESTING REGISTRATION USING IMSI ", imsi, " UE")
	// 	cfg.Ue.Msin = imsi
	// 	go ue.RegistrationUe(cfg, uint8(i), &wg)
	// 	wg.Add(1)

	// 	sleepTime := time.Duration(rand.Intn(100)+1) * time.Millisecond
	// 	time.Sleep(sleepTime)
	// }

	wg.Wait()
	// endTime := time.Now()
	// executionTime := endTime.Sub(startTime)
	// log.Info("Total Registeration Time =", executionTime)
}
