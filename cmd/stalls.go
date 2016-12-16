package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/bfosberry/rancher-stalls/haproxy"
	"github.com/bfosberry/rancher-stalls/services"
	"github.com/rancher/go-rancher-metadata/metadata"
)

const (
	metadataURL          = "http://rancher-metadata/2015-12-19"
	serviceNameEnvVar    = "SERVICE_NAME"
	basePortEnvVar       = "BASE_PORT"
	backendPortEnvVar    = "BACKEND_PORT"
	pollSleepTimeEnvVar  = "POLL_SLEEP_TIME_MS"
	configFilename       = "/haproxy.cfg"
	defaultPollSleepTime = "5000"
)

type configFetcher func() (string, error)

func main() {
	metadataClient := metadata.NewClient(metadataURL)
	servicesClient := services.NewServicesClient(metadataClient)

	serviceName, backendPort, basePort, pollSleepTime := loadConfig()

	fmt.Printf("Starting HAProxy with %s, %d, %d, %d\n", serviceName, backendPort, basePort, pollSleepTime)

	fetcher := func() (string, error) {
		return getHaproxyConfig(servicesClient, serviceName, backendPort, basePort)
	}

	haproxyWg := &sync.WaitGroup{}
	if err := updateHaproxy(fetcher, haproxyWg); err != nil {
		fail(err)
	}

	metadataClient.OnChange(pollSleepTime/1000, func(_ string) {
		if err := updateHaproxy(fetcher, haproxyWg); err != nil {
			fail(err)
		}
	})

	for {
		time.Sleep(1 * time.Minute)
	}
}

func updateHaproxy(fetcher configFetcher, wg *sync.WaitGroup) error {
	config, err := fetcher()
	if err != nil {
		return err
	}

	if err := writeHaproxyConfig(config, configFilename); err != nil {
		return err
	}

	return haproxy.StartHAProxy(configFilename, wg)
}

func writeHaproxyConfig(config, filename string) error {
	return ioutil.WriteFile(filename, []byte(config), os.ModePerm)
}

func getHaproxyConfig(servicesClient services.Client, serviceName string, backendPort, basePort int) (string, error) {
	service, err := servicesClient.GetService(serviceName, basePort)
	if err != nil {
		return "", err
	}
	config, err := haproxy.GenerateConfig(service, backendPort)
	if err != nil {
		return "", err
	}
	return config, nil
}

func loadConfig() (string, int, int, int) {
	serviceName := loadRequiredEnvVar(serviceNameEnvVar)
	backendPort := loadRequiredEnvVar(backendPortEnvVar)
	backendPortInt, err := strconv.Atoi(backendPort)
	if err != nil {
		fail(err)
	}

	basePort := loadRequiredEnvVar(basePortEnvVar)
	basePortInt, err := strconv.Atoi(basePort)
	if err != nil {
		fail(err)
	}

	pollSleepTime := loadOptionalEnvVar(pollSleepTimeEnvVar, defaultPollSleepTime)
	pollSleepTimeInt, err := strconv.Atoi(pollSleepTime)
	if err != nil {
		fail(err)
	}
	return serviceName, backendPortInt, basePortInt, pollSleepTimeInt
}

func loadRequiredEnvVar(key string) string {
	envVar := os.Getenv(key)
	if envVar == "" {
		fail(fmt.Errorf("Missing required env var %s", key))
	}
	return envVar
}

func loadOptionalEnvVar(key string, defaultValue string) string {
	envVar := os.Getenv(key)
	if envVar == "" {
		envVar = defaultValue
	}
	return envVar
}

func fail(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}
