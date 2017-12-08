
all: compile client

compile:
	go get
	go build

client:
	make compile -C passwordc
