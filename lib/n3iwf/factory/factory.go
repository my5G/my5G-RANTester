/*
 * N3IWF Configuration Factory
 */

package factory

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

var N3iwfConfig Config

func checkErr(err error) {
	if err != nil {
		err = fmt.Errorf("[Configuration] %s", err.Error())
		log.Fatal(err)
	}
}

// TODO: Support configuration update from REST api
func InitConfigFactory(f string) {
	content, err := ioutil.ReadFile(f)
	checkErr(err)

	N3iwfConfig = Config{}

	err = yaml.Unmarshal([]byte(content), &N3iwfConfig)
	checkErr(err)

	log.Infof("Successfully initialize configuration %s", f)
}
