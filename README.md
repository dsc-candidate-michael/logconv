# Summary

Logconv is a package that can run batch jobs to parse your logs and convert them to statsd-compatible messages.

## Overall design

### Logconv Go package

The logconv has 4 main components

1. *LogParser*: This interface is small and allows users to provide different ways to parse an individual log. It should return a *ReqDetail*, which is a simple struct that has the data we want to report in the statsd-compatible messages.

2. *LogObserver*: This will observe the specified log file. Whenever a new log is written to the file, it will use the *LogParser* to attempt to parse the log. If successful, it sends a *ReqDetail* over a channel.   

3. *ReqDetailConsumer*: This will consume the ReqDetail messages sent over the aforementioned channel. It also provides a mechanism for flushing the data.

4. *LogConv*: This setups up the *LogObserver* to *ReqDetailConsumer* pipeline. Then, it will use a ticker in order to basically flush the data to stdout every n seconds. We use a pipe and `tee` in order to also write this to file when we deploy to Kubernetes. 

### Kubernetes

I decided early on that I was going to have logconv container essentially be a "sidecar" container in the nginx pods. This seemed to make accessing the nginx logfiles simple, and ensured that we would always have a 1 to 1 mapping of nginx to logconv containers. 

There are a few Kubernetes constructs being used in the deployment example. 

1. A Deployment is used as any easy way to create a ReplicaSet of nginx-logcongv pods. This ultimately can let you spin up extra replicas pretty easily. 

2. A LoadBalancer service is used for the nginx-logconv Deployment. 

3. An Ingress is used to make the nginx-logconv LoadBalancer service accessible externally. 

4. A ConfigMap is used to configure the nginx-logconv pods. 

## Build, test, and deploy steps

### Buid

Running `make` will attempt to build the project locally.

Running `make image` will build a Docker image which may make things a bit easier. 

### Test

To run tests against your local build, you can run `make test`. 

To run the tests in a Docker container, you can run `make test-container`. 

### Deploy

1. `make create-cluster` will start minikube, enable the necessary addons, and use the docker daemon from your host so we can leverage the docker images we build locally.  
2. Run `eval $(minikube docker-env)` in order to use the docker daemon from your host so we can leverage the dodcker images we build locally. 
3. `make configure` will create a ConfigMap to be used by the nginx-logconv pods. 
4. `make image` will create the necessary images for our nginx-logconv pods. 
5. `make deploy` will create the Ingress, Deployment, and LoadBalancer Service. 

You can now access the Nginx service by grabbing the IP address 
`minikube ip`. Throw this IP into curl, wget, or your favorite browser. 

You can monitor the statsd-compatible messages by running

`kubectl logs <pod_name> logconv`

Or, if you happen to have [kubetail](https://github.com/johanhaleby/kubetail) installed, you can just run `kubetail nginx` to see the ouput from all of the pods. 

## Caveats

I had to change `x_forwarded_proto_or_scheme` to `$http_x_forwarded_proto_or_scheme` in order for the nginx directive to work as expected in nginx 1.15. Apart from that, it should be the same as the one in the spec.

## Nice to haves, but ran out of time:

### Cleanup tests
* Currently, the tests have some duplicate setup / teardown logic that could certainly be refactored a bit.
* Some of the tests could be considered 'whitebox' tests and have some internal knowledge of the package. 
* Would be nice to have a bit better test coverage as well.

### Use better logging system throughout the codebase
* Since I am relying on stdout piping and `tee` in order to log the required output to file and to stdout, I really should have a better logging system through the codebase so that other developers wouldn't accidentally log to stdout and accidentally corrupt the output data. 

### Code cleanup
* The configuration of the various different interfaces and structs in this pkg is kind of convuluted. Seems like this could be simplified. 
* The pkg is currently exposing all of the things. We should only be exposing methods and variables that users of the pkg would need to make use of. 

### Make it a bit more generalizable
* While I tried to make most of the dependencies on interfaces rather than structs, there are some parts of the codebase that are not really extendable as is. In particular, `reqdetailconsumer` is very much tied to the `ReqDetail` data structure and to statsd as well. 

### Health checks
* Would be easy to add basic http liveness and readiness probes to Nginx. It would be slightly harder to add this for the logconv container in its current form.

