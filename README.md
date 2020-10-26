# PEX Service 

## Introduction
PexService is a set of API end points designed as an example of a **resilient, high throughput** API. 
This means it can handle a fairly substantial number of API requests in very short time and carry out all that's required of it while maintaining a fairly stable uptime on minimal compute resources (a single Linux server).

Its goal is to showcase how a simple design implemented in Go programming language can handle complex tasks in production easily.

After going through and carrying out the instructions in this document, you would have achieved the following.
+ Gained some understanding of how key aspects of the pexservice code works.
+ Setup a simple Linux server to run the pex service.
+ Installed Go.
+ Ran a simulated load test using [ApacheBench](https://httpd.apache.org/)

## How it works

In mathematics, the [Fibonacci](https://en.wikipedia.org/wiki/Fibonacci_number) numbers, commonly denoted Fₙ, form a sequence, called the Fibonacci sequence, such that each number is the sum of the two preceding ones, starting from 0 and 1. 

Using the Fibonacci sequence challenge as example workload, pex service will carry out thousands of calculations, persists or store the results as it goes and in the case of a crash, recover, and start from where it stopped.


This API has five end points. I've described what each endpoint does in the table below.

<table> 
        <tr> <td> End Point </td><td> Function </td></tr>
        <tr> <td> / </td><td> Home endpoint - Displays a description of the other end-points </td></tr>
        <tr><td> /current </td><td> Current - Displays the current Fibonacci calculated by the service </td> </tr>
        <tr> <td> /next </td><td> Next - Calculates and displays the next Fibonacci number relative to the current one </td></tr>
        <tr> <td> /previous </td><td> Previous - Calculates and displays the previous Fibonacci number relative to the current one </td></tr>
        <tr> <td> /reset </td><td> Resets the Fibonacci number calculation to the beginning </td></tr>        
</table>


#### System flow

![alt text](https://github.com/Prosp3r/pexservice/blob/master/test_shots/pex_inaction.png)


In the illustration above, points 1, 2 and 3 are executed at the start of the system, point 4 is triggered by requests made on the API endpoints with read/write capabilities e.g. `/next` , `/previous` , `current` and `/reset`.

1. When the service is started, it reads the previous calculations from the `.csv` file store and sets it in memory.
2. It then starts a Go routine that keeps updating the `.csv` store with changes to the in-memory store independent of the main program.
3. The service also starts a http server using the gorilla mux package that listens for request on all endpoints.
4. When requests come in through any of the endpoints, a calculation is made where necessary and the results stored in the in-memory store.

To prevent read conflicts in the in-memory store, all reads and writes are protected with the sync package mutex lock function.







#### Data Persistence
Part of being resilient is expecting the worst and planning for it when designing the software.
I decided to persist the calculations with as light a footprint as could be managed. 
To do this, I save the calculations per time to a `.csv` file on the system.
This way, in case the entire system crashes due to circumstances beyond our control all previous calculations are left intact and the system can pick up where it left off by loading the previous data from the `.csv` file.


#### Keep Alive
This feature simply means the system expects that at some point the pexservice could be terminated.
To make sure the system is always up and running, we've adopted two methods of deployment.

1. Deploying it as a **systemd daemon**
2. Deploying it as a **Docker container** that can be deployed in a Kubernetes cluster.

In this document we are focused on the systemd daemon.


## Setting up pexservice

#### Prerequisites
+ An Ubuntu Linux server or VPS with minimum 512MiB and 1 Core CPU.
+ Command line prompt access to the Ubuntu server
  ###### _A VPS could be acquired from any of the popular vendors like [Amazon AWS](https://aws.amazon.com), [DigitalOcean](https://digitalocean.com)_


_While these instructions may work on other Linux server types, I specify Ubuntu because the process has been thoroughly tested on Ubuntu(20.04) Linux servers._


#### Installation

To setup as user Ubuntu
Log on to your server through terminal with the user Ubuntu.


       $ cd ~
       $ sudo apt update
       $ sudo apt install git

You can confirm that you have installed Git correctly by running the following command:

        $ git --version
 

. 

        Output 
        git version 2.18.0

        
Clone this repository to your local directory using the following command:

       $ git clone https://github.com/Prosp3r/pexservice.git

Our deployment uses the make utility to carry out some of the important installations needed to get the system up and running.
The following commands will install [GNU Make](https://www.gnu.org/software/make/) and setup the [Go](https://golang.org/) environment needed for pexservice to run and starts the pexservice as a [**system service**](http://manpages.ubuntu.com/manpages/bionic/man5/systemd.service.5.html) all at the same time.

      
      $ sudo apt install make     
      $ cd /home/ubuntu/pexservice/      
      $ sudo make all



## Load testing with Apache Bench
To test the ability of our service to handle large amount of connections in fairly short time, we will be using [Apache Bench](https://httpd.apache.org/docs/2.4/programs/ab.html) utility which comes installed by default in some operating systems. 
We will need to install it on our server however.

Run the following command to install Apache Bench: 

        $ sudo apt install apache2-utils
        $ ab -c 10 -n 10000 -r http://localhost:8081/

The second command above sends ten thousand requests at 10 concurrent requests to the / endpoint.
You can vary the requests by changing the parameters.

#### Our results

When we tested the service on a variety of systems, here are some results we got.

<table> 
        <tr><td> System configuration </td><td> Our Results </td><td> Screen Shot </td>  </tr>
        <tr><td> AWS - t2.nano Ubuntu 20.04[LTS](HVM) (0.5GiB, 1 vCPUs) </td><td> 4511.26 [#/sec], Avg time per request 0.222[ms] </td><td> https://bit.ly/35xOuDT </td>  </tr>
        <tr><td> DigitalOcean - Ubuntu [LTS] x64 20.04 (8GB, 4 vCPUs) </td><td> 9811.24 [Requests/sec], Avg time per request 0.102[ms] </td><td> https://bit.ly/2Tt3UUl </td>  </tr>
        
</table>


## Extras


#### Using Docker containers

To use the Docker image that could also run on Kubernetes or other orchestration systems, run the following command on your terminal.

        $ docker run -d -p 8081:8081  sirpros/pexservice:v1.1
        
  - Keep alive for this is possible through Kubernetes cluster.
  - Load performance on a single Docker container vary depending on the Docker host system capacity, requests throughput of between 560 to 9811 requests/second are achieved well for both non-write/read and write/read endpoints. 
  The rates varies wide depending on the amount of CPU power being on host. 
  Persistence is on container though could be modified to work with cluster wide storage.


## Contributing
We always welcome contributions in various forms.
Your contributions could be in form of pointing out errors or actually contributing some code updates.
Whatever category your desired contribution fall into, here's a simple guide for making that impact quickly.

#### Pointing out something

1. For people who don't have the time or for some reason are unable to make a code contribution themselves, please use the "Issues" link at the top of this document.

 + ![Click to view issues](https://github.com/Prosp3r/pexservice/blob/master/test_shots/pex_art_issues.fw.png)

 + ![Click to open new a issue](https://github.com/Prosp3r/pexservice/blob/master/test_shots/pex_art_newissues.fw.png)

Fill the new issue form and we'll take it from there.


## Conclusion

[Go(Golang)](https://golang.org) is a very powerful, feature packed programming language for developing systems that can handle fairly large amount of work loads.
It concurrency features particularly, can enable developers design very efficient systems with near unlimited capability with little system resource footprints.

If you've are contemplating learning a new language for the long term, give [Go](https://golang.org) a try.
