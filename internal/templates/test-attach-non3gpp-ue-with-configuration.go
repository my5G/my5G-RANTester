package templates

import (
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/non3gppue"
	"sync"

	log "github.com/sirupsen/logrus"
)

func TestAttachNon3gppUeWithConfiguration() {

	wg := sync.WaitGroup{}

	cfg, err := config.GetConfig()
	if err != nil {
		//return nil
		log.Fatal("Error in get configuration")
	}

	go non3gppue.RegistrationUe(cfg, 1, &wg)

	wg.Add(1)

	wg.Wait()
}
