package main

import (
	"encoding/json"
	"fmt"
	"github.com/takama/daemon"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (

	// name of the service
	name        = "passwordd"
	description = "Simple Password storing daemon"

	// port which daemon should be listen
	//port = "127.0.0.1:6969"
	port = "0.0.0.0:6969"

	// the current version
	version = "1.0.1"
)

// dependencies that are NOT required by the service, but might be used
var dependencies = []string{}

var stdlog, errlog *log.Logger

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

type Request struct {
	Action string `json:"action"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}

type Response struct {
	Result string `json:"result"`
	Value  string `json:"value"`
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

	usage := "Usage: pwd-daemon install | remove | start | stop | status"

	// if received any kind of command, do it
	// else just run the daemon on the console
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		case "version":
			return version, nil
		default:
			return usage, nil
		}
	}
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Set up listener for defined host and port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return "Possibly was a problem with the port binding", err
	}

	// set up channel on which to send accepted connections
	listen := make(chan net.Conn, 100)
	go acceptConnection(listener, listen)

	passwords := make(map[string]string)
	// loop work cycle with accept connections or interrupt
	// by system signal
	for {
		select {
		case conn := <-listen:
			go handleClient(conn, passwords)

		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			stdlog.Println("Stoping listening on ", listener.Addr())
			listener.Close()
			if killSignal == os.Interrupt {
				return "Daemon was interrupted by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}
}

// Accept a client connection and collect it in a channel
func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		listen <- conn
	}
}

func handleClient(client net.Conn, passwords map[string]string) {
	for {
		request := &Request{}
		bad_response := &Response{"error", ""}
		buf := make([]byte, 4096)
		numbytes, err := client.Read(buf)
		error_data_out, _ := json.Marshal(bad_response)
		if numbytes == 0 || err != nil {
			client.Write(error_data_out)
			client.Close()
			return
		}

		jsonErr := json.Unmarshal(buf[:numbytes], &request)
		if jsonErr != nil {
			client.Write(error_data_out)
			client.Close()
			return
		}

		if request.Action == "get" {
			secret, exists := passwords[request.Key]
			good_response := &Response{"ok", secret}
			good_data_out, _ := json.Marshal(good_response)
			if exists == true {
				client.Write(good_data_out)
			} else {
				client.Write(error_data_out)
			}
		} else if request.Action == "set" {
			_, exists := passwords[request.Key]
			good_response := &Response{"ok", ""}
			good_data_out, _ := json.Marshal(good_response)
			if exists == true {
				client.Write(error_data_out)
			} else {
				passwords[request.Key] = request.Value
				client.Write(good_data_out)
			}
		} else if request.Action == "overwrite" {
			_, exists := passwords[request.Key]
			good_response := &Response{"ok", ""}
			good_data_out, _ := json.Marshal(good_response)
			if exists == true {
				delete(passwords, request.Key)
				passwords[request.Key] = request.Value
				client.Write(good_data_out)
			}
		} else {

			client.Write(error_data_out)
		}
		client.Close()
		return

	}
}

func init() {
	stdlog = log.New(os.Stdout, "", 0)
	errlog = log.New(os.Stderr, "", 0)
}

func main() {
	srv, err := daemon.New(name, description, daemon.SystemDaemon, dependencies...)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
