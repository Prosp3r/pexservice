# Fibonacci Sequence Challenge by Prosper Onogberie

* API will handle over ~3k requests per second on a 512MiB 1 Core Linux Box

* This API has five end points
  - /current
  - /next
  - /previous
  - /reset
  - /
  Each end point will return the curresponding fibonacci sequence for the current session.

## Installation

- Gitclone this repo and run the binary ./pexservice


OR

On Kubernetes or other Orchestration systems.
- docker run -d -p 8081:8081  sirpros/pexservice:v1.1
  - Keepalive for this is possible through Kubernetes cluster.
  - Load performance on a single docker container, requests of upto 560/second are handled well for non write/read endpoints. the rates varies wide on others depending on the amount of data being read and written. Persistence is on container but could be modified to work with clusterwide storage.

## Keep Alive
- KeepAlive feature is achieved with
  - docker container deployment on Kubernetes cluster
  - Systemd