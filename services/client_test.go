package services_test

import (
	"errors"
	"testing"

	"github.com/bfosberry/rancher-stalls/services"
	"github.com/rancher/go-rancher-metadata/metadata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	serviceName = "foo"
)

func TestGetService(t *testing.T) {
	fakeClient := &fakeMetadataClient{
		service: metadata.Service{
			Name: serviceName,
			Containers: []metadata.Container{
				metadata.Container{
					CreateIndex: 0,
					PrimaryIp:   "1.2.3.4",
				},
				metadata.Container{
					CreateIndex: 1,
					PrimaryIp:   "1.2.3.5",
				},
			},
		},
	}
	servicesClient := services.NewServicesClient(fakeClient)
	services, err := servicesClient.GetService(serviceName)

	assert.Nil(t, err)
	require.NotNil(t, services)

	assert.Equal(t, serviceName, services.Name)
	containers := services.Containers
	require.Equal(t, 2, len(containers))

	assert.Equal(t, 0, containers[0].Index)
	assert.Equal(t, "1.2.3.4", containers[0].IP)

	assert.Equal(t, 1, containers[1].Index)
	assert.Equal(t, "1.2.3.5", containers[1].IP)
}

func TestGetServiceError(t *testing.T) {
	testError := errors.New("test error")
	fakeClient := &fakeMetadataClient{
		service: metadata.Service{},
		err:     testError,
	}

	servicesClient := services.NewServicesClient(fakeClient)
	services, err := servicesClient.GetService(serviceName)

	assert.Nil(t, services)
	assert.Equal(t, err, testError)

}

type fakeMetadataClient struct {
	service metadata.Service
	err     error
}

func (f *fakeMetadataClient) GetSelfServiceByName(name string) (metadata.Service, error) {
	return f.service, f.err
}
