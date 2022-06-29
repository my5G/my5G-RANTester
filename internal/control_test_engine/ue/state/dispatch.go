package state

import (
	"my5G-RANTester/internal/control_test_engine/ue/context"
	data "my5G-RANTester/internal/control_test_engine/ue/data/service"
	"my5G-RANTester/internal/control_test_engine/ue/nas"

	log_time "my5G-RANTester/internal/analytics/log_time"
)

func DispatchState(ue *context.UEContext, message []byte) {

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

	log_time.LogUeTime(ue.GetUeId2(), state)
}
