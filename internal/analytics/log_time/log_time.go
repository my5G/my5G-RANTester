package analytics

import (
	"fmt"
	"time"
	log "github.com/sirupsen/logrus"
)

var enabled = true
func ChangeAnalyticsState(state bool){
	enabled = state
}

var gnbid = 0
func SetGnodebId(id int){
	gnbid = id
}

func LogUeTime(id string, task string) {

	now := time.Now()

	go ShowUeLog(id, task, now)
}

func ShowUeLog(id string, task string, now time.Time) {
	if enabled {
		nsec_now := now.UnixNano()
		log.Info(fmt.Sprintf("[ANALYTICS] %d, %s, %s, %d", gnbid, id, task, nsec_now))
	}
}