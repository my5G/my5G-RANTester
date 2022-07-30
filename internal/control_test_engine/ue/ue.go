package ue

import (
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/ue/context"
	"my5G-RANTester/internal/control_test_engine/ue/nas/service"
	"my5G-RANTester/internal/control_test_engine/ue/nas/trigger"
	"my5G-RANTester/lib/nas/security"
	"os"
	"os/signal"
	"sync"
	
	"github.com/gookit/event"
)

func RegistrationUe(conf config.Config, id int64, wg *sync.WaitGroup) {
	RegistrationUe2(conf, id, wg, -1)
}

func RegistrationUe2(conf config.Config, id int64, wg *sync.WaitGroup, delayDsc int) {

	// new UE instance.
	ue := &context.UEContext{}

	// new UE context
	ue.NewRanUeContext(
		conf.Ue.Msin,
		security.AlgCiphering128NEA0,
		security.AlgIntegrity128NIA2,
		conf.Ue.Key,
		conf.Ue.Opc,
		"c9e8763286b5b9ffbdf56e1297d0887b",
		conf.Ue.Amf,
		conf.Ue.Sqn,
		conf.Ue.Hplmn.Mcc,
		conf.Ue.Hplmn.Mnc,
		conf.Ue.Dnn,
		int32(conf.Ue.Snssai.Sst),
		conf.Ue.Snssai.Sd,
		id)

	// In case the disconnection delay is different of -1 (it's enabled),
	// listen for a disconnection event
	running := true
	if delayDsc != -1 {
		event.On(ue.GetMsin(), event.ListenerFunc(func(e event.Event) error {
			time.Sleep(time.Duration(delayDsc) * time.Millisecond)
			log_time.LogUeTime(0, ue.GetMsin(), "StartDeregistration")
			ue.Terminate()
			wg.Done()
	
			ue = nil // Clear UE pointer
			running = false
			return nil
		}))
	}

	// starting communication with GNB and listen.
	err := service.InitConn(ue)
	if err != nil {
		log.Fatal("Error in", err)
	} else {
		log.Info("[UE] UNIX/NAS service is running")
		// wg.Add(1)
	}

	// registration procedure started.
	trigger.InitRegistration(ue)

	if delayDsc != -1 {
		// Wait until finishes
		for running != false {
			time.Sleep(time.Duration(5) * time.Millisecond)
		}
	}
	else {
		// Use a signal to verify when it needs to disconnect

		// control the signals
		sigUe := make(chan os.Signal, 1)
		signal.Notify(sigUe, os.Interrupt)

		// Block until a signal is received.
		<-sigUe
		ue.Terminate()
		wg.Done()
		// os.Exit(0)
	}
}
