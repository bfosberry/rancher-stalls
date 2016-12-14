package haproxy

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

const (
	pidFilename = "/var/run/haproxy.pid"
)

// StartHAProxy starts or reloads haproxy
func StartHAProxy(configFilename string, wg *sync.WaitGroup) error {
	pid := getPid()
	command := "haproxy"

	pidArg := ""
	if pid != "" {
		pidArg = fmt.Sprintf("-sf %s", pid)
	}
	args := fmt.Sprintf("-f %s -p %s %s", configFilename, pidFilename, pidArg)
	cmd := exec.Command(command, strings.Split(args, " ")...)
	if err := cmd.Start(); err != nil {
		return err
	}
	wg.Add(1)
	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Println(err.Error())
		}
		wg.Done()
	}()
	return nil
}

func getPid() string {
	out, _ := exec.Command("cat", pidFilename).Output()
	return string(out)
}
