package trigger

import (
	"my5G-RANTester/internal/control_test_engine/ue/context"
	"my5G-RANTester/internal/control_test_engine/ue/nas/message/nas_control/mm_5gs"
	"my5G-RANTester/internal/control_test_engine/ue/nas/message/sender"
	"my5G-RANTester/lib/nas/nasMessage"

	"my5G-RANTester/internal/analytics/log_time"
)

func InitRegistration(ue *context.UEContext) {
	// registration procedure started.
	registrationRequest := mm_5gs.GetRegistrationRequest(
		nasMessage.RegistrationType5GSInitialRegistration,
		nil,
		nil,
		false,
		ue)

	// send to GNB.
	sender.SendToGnb(ue, registrationRequest)
	log_time.LogUeTime(ue.ue_id, 2)

	// change the state of ue for deregistered
	ue.SetStateMM_DEREGISTERED()
}
