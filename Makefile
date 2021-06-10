VERSION := $(shell head debian/changelog  -n 1 | cut -d \( -f 2 | cut -d \) -f 1 | cut -d \- -f 1)
PACKAGE := $(shell head debian/changelog  -n 1 | cut -d \  -f 1 )
ARCH := amd64

INSTALL_PATH=$$(pwd)/debian/passwordd/
GOPATH:=/tmp/GOPATH
export GOPATH := /tmp/GOPATH
export GOBIN := $(GOPATH)/bin

all: compile client

compile:
	go get
	go build

client:
	make compile -C passwordc

install: compile client
	@install -D -m 755 passwordd $(INSTALL_PATH)/usr/bin/passwordd
	@install -D -m 644 passwordd.service $(INSTALL_PATH)/lib/systemd/system/passwordd.service
	#@install -D -m 644 passwordd.service $(INSTALL_PATH)/etc/systemd/system/passwordd.service
	@install -D -m 755 passwordc/passwordc $(INSTALL_PATH)/usr/bin/passwordc
	@install -D -m 755 passworddsync/passworddsync.py $(INSTALL_PATH)/usr/bin/passworddsync
	@install -D -m 644 passworddsync/passworddsync.service $(INSTALL_PATH)/lib/systemd/system/passworddsync.service
	@install -D -m 755 passworddsync/passworddsync-script $(INSTALL_PATH)/usr/bin/passworddsync-script
	
	@install -d -m 755 $(INSTALL_PATH)/usr/share/doc/passwordd/examples
	@install -D -m 644 passworddsync/passworddsync-*conf $(INSTALL_PATH)/usr/share/doc/passwordd/examples


clean:
	rm -rf $(INSTALL_PATH)
	rm -f passwordc/passwordc
	rm -f passwordd

deb:
	dpkg-buildpackage -uc -us

.PHONY: upload
upload: deb
	scp ../$(PACKAGE)_$(VERSION)_$(ARCH).deb root@repo.data.kit.edu:/var/www/debian/buster
	scp ../$(PACKAGE)_$(VERSION)_$(ARCH).deb root@repo.data.kit.edu:/var/www/debian/bullseye
	scp ../$(PACKAGE)_$(VERSION)_$(ARCH).deb root@repo.data.kit.edu:/var/www/ubuntu/bionic
	scp ../$(PACKAGE)_$(VERSION)_$(ARCH).deb root@repo.data.kit.edu:/var/www/ubuntu/focal
