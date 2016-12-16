package haproxy

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bfosberry/rancher-stalls/command"
)

var (
	// Runner holds the global command executor used for this package
	Runner command.Runner = command.ExecRunner
)

const (
	pidFilename = "/var/run/haproxy.pid"
)

// StartHAProxy starts or reloads haproxy
func StartHAProxy(configFilename string, wg *sync.WaitGroup) error {
	pid := getPid()
	command := "/usr/local/sbin/haproxy"

	pidArg := ""
	if pid != "" {
		pidArg = fmt.Sprintf(" -sf %s", pid)
	}
	argsStr := fmt.Sprintf("-D -f %s -p %s%s", configFilename, pidFilename, pidArg)
	_, err := Runner(command, strings.Split(argsStr, " "), true, nil)
	return err
}

func getPid() string {
	out, err := Runner("cat", []string{pidFilename}, true, nil)
	if err != nil {
		out = ""
	}
	return out
}
