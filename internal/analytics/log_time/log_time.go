package analytics

import (
	"fmt"
	"time"
	log "github.com/sirupsen/logrus"
)

func LogUeTime(id int, task_id int) {

	nsec_now := now.UnixNano()

	go ShowUeLog(id, task_id, nsec_now)
}

func ShowUeLog(id int, task_id int, time int64) {
	if id > 0 {
		log.Info(fmt.Sprintf("[Lando] %d, %d, %d", id, task_id, time))
	}
}