package config

import (
	_ "embed"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

const (
	vendorName       = "my5g"
	configDir        = "RANTester"
	configName       = "config.yml"
	workingDirectory = "."
)

//go:embed config.yml
var defaultConfigData []byte

func getDefaultConfigDir() (string, error) {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		cfgDir = getDirOfRunningBinary()
	}

	target := filepath.Join(cfgDir, vendorName, configDir)

	return filepath.Abs(target)
}

func getDefaultConfigPath() (string, error) {
	parentDir, err := getDefaultConfigDir()
	if err != nil {
		parentDir = workingDirectory
	}

	return filepath.Join(parentDir, configName), nil
}

func getDirOfRunningBinary() string {
	_, b, _, _ := runtime.Caller(0)

	return filepath.Dir(path.Join(path.Dir(b)))
}

// Load attempts to load the given config paths, yielding the first valid path, or a default config file if
// no user-specified paths are found.
func Load(paths ...string) (*Config, error) {
	toCheck := make([]string, 0)

	for _, path := range paths {
		if path != "" {
			toCheck = append(toCheck, path)
		}
	}

	defaultConfigPath, errGetDefault := getDefaultConfigPath()
	if errGetDefault == nil && len(toCheck) < 1 {
		toCheck = append(toCheck, defaultConfigPath)
	}

	var target string

	for _, cfgPath := range toCheck {
		if _, errFound := os.Stat(cfgPath); errors.Is(errFound, os.ErrNotExist) {
			log.Warnf("could not find config file '%v'", cfgPath)
			continue // nothing there, go to next one
		}

		target = cfgPath
		break // no need to look any further
	}

	// if the above loop did not find a valid config file,
	// create the default config file and set it as the target
	if target == "" {
		target = defaultConfigPath
		if errMake := makeConfig(target); errMake != nil {
			return nil, errMake
		}
	}

	log.Infof("using config file at: %v", target)

	file, errOpen := ioutil.ReadFile(target)
	if errOpen != nil {
		return nil, fmt.Errorf("could not read config file '%v': %v", target, errOpen)
	}

	cfg := &Config{}
	if errUnmarshal := yaml.Unmarshal(file, &cfg); errUnmarshal != nil {
		return nil, fmt.Errorf("could not unmarshal file data from '%v': %v", target, errUnmarshal)
	}

	return cfg, nil
}

func makeConfig(target string) error {
	if _, err := os.Stat(target); err == nil {
		log.Info("config file exists")
		return nil
	}

	configDirParent := filepath.Dir(target)

	if mkdirErr := os.MkdirAll(configDirParent, os.ModePerm); mkdirErr != nil {
		return fmt.Errorf("could not create default config dir: %v", mkdirErr)
	}

	log.Infof("created config directory: %v", configDirParent)

	writeErr := os.WriteFile(target, defaultConfigData, 0755)
	if writeErr != nil {
		const fmtWriteErr = "no config file found, also could not write default config file '%v': %v"
		return fmt.Errorf(fmtWriteErr, target, writeErr)
	}

	log.Infof("created config file: %v", target)

	return nil
}
