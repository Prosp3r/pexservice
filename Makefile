#BUILDPATH=$(CURDIR)
GO=$(shell which go)
GOINSTALL=$(GO) install
GOGET=$(GO) get

#GOCLEAN=$(GO) clean


EXENAME=pexservice
FILE=/lib/systemd/system/pexservice.service
HOMESERVICEFILE=/home/ubuntu/pexservice/pexservice.service
#export GOPATH=$(CURDIR)

myname:
	@echo "Pex Service v1.1"

makedir:
	#@if [ ! -d $(BUILDPATH)/bin ] ; then mkdir -p $(BUILDPATH)/bin ; fi
	#@if [ ! -d $(BUILDPATH)/bin ] ; then mkdir -p $(BUILDPATH)/bin ; fi

get:
	#$(GOGET) github.com/gorilla/mux

#build:
#	$(GOINSTALL) $(EXENAME)
build:
	#sudo useradd pexservice -s /sbin/nologin -M
	sudo useradd {pexservice -s /sbin/nologin -M} || echo "User already exists."
	@echo "Creating service working directory"
	sudo /bin/mkdir -p pexservice
	sudo /bin/chown pexservice:adm pexservice
	sudo /bin/chmod 755 pexservice

	@echo "Building binary..."
	go build -o ./pexservice/pexservice main.go
docker:
	#docker build -f Dockerfile  --no-cache -t pexservice .
	#docker push sirpros/pexservice:latest
movefiles:
	@echo "entering tmp creating OS files..."
	cd /tmp
	@echo "creating OS user..."

	@echo "Creating log file"
	sudo /bin/mkdir -p /var/log/pexservice
	sudo /bin/chown syslog:adm /var/log/pexservice
	sudo /bin/chmod 755 /var/log/pexservice

	@echo "moving OS system files..."
	#ifeq ("$(wildcard /lib/systemd/system/pexservice))","")	
	#@echo "pexservice exists in systemd...skip copy and setting permissions ..."
	#else
	#	sudo mv /home/prosper/pexservice/pexservice.service /lib/systemd/system/.
	#@echo "changing file permissions on OS files..."
	#	sudo chmod 755 /lib/systemd/system/pexservice.service
	#endif
	if [ -f "$(FILE)" ]; then \
		echo "$(FILE) exists."; \
	else \
		echo "$(FILE) does not exist."; \
		cp $(HOMESERVICEFILE) /lib/systemd/system/. ;\
		chmod 755 /lib/systemd/system/pexservice.service; \
	fi

	@echo "enabling pexservice..."
	sudo systemctl enable pexservice.service
	@echo "starting pexservice..."
	sudo systemctl start pexservice
	@echo "monitoring pexservice ..."
	sudo systemctl status pexservice
	#sudo journalctl -f -u pexservice
	
all: get build movefiles