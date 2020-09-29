# Fibonacci Sequence Challenge

* API will handle over ~6k requests per second on a 512MiB 1 Core Linux Box

* This API has five end points
  - /current
  - /next
  - /previous
  - /reset
  - /
  Each end point will return the curresponding fibonacci sequence for the current entry position

## Installation

- Gitclone this repo and run the binary ./pexservice


OR

On Kubernetes or other Orchestration systems.
- docker run -d -p 8081:8081  sirpros/pexservice:v1.1
  - Keepalive for this is possible through Kubernetes cluster.
  - Load performance on a single docker container vary depending on the docker host system capacity, requests throughput of between 560 to 9811 requests/second are achieved well for both non-write/read and write/read endpoints. 
  The rates varies wide depending on the amount of CPU power being on host. 
  Persistence is on container though could be modified to work with clusterwide storage.


## Keep Alive
- KeepAlive feature is achieved with
  - docker container deployment on Kubernetes cluster
  - Systemd