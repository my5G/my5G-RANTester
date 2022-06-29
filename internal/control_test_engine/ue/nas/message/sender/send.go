package sender

import (
	"fmt"
	"my5G-RANTester/internal/control_test_engine/ue/context"

	"my5G-RANTester/internal/analytics/log_time"
)

func SendToGnb(ue *context.UEContext, message []byte) {

	conn := ue.GetUnixConn()
	
	_, err := conn.Write(message)
	if err != nil {
		fmt.Println("Tratar o erro")
	}
	
	log_time.LogUeTime(ue.ue_id, 3)
}
