# PEX Service 

## Introduction
PexService is a set of API end points designed as an example of a **resillient, high throughput** API. 
This means it can handle a fairly substantial number of api requests in very short time and carry out all that's required of it while maintaining a fairly stable uptime on a single linux server.

Its goal is to showcase how a simple design implemented in Go programming language can handle complex tasks in production easily.

After going through and carrying out the instructions in this document, you would have achieved the following.
+ Gotten some understanding of how key aspects of the pexservice code works.
+ Setup a simple linux server to run the pex service.
+ Installed Go.
+ Ran a simulated load test using [ApacheBench](https://httpd.apache.org/)

## How it works

Using the fibonacci sequence challenge as example workload, pex service will carry out thousands of calculations, persists or store the results as it goes and in the case of a crash, recover, and start from where it stoped.


* This API has five end points
  - /current
  - /next
  - /previous
  - /reset
  - /
  Each end point will return the curresponding fibonacci sequence for the current entry position

### Data Persistence


### Keep Alive
- KeepAlive feature is achieved with
  - docker container deployment on Kubernetes cluster
  - Systemd daemon 



## Prerequisites
+ An Ubuntu linux server or VPS.
+ Command line prompt access to the Ubuntu server

A VPS could be gotten from any of the following.

+ [Amazon AWS](https://aws.amazon.com)
+ [DigitalOcean](https://digitalocean.com)

While these instructions may work on other linux server types, I specify Ubuntu because the process has been thoroughly tested on Ubuntu(8+) linux servers.


## Installation

To setup as user ubuntu

- Unzip or Gitclone this repo and run the binary to home directory of user [ubuntu]
- git clone https://github.com/Prosp3r/pexservice.git

RUN The following
  - $ sudo apt update
  - $ sudo apt install make
  - $ cd /home/ubuntu/pexservice/
  - $ sudo make all

OR

On use the docker image that could also run on Kubernetes or other Orchestration systems.
- docker run -d -p 8081:8081  sirpros/pexservice:v1.1
  - Keepalive for this is possible through Kubernetes cluster.
  - Load performance on a single docker container vary depending on the docker host system capacity, requests throughput of between 560 to 9811 requests/second are achieved well for both non-write/read and write/read endpoints. 
  The rates varies wide depending on the amount of CPU power being on host. 
  Persistence is on container though could be modified to work with clusterwide storage.


## Testing with Apache Bench
  - Install Apache bench
    $ sudo apt install apache2-utils
    $ ab -c 10 -n 10000 -r http://localhost:8081/
  The above command sends ten thousand hits at 10 concurrent connections to the / endpoint




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
