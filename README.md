# Fibonacci Sequence Challenge

* API will handle over ~6k requests per second on a 512MiB 1 Core Linux Box

* This API has five end points
  - /current
  - /next
  - /previous
  - /reset
  - /
  Each end point will return the curresponding fibonacci sequence for the current entry position

## Requirements
Any Linux box with Go installed

## Installation


To setup as user ubuntu

- Unzip or Gitclone this repo and run the binary to home directory of user [ubuntu]

RUN The following

  - $ cd /home/ubuntu/
  - $ sudo make all

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
  - Systemd daemon 

## Achieving 1k requests and above. 
This one is a bit tricky.
The best way to eliminate restrictions placed by operating systems and network interfaces is to run tests from localhost.
Or Increase the maximum number of TCP IP connections allowed in linux.


Included are various screenshots of Apache bench tests.
On 512MiB Ram 1 CPU on AWS without restrictions on Packet rates, this API's throughput was 4511/requests per second with zero errors on read/write endpoint.
On DigitalOcean with 8GB Ram and 4 CPUS it handled 8492 requests per secodn with zero error on read/write endpoint.
Results are significantly better but vary on non write end points.

See test_shots folder for examples screen shots of Apache bench tests.

![alt text](https://github.com/Prosp3r/pexservice/blob/master/test_shots/Screen%20Shot%202020-09-29%20at%205.34.58%20PM.png)
![alt text](https://github.com/Prosp3r/pexservice/blob/master/test_shots/Screen%20Shot%202020-09-29%20at%205.34.42%20PM.png)
![alt text](https://github.com/Prosp3r/pexservice/blob/master/test_shots/Screen%20Shot%202020-09-29%20at%205.35.32%20PM.png)

