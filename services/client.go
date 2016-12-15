package services

import (
	"strconv"

	"github.com/rancher/go-rancher-metadata/metadata"
)

// MetadataServicesClient is a subset of metadata.Client with only
// the fields the services client needs
type MetadataServicesClient interface {
	GetSelfServiceByName(string) (metadata.Service, error)
}

// Service represents a rancher service
type Service struct {
	Name       string
	Containers []Container
}

// Container represents a container on rancher
type Container struct {
	IP           string
	Index        int
	ExternalPort int
}

// Client represents an interfaces for clients to rancher services
type Client interface {
	// GetServices returns a Service by name
	GetService(string, int) (*Service, error)
}

// NewServicesClient returns a new services client implementation
func NewServicesClient(metadataClient MetadataServicesClient) Client {
	return &client{
		metadataClient: metadataClient,
	}
}

type client struct {
	metadataClient MetadataServicesClient
}

func (c *client) GetService(name string, basePort int) (*Service, error) {
	metadataService, err := c.metadataClient.GetSelfServiceByName(name)
	if err != nil {
		return nil, err
	}

	service := &Service{
		Name: name,
	}

	containers := make([]Container, 0, len(service.Containers))
	for _, c := range metadataService.Containers {
		serviceIndex, err := strconv.Atoi(c.ServiceIndex)
		if err != nil {
			serviceIndex = c.CreateIndex
		}
		container := Container{
			IP:           c.PrimaryIp,
			Index:        serviceIndex,
			ExternalPort: basePort + serviceIndex,
		}
		containers = append(containers, container)
	}
	service.Containers = containers
	return service, nil
}
