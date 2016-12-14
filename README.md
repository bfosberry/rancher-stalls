# rancher-stalls
Incremental port Haproxy implementation for Rancher services

This application is in Beta

Stalls provides an incremental port mapping proxy for scaled rancher services. This is for use with singleton containers which require direct proxied port mappings, but need to be scaled. The mapping should look like this:

* 5000 -> Container1:8080
* 5001 -> Container2:8080
* 5002 -> Container3:8080

To use stalls it MUST be executed within a Rancher service since it uses the metadata api for service discovery. It must also be run inside the same stack as the target service, and should be run on each host within the cluster with a restart policy enabled.

## Configuration

Stalls requires the following environment variables:

* `SERVICE_NAME`: The name of the Service to be proxied
* `BACKEND_PORT`: The static private port each container of the service listens on  
* `BASE_PORT`: The port on which to start incrementally mapping

The service should be run on each host of the cluster, or at least a few for availability. 

The target service should expose the `BACKEND_PORT` port in it's Dockerfile, and should be configured using a managed network layer. This will allow for easy routing beteen containers on one host and the stalls instance on another. There is no need to directly map ports for the containers in your service, and doing so may actually reduce tennancy options due to port reservation conflicts. 

## Features

Stalls will use the Rancher metadata service to discover the backend service and reload haproxy any time this changes. This means that once this service is run and configured as part of your stack, it should automatically scale with your hosts, and with your containers. 

Be sure to pre-map an appropriate number of ports for your service to support scale without downtime. A service upgrade extending port mappings may well cause connections to drop and error for a short time. Since Rancher does not support mapping port ranges to services, this process can be quite tedious.