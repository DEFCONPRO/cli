//go:generate go run github.com/travisjeffery/mocker/cmd/mocker --dst ../../../mock/confluent_current.go --pkg mock --selfpkg github.com/confluentinc/cli confluent_current.go ConfluentCurrent

package local

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

/*
Directory Structure:

CONFLUENT_CURRENT/
	confluent.current
	confluent.000000/
		[service]/
			data/
			logs/ # TODO
			[service].config
			[service].log
			[service].pid
*/

type ConfluentCurrent interface {
	RemoveTrackingFile() error

	GetCurrentDir() (string, error)
	RemoveCurrentDir() error

	GetDataDir(service string) (string, error)

	GetConfigFile(service string) (string, error)
	SetConfig(service string, config []byte) error

	GetLogFile(service string) (string, error)

	GetPidFile(service string) (string, error)
	HasPidFile(service string) (bool, error)
	GetPid(service string) (int, error)
	SetPid(service string, pid int) error
	RemovePidFile(service string) error
}

type ConfluentCurrentManager struct {
	trackingFile string
	currentDir   string
	pidFiles     map[string]string
}

func NewConfluentCurrentManager() *ConfluentCurrentManager {
	cc := new(ConfluentCurrentManager)
	cc.pidFiles = make(map[string]string)
	return cc
}

func (cc *ConfluentCurrentManager) GetCurrentDir() (string, error) {
	if cc.currentDir != "" {
		return cc.currentDir, nil
	}

	if !fileExists(cc.getTrackingFile()) {
		cc.currentDir = getRandomChildDir(cc.getRootDir())
		if err := os.MkdirAll(cc.currentDir, 0777); err != nil {
			return "", err
		}
		if err := ioutil.WriteFile(cc.getTrackingFile(), []byte(cc.currentDir), 0644); err != nil {
			return "", err
		}
	} else {
		data, err := ioutil.ReadFile(cc.getTrackingFile())
		if err != nil {
			return "", err
		}
		cc.currentDir = string(data)
	}

	return cc.currentDir, nil
}

func (cc *ConfluentCurrentManager) getServiceDir(service string) (string, error) {
	dir, err := cc.GetCurrentDir()
	if err != nil {
		return "", err
	}

	dir = filepath.Join(dir, service)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return "", err
	}

	return dir, nil
}

func (cc *ConfluentCurrentManager) GetDataDir(service string) (string, error) {
	dir, err := cc.getServiceDir(service)
	if err != nil {
		return "", err
	}

	dir = filepath.Join(dir, "data")
	if service == "ksql-server" {
		// TODO: Investigate if this is actually necessary
		dir = filepath.Join(dir, "kafka-streams")
	}
	if err := os.MkdirAll(dir, 0777); err != nil {
		return "", err
	}

	return dir, nil
}

func (cc *ConfluentCurrentManager) GetConfigFile(service string) (string, error) {
	return cc.getServiceFile(service, fmt.Sprintf("%s.config", service))
}

func (cc *ConfluentCurrentManager) SetConfig(service string, config []byte) error {
	file, err := cc.GetConfigFile(service)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, config, 0644)
}

func (cc *ConfluentCurrentManager) GetLogFile(service string) (string, error) {
	return cc.getServiceFile(service, fmt.Sprintf("%s.log", service))
}

func (cc *ConfluentCurrentManager) HasPidFile(service string) (bool, error) {
	file, err := cc.GetPidFile(service)
	if err != nil {
		return false, err
	}

	return fileExists(file), nil
}

func (cc *ConfluentCurrentManager) GetPidFile(service string) (string, error) {
	if file, ok := cc.pidFiles[service]; ok {
		return file, nil
	}

	file, err := cc.getServiceFile(service, fmt.Sprintf("%s.pid", service))
	if err != nil {
		return "", err
	}

	cc.pidFiles[service] = file
	return cc.pidFiles[service], nil
}

func (cc *ConfluentCurrentManager) GetPid(service string) (int, error) {
	file, err := cc.GetPidFile(service)
	if err != nil {
		return 0, err
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(data))
}

func (cc *ConfluentCurrentManager) SetPid(service string, pid int) error {
	file, err := cc.GetPidFile(service)
	if err != nil {
		return err
	}

	data := []byte(strconv.Itoa(pid))
	return ioutil.WriteFile(file, data, 0644)
}

func (cc *ConfluentCurrentManager) RemoveTrackingFile() error {
	return os.Remove(cc.trackingFile)
}

func (cc *ConfluentCurrentManager) RemoveCurrentDir() error {
	return os.RemoveAll(cc.currentDir)
}

func (cc *ConfluentCurrentManager) RemovePidFile(service string) error {
	return os.Remove(cc.pidFiles[service])
}

func (cc *ConfluentCurrentManager) getRootDir() string {
	if dir := os.Getenv("CONFLUENT_CURRENT"); dir != "" {
		return dir
	}
	return os.TempDir()
}

func (cc *ConfluentCurrentManager) getTrackingFile() string {
	if cc.trackingFile != "" {
		return cc.trackingFile
	}

	cc.trackingFile = filepath.Join(cc.getRootDir(), "confluent.current")
	return cc.trackingFile
}

func (cc *ConfluentCurrentManager) getServiceFile(service, file string) (string, error) {
	dir, err := cc.getServiceDir(service)
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, file), nil
}

func getRandomChildDir(parentDir string) string {
	rand.Seed(time.Now().Unix())

	for {
		childDir := fmt.Sprintf("confluent.%06d", rand.Intn(1000000))
		path := filepath.Join(parentDir, childDir)
		if !fileExists(path) {
			return path
		}
	}
}

func fileExists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}