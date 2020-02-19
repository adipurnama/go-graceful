# Golang graceful shutdown sample

Sample golang http webservice with graceful shutdown mechanism.
GCP & AWS provides ephemeral VM deployment (preemptive / spot instance) with cheaper cost (up to 80%).
The caveat are its host can only live up to 24 hours & can be halted within ~30s.

## Solution

1. Backend service API provide healthcheck endpoint, in this case "/actuator/health".
2. Load Balancer monitor the healthcheck for every 2-5 seconds, any unhealthy nodes will be removed from serving backend group.
3. When kill/interrupt OS signal is received, service will mark itself as unhealthy (in this case "OUT_OF_SERVICE"), and then wait for n seconds so LB will remove it from backend group.
4. Wait for t seconds for ongoing request to be completed.
5. Shutdown the service.

For any given "h" halt signal timeout, n + t < h.
