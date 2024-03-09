package ue

import (
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/ue/context"
	"my5G-RANTester/internal/control_test_engine/ue/nas/service"
	"my5G-RANTester/internal/control_test_engine/ue/nas/trigger"
	"my5G-RANTester/internal/monitoring"
	"my5G-RANTester/lib/nas/security"
	"os"
	"os/signal"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func RegistrationUe(conf config.Config, id uint8, wg *sync.WaitGroup, ueRegistrationSignal chan int) {

	// wg := sync.WaitGroup{}

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
		id,
		conf.GNodeB.PlmnList.GnbId)

	// starting communication with GNB and listen.
	ueTerminationSignal := make(chan int, 1)

	err := service.InitConn(ue, ueRegistrationSignal, ueTerminationSignal)
	if err != nil {
		log.Warn("Error in ", err)
		wg.Done()
		ueRegistrationSignal <- 0
		return
	} else {
		log.Info("[UE] UNIX/NAS service is running")
		// wg.Add(1)
	}

	// registration procedure started.
	trigger.InitRegistration(ue)

	// wg.Wait()

	// control the signals
	sigUe := make(chan os.Signal, 1)
	signal.Notify(sigUe, os.Interrupt)

	// Block until a signal is received.
	select {
	case <-sigUe:
		ue.Terminate()
		wg.Done()
	case <-ueTerminationSignal:
		ueRegistrationSignal <- 0
		ue.Terminate()
		wg.Done()
	case <-time.After(60 * time.Second):
		wg.Done()
	}
	// <- sigUe
	// ue.Terminate()
	// wg.Done()
	// os.Exit(0)

}

func RegistrationUeMonitor(conf config.Config,
	id uint8, monitor *monitoring.Monitor, wg *sync.WaitGroup, start time.Time, ueRegistrationSignal chan int) {

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
		id,
		conf.GNodeB.PlmnList.GnbId)

	// starting communication with GNB and listen.
	ueTerminationSignal := make(chan int, 1)
	err := service.InitConn(ue, ueRegistrationSignal, ueTerminationSignal)
	if err != nil {
		log.Fatal("Error in", err)
	} else {
		log.Info("[UE] UNIX/NAS service is running")
		// wg.Add(1)
	}

	// registration procedure started.
	trigger.InitRegistration(ue)

	for {

		// UE is register in network
		if ue.GetStateMM() == 0x03 {
			elapsed := time.Since(start)
			monitor.LtRegisterLocal = elapsed.Milliseconds()
			log.Warn("[TESTER][UE] UE LATENCY IN REGISTRATION: ", monitor.LtRegisterLocal, " ms")
			break
		}

		// timeout is 10 000 ms
		if time.Since(start).Milliseconds() >= 10000 {
			log.Warn("[TESTER][UE] TIME EXPIRED IN UE REGISTRATION 10 000 ms")
			break
		}
	}

	wg.Done()
	// ue.Terminate()
	// os.Exit(0)
}
