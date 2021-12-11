SHELL := /bin/bash

.PHONY: all pmond worker pmon2

mainfiles = pmon2.go pmond.go worker.go
config = config/config.yml
config_dir = /etc/pmon2/config
datadir = `cat $(config) | grep 'data:' | cut -d \" -f 2`
logsdir = `cat $(config) | grep 'logs:' | cut -d \" -f 2`
sockfile = `cat $(config) | grep 'sock:' | cut -d \" -f 2`
bindir = `cat $(config) | grep 'bin:' | cut -d \" -f 2`

service = `cat /proc/1/status | grep Name | cut -f 2`

build: clean build-pmond build-worker build-pmon2
build-pmond:$(mainfiles)
	@echo -n "building pmond..."
	@go build -ldflags "-X main.mode=pmond" -o pmond $(mainfiles)
	@echo "ok"
build-worker:
	@echo -n "building worker..."
	@go build -ldflags "-X main.mode=worker" -o worker $(mainfiles)
	@echo "ok"
build-pmon2:
	@echo -n "building pmon2..."
	@go build -o pmon2 $(mainfiles)
	@echo "ok"

install:
	@sudo mkdir -p $(config_dir)
	@sudo mkdir -p $(datadir)
	@sudo mkdir -p $(logsdir)
	@sudo mkdir -p $(bindir)
	
	@sudo cp $(config) $(config_dir)
	@if [ $(service) = systemd ]; then \
		sudo cp service/systemd/pmon2.service /etc/systemd/system/pmon2.service; \
	else \
		sudo cp service/init/pmon2.conf /etc/init/pmon2.conf; \
	fi
	@sudo cp logrotate/pmon2 /etc/logrotate.d/pmon2
	
	@sudo systemctl daemon-reload
	@sudo systemctl stop pmon2.service
	@sudo cp pmond $(bindir)
	@sudo cp worker $(bindir)
	@sudo cp pmon2 /usr/local/bin/.
	@sudo systemctl start pmon2.service

	@echo success
clean:
	@go clean