package trigger

import (
	"my5G-RANTester/internal/control_test_engine/ue/context"
	"my5G-RANTester/internal/control_test_engine/ue/nas/message/nas_control/mm_5gs"
	"my5G-RANTester/internal/control_test_engine/ue/nas/message/sender"
	"my5G-RANTester/lib/nas/nasMessage"

	"my5G-RANTester/internal/analytics/log_time"
)

func InitRegistration(ue *context.UEContext) {
	InitRegistration(ue, 0)
}

func InitRegistration(ue *context.UEContext, ue_id int) {

	// registration procedure started.
	registrationRequest := mm_5gs.GetRegistrationRequest(
		nasMessage.RegistrationType5GSInitialRegistration,
		nil,
		nil,
		false,
		ue)

	// send to GNB.
	sender.SendToGnb(ue, registrationRequest, ue_id)
	log_time.LogUeTime(ue_id, 2)

	// change the state of ue for deregistered
	ue.SetStateMM_DEREGISTERED()
}
