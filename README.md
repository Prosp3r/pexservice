# Fibonacci Sequence Challenge by Prosper Onogberie

** This API will hancle ~1k requests per second on a 512MiB 1 Core Linux Box **

* This API has three end points on port 8081
  - /current
  - /next
  - /previous
  Each end point will return the curresponding fibonacci sequence for the current session.

## Installation
- Gitclone this repo and run the binary ./main 
or
- Docker run -d -p 8081:8081  sirpros/fibonacci

## Keep Alive
- KeepAlive feature is achieved with
  - docker container deployment on Kubernetes cluster
  - Systemd