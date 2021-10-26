# Password Hash Service

## Test, Build, and Run
Make file
```sh
make test
make
make run
```
Use cURL command
```sh
curl -X POST --data "password=angryMonkey" http://localhost:8080/api/v1/hash
curl http://localhost:8080/api/v1/hash/1
curl http://localhost:8080/api/v1/stats
curl -X POST http://localhost:8080/api/v1/shutdown
```

## Password Hash Service API
The Password Hash Service API is documented on this [OpenAPI](api/swagger.yaml) document.
The OpenAPI document can be used to drive the API on a local test environment on port 8080.
The local test service can be run by running the following commands on the terminal:
```sh
go build -o pwdhash main.go && ./pwdhash
```
### Design
The hash might not be immediately available after the initial API call. This means the API
SHOULD provide provide a mechanism to check if the hash is ready as well as an estimate of when
the hash will be available. This design can help the clients determine when the next call
can be made in order to avoid unnecessary calls that can overload the service. Finally, we
need to pay careful attention to the API such that it can be hosted anywhere such as AWS
without having to worry about extended infrastructure such as proxy or load balancer.

#### Hash: Simple Approach
1. Creating a hash should return 200 status code. This is the go-to status code.
2. Creating a hash should return 202 status code. This is the status that confirms the call will be processed. 
3. Checking to see if a hash ID is ready returns status 200 with or no content.
4. Checking to see if a hash ID is ready returns status 204 or no content or 200 with content.

There is nothing wrong with using these statuses to represent a successful call to the service.
For example, most APIs return 200 so signal success while others return 202 to signal that 
the request was accepted. However, both have very small drawbacks that are worth mentioning.
Status 200 is a cacheable status which means the response COULD be cached by the browser or
web caches. The same issue is true for status 204 when checking for a hash status that is
not ready.

Status 202 is not cacheable, but does not provide enough semantics for our API. For example,
status 202 is non-comittal. This means that the status does not convey the idea that some
resource was created after an API call which is the case for the creation API.

Because the creation API returns the hash ID / job ID, we nee to find a better status that
represent the creation of a new resource. This means that status 200 and 202 are not ideal
to represent the creation of a hash that might take long time. Finally, using 204 and 200
might reasonable statuses, but they cannot convey the notion of of a long running task.

#### Hash: Better Approach
1. Creating a hash should return 201.
2. Checking to see if hash ID is ready returns 503 and 200 when the hash is not ready and ready respectively.

There are advantages when using statuses 201, 503, and 200 within this API context. The
firs benefit is that status 201 is not cacheable and represents the idea that a resource
was created. Furthermore, the API can return the location the newly created resource in
the "Location" header. This fits perfectly with the API since the newly created resource
will be used to check if hash is ready or not.

Another benefit of using status 503 and 200 for checking the status is that 503 can carry
the notion of time. For example, status 503 expresses the meaning that a service is
temporarily unavailable. Also, status 503 is further enhanced by allowing servers to 
respond with the "Retry-After" header. This header can provide a hint as to when the client
can retry the response in order to not overload the service. Finally, status 200 for
the response of hash is ideal because it SHOULD lower the load on the service after it
is cached.

#### Service Configuration
The service can be configured through the configuration file config.json. The configuration
file is simple and self explanatory. For information look at the documentation in the config
package.

#### Shutdown
The shutdown can be indefinite. This means should not return 503 since we do know for how
long the service will be down. In this case, we will return 500 to signal an error. Also,
we will log the error if a request is sent when the server is shutting down.

#### Tasking
It was assumed that the tasks needed to wait for some time. That time wait can be configured
through the configuration file. The tasks are managed through the task manager in the
task package. 

#### Logging
It is important to have different levels of logging in order to be able to diagnose
production issues. In order to be able to correlate issues, we need to be able to track
an action across logs. The "X-Request-ID" header will be used to correlate log entries.
This request for now will use the syntax "REQID(integer)" to correlate the logs. The
service will attempt to get the request header "X-Request-ID" . If such header does not
exist, an increasing number will be used for the request ID.

Logs can be printed to the standard out, error, or file. The different type of logs can
be configured through the configuration file.

#### Reporting Issues
The idea of using a request ID is useful when reporting an issue. For example, the 
request ID can be used in an API response so that logs can be correlated later when
diagnosing the issue. Perhaps this can done on the next iteration of the product.

#### Security
##### HTTPS Support
Currently, the service only works with HTTP. The service can be enhanced by adding
some configuration items to use HTTPS. The configuration can have the locations for
the certificates to be used by the service. This can be done on the next iteration.

##### Denial of Service Attack
*Password Length*
Password should not be that big and should be around 50 characters long at most.
Currently, there is no limit on the password length. However, this can be easily
mitigated by adding a configuration item and changing the password strength to
fail when the maximum is reached. This should help mitigate some of the memory
consumption on an attack.

*Concurrent Request*
The system starts thrashing after 5 million concurrent request during a simple test.
This means that the service should know many pending requests there are and avoid
allowing more requests to keep the system running. This can be easily achieved by
adding a pending task counter in the task manager. Therefore, a solution for this
problem would be to add a configuration item to tell the task manger to stop
more new hash requests when the pending task exceeds a certain threshold.

This can happen with legitimate users and the API should be able to tell the clients
to wait for some time. This can be easily done by adding a configuration item to tell
clients how long to wait with 503 status and a "Retry-After" header. Legitimate user
will use the "Retry-After" header and wait to allow the service to come back to a
healthy state.

Bad actors will not not care about the "Retry-After" header and will continue to
make requests to the service. In such case, the IP of the client can be tied
503 request and "Retry-After". If the client does not obey the "Retry-After" 
header it can be considered a hostile user and blocked from using the service for
some time to help the service come back to a healthy state.

### Code Coverage and Quality.
The code coverage should be above 90% on all the packages. This will ensure a
minium quality for the product.


