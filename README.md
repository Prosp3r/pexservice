# PEX Service 

## Introduction
PexService is a set of API end points designed as an example of a **resillient, high throughput** API. 
This means it can handle a fairly substantial number of api requests in very short time and carry out all that's required of it while maintaining a fairly stable uptime on minimal compute resources (a single linux server).

Its goal is to showcase how a simple design implemented in Go programming language can handle complex tasks in production easily.

After going through and carrying out the instructions in this document, you would have achieved the following.
+ Gained some understanding of how key aspects of the pexservice code works.
+ Setup a simple linux server to run the pex service.
+ Installed Go.
+ Ran a simulated load test using [ApacheBench](https://httpd.apache.org/)

## How it works

Using the fibonacci sequence challenge as example workload, pex service will carry out thousands of calculations, persists or store the results as it goes and in the case of a crash, recover, and start from where it stoped.


This API has five end points. I've describe what each endpoint does in the table below.

<table> 
        <tr> <td> End Point </td><td> Function </td></tr>
        <tr> <td> / </td><td> Home endpoint - Displays a description of the other end-points </td></tr>
        <tr><td> /current </td><td> Current - Displays the current fibonacci calculated by the service </td> </tr>
        <tr> <td> /next </td><td> Next - Calculates and displays the next fibonacci number relative to the current one </td></tr>
        <tr> <td> /previous </td><td> Previous - Calculates and displays the previous fibonacci number relative to the current one </td></tr>
        <tr> <td> /reset </td><td> Resets the fibonacci number calculation to the begining </td></tr>        
</table>



#### Data Persistence
Part of being resilient is expecting the worst and planning for it when designing the software.
I decided to persist the calculations with as light a footprint as could be managed. 
To do this, I save the calculations per time to a `.csv` file on the system.
This way, in case the entire system crashes due to circumstances beyond our control all previous calculations are left intact and the system can pick up where it left off by loading the previous data from the `.csv` file.


#### Keep Alive
This feature simply means the system expects that at some point the pexservice could be terminated.
To make sure the system is always up and running, we've adopted two methods of deployment.

1. Deploying it as a **systemd daemon**
2. Deploying it as a **docker container** that can be deployed in a kubernetes cluster.

In this document we are focused on the systemd daemon.



## Prerequisites
+ An Ubuntu linux server or VPS with minimum 512MiB and 1 Core CPU.
+ Command line prompt access to the Ubuntu server
  _A VPS could be gotten from any of the following._
  + [Amazon AWS](https://aws.amazon.com)
  + [DigitalOcean](https://digitalocean.com)

While these instructions may work on other linux server types, I specify Ubuntu because the process has been thoroughly tested on Ubuntu(8+) linux servers.


## Installation

To setup as user ubuntu
Log on to your server through terminal with the user ubuntu.


Unzip or Gitclone this repo and run the binary to home directory of user [ubuntu]

       $ git clone https://github.com/Prosp3r/pexservice.git

RUN The following commands in your server terminal

      $ sudo apt update     
      $ sudo apt install make     
      $ cd /home/ubuntu/pexservice/      
      $ sudo make all



## Testing with Apache Bench
  + Install Apache bench

        $ sudo apt install apache2-utils
        $ ab -c 10 -n 10000 -r http://localhost:8081/

  The above command sends ten thousand connections at 10 concurrent connections to the / endpoint



## Extras


#### Using Docker containers

To use the docker image that could also run on Kubernetes or other orchestration systems, run the following command on your terminal.

        $ docker run -d -p 8081:8081  sirpros/pexservice:v1.1
        
  - Keepalive for this is possible through Kubernetes cluster.
  - Load performance on a single docker container vary depending on the docker host system capacity, requests throughput of between 560 to 9811 requests/second are achieved well for both non-write/read and write/read endpoints. 
  The rates varies wide depending on the amount of CPU power being on host. 
  Persistence is on container though could be modified to work with clusterwide storage.


#### Achieving 1k requests and above. 
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
