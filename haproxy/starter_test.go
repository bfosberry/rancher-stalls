package haproxy_test

import (
	"sync"
	"testing"

	"github.com/bfosberry/rancher-stalls/command"
	"github.com/bfosberry/rancher-stalls/haproxy"
	"github.com/stretchr/testify/assert"
)

const (
	pidFilename    = "/var/run/haproxy.pid"
	configFilename = "/haproxy.cfg"
	haproxyPath    = "/usr/local/sbin/haproxy"
)

func TestStartHAProxy(t *testing.T) {
	commandChan := make(chan string, 2)
	argsChan := make(chan []string, 2)
	runner := command.NewDummyRunner(commandChan, argsChan, "", nil, nil)

	haproxy.Runner = runner
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		assert.Nil(t, haproxy.StartHAProxy(configFilename))
		wg.Done()
	}()

	assert.Equal(t, "cat", <-commandChan)
	assert.Equal(t, []string{pidFilename}, <-argsChan)
	assert.Equal(t, haproxyPath, <-commandChan)
	assert.Equal(t, []string{"-D", "-f", configFilename, "-p", pidFilename}, <-argsChan)
}
