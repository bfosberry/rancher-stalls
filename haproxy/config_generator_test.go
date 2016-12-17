package haproxy_test

import (
	"testing"

	"github.com/bfosberry/rancher-stalls/haproxy"
	"github.com/bfosberry/rancher-stalls/services"
	"github.com/stretchr/testify/assert"
)

const (
	backendPort    = 2000
	expectedConfig = `
frontend api
    mode tcp
    
    bind *:10000
    acl dest0 dst_port 10000
    use_backend Backend0 if dest0
    
    bind *:10001
    acl dest1 dst_port 10001
    use_backend Backend1 if dest1
    
backend Backend0
    mode tcp
    server Backend0 1.2.3.4:2000 check

backend Backend1
    mode tcp
    server Backend1 1.2.3.5:2000 check
`
)

func TestGenerateConfig(t *testing.T) {
	testService := &services.Service{
		Containers: []services.Container{
			services.Container{
				Index:        0,
				IP:           "1.2.3.4",
				ExternalPort: 10000,
			},
			services.Container{
				Index:        1,
				IP:           "1.2.3.5",
				ExternalPort: 10001,
			},
		},
	}
	config, err := haproxy.GenerateConfig(testService, backendPort)
	assert.Nil(t, err)
	assert.Equal(t, expectedConfig, config)
}
