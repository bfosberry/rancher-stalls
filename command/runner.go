package command

import (
	"bytes"
	"os/exec"
)

// DoneFunc indicates an action is complete with a possible error
type DoneFunc func(error)

// Runner takes a command and args, executes that command
// and triggers an onDone callback when complete
type Runner func(command string, args []string, inline bool, onDone DoneFunc) (string, error)

// ExecRunner provides a Runner which uses the exec lib
func ExecRunner(command string, args []string, inline bool, onDone DoneFunc) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if inline {
		err := cmd.Run()
		return out.String(), err
	}
	if err := cmd.Start(); err != nil {
		return out.String(), err
	}
	if onDone != nil {
		go func(c *exec.Cmd) {
			onDone(cmd.Wait())
		}(cmd)
	}
	return out.String(), nil
}

// NewDummyRunner returns a Runner which can be used for testing
func NewDummyRunner(commandPtr *string, argsPtr *[]string, output string, err error, doneErr error) Runner {
	return func(command string, args []string, inline bool, onDone DoneFunc) (string, error) {
		*commandPtr = command
		*argsPtr = args
		if onDone != nil {
			go onDone(doneErr)
		}
		return output, err
	}
}
