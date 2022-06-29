package trigger

import (
	"my5G-RANTester/internal/control_test_engine/ue/context"
	"my5G-RANTester/internal/control_test_engine/ue/nas/message/nas_control/mm_5gs"
	"my5G-RANTester/internal/control_test_engine/ue/nas/message/sender"
	"my5G-RANTester/lib/nas/nasMessage"

	log_time "my5G-RANTester/internal/analytics/log_time"
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
	log_time.LogUeTime(ue.GetUeId2(), 2)

	// change the state of ue for deregistered
	ue.SetStateMM_DEREGISTERED()
}
