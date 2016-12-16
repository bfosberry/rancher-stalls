package command_test

import (
	"testing"
	"time"

	"github.com/bfosberry/rancher-stalls/command"
	"github.com/stretchr/testify/assert"
)

func TestExecRunner(t *testing.T) {
	var runner command.Runner = command.ExecRunner
	output, err := runner("printenv", []string{}, true, nil)
	assert.NotEqual(t, "", output)
	assert.Nil(t, err)
}

func TestExecRunnerError(t *testing.T) {
	var runner command.Runner = command.ExecRunner
	output, err := runner("hoobastank", []string{}, true, nil)
	assert.Equal(t, "", output)
	assert.NotNil(t, err)
}

func TestExecRunnerBackground(t *testing.T) {
	var runner command.Runner = command.ExecRunner
	doneChan := make(chan bool)
	output, err := runner("printenv", []string{}, false, func(err error) {
		doneChan <- true
	})
	assert.Equal(t, "", output)
	assert.Nil(t, err)
	select {
	case <-doneChan:
	case <-time.After(1 * time.Second):
		t.Error("Timed out waiting for callback")
	}
}

func TestExecRunnerBackgroundLongRunning(t *testing.T) {
	var runner command.Runner = command.ExecRunner
	doneChan := make(chan bool)
	output, err := runner("sleep", []string{"1s"}, false, func(err error) {
		doneChan <- true
	})
	assert.Equal(t, "", output)
	assert.Nil(t, err)
	select {
	case <-doneChan:
		t.Error("Received callback too soon")
	case <-time.After(10 * time.Millisecond):
		return
	}
}

func TestExecRunnerBackgroundError(t *testing.T) {
	var runner command.Runner = command.ExecRunner
	doneChan := make(chan bool)
	output, err := runner("hoobastank", []string{}, false, func(err error) {
		doneChan <- true
	})
	assert.Equal(t, "", output)
	assert.NotNil(t, err)
	select {
	case <-doneChan:
		t.Error("Received callback incorrectly")
	case <-time.After(10 * time.Millisecond):
	}
}
