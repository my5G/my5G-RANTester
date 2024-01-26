package templates

import (
	log "github.com/sirupsen/logrus"
)

func TestMultiUesMultiGNBs(numUes int, numGNBs int) {

	log.Info("Num of UEs =", numUes)
	log.Info("Num of GNBs =", numGNBs)

	// wg := sync.WaitGroup{}

	// cfg, err := config.GetConfig()
	// if err != nil {
	// 	//return nil
	// 	log.Fatal("Error in get configuration")
	// }

	// go gnb.InitGnb(cfg, &wg)
	// wg.Add(1)

	// time.Sleep(1 * time.Second)

	// msin := cfg.Ue.Msin
	// for i := 1; i <= numUes; i++ {

	// 	imsi := imsiGenerator(i, msin)
	// 	log.Info("[TESTER] TESTING REGISTRATION USING IMSI ", imsi, " UE")
	// 	cfg.Ue.Msin = imsi
	// 	go ue.RegistrationUe(cfg, uint8(i), &wg)
	// 	wg.Add(1)

	// 	time.Sleep(10 * time.Second)
	// }

	// wg.Wait()
}
