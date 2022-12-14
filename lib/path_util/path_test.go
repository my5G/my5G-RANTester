package path_util

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestFree5gcPath(t *testing.T) {
	log.Infoln(Gofree5gcPath("free5gc/abcdef/abcdef.pem"))
}
