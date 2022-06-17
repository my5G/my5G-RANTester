package state

import (
	"my5G-RANTester/internal/control_test_engine/ue/context"
	data "my5G-RANTester/internal/control_test_engine/ue/data/service"
	"my5G-RANTester/internal/control_test_engine/ue/nas"

	"fmt"
	"github.com/bradhe/stopwatch"
	log "github.com/sirupsen/logrus"
)

func DispatchState(ue *context.UEContext, message []byte) {
	DispatchState2(ue, message, nil)
}

func DispatchState2(ue *context.UEContext, message []byte, watch stopwatch.Watch) {

	// if state is PDU session inactive send to analyze NAS
	var state = ue.GetStateSM()

	switch state {

	case context.SM5G_PDU_SESSION_INACTIVE:
		nas.DispatchNas(ue, message)
	case context.SM5G_PDU_SESSION_ACTIVE_PENDING:
		nas.DispatchNas(ue, message)
	case context.SM5G_PDU_SESSION_ACTIVE:
		data.InitDataPlane(ue, message)
	}

	// go LogTime(state, watch)
}

func LogTime(state int, watch stopwatch.Watch) {

	if watch != nil {
		var stateStr string

		switch state {

		case context.SM5G_PDU_SESSION_INACTIVE:
			stateStr = "SM5G_PDU_SESSION_INACTIVE"
		case context.SM5G_PDU_SESSION_ACTIVE_PENDING:
			stateStr = "SM5G_PDU_SESSION_ACTIVE_PENDING"
		case context.SM5G_PDU_SESSION_ACTIVE:
			stateStr = "SM5G_PDU_SESSION_ACTIVE"
		}

		watch.Stop()
		log.Info(fmt.Sprintf("[UE][%s] Milliseconds elapsed: %v", stateStr, watch.Milliseconds()))
	}
}