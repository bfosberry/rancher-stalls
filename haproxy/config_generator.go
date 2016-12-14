package haproxy

import (
	"bytes"
	"text/template"

	"github.com/bfosberry/rancher-stalls/services"
)

const (
	configTemplate = `{{$port := .BackendPort}}
frontend api
    mode tcp
    {{range .Containers}}
    bind *:{{.ExternalPort}}
    acl dest{{.Index}} dst_port {{.ExternalPort}}
    use_backend Backend{{.Index}} if dest{{.Index}}
    {{end}}{{range .Containers}}
backend Backend{{.Index}}
    mode tcp
    server Backend{{.Index}} {{.IP}}:{{$port}} check
{{end}}`
)

type templateData struct {
	BackendPort int
	Containers  []services.Container
}

// GenerateConfig builds a haproxy config from a template for the specified
// backend port and services
func GenerateConfig(service *services.Service, backendPort int) (string, error) {
	t := template.New("config")
	t = template.Must(t.Parse(configTemplate))
	data := &templateData{
		BackendPort: backendPort,
		Containers:  service.Containers,
	}
	buf := bytes.NewBuffer(nil)
	err := t.Execute(buf, data)
	return buf.String(), err

}
