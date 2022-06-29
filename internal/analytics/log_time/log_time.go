package analytics

import (
	"fmt"
	"time"
	log "github.com/sirupsen/logrus"
)

func LogUeTime(id int, task_id int) {

	now := time.Now()

	go ShowUeLog(id, task_id, now)
}

func ShowUeLog(id int, task_id int, now time.Time) {
	if id > 0 {
		nsec_now := now.UnixNano()
		log.Info(fmt.Sprintf("[Lando] %d, %d, %d", id, task_id, nsec_now))
	}
}