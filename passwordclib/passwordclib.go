package passwordclib

import (
	"encoding/json"
	"net"
)

func GetPassword(key string) (string, error) {
	conn := connect()
	return sendAndCheckResponse(Request{"get", key, ""}, conn)
}

func SetPassword(key string, value string) (string, error) {
	conn := connect()
	return sendAndCheckResponse(Request{"set", key, value}, conn)
}

func Version() string {
	return version
}

const (
	port    = "127.0.0.1:6969"
	version = "1.0.0"
)

type Error struct {
	Description string
}

func (e Error) Error() string {
	return e.Description
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


func connect() net.Conn {
	c, err := net.Dial("tcp", port)
	if err != nil {
		return nil
	}
	return c
}

func sendAndCheckResponse(request Request, conn net.Conn) (string, error) {
	if conn == nil {
		return "", Error{"could not connect"}
	}
	defer conn.Close()
	sendRequest(request, conn)
	return readResponse(conn)
}

func sendRequest(request Request, conn net.Conn) {
	data_out, _ := json.Marshal(request)
	conn.Write(data_out)
}

func readResponse(conn net.Conn) (string, error) {
	response := Response{"error", ""}
	buf := make([]byte, 4096)
	numbytes, err := conn.Read(buf)
	if numbytes == 0 {
		return "", Error{"no data received"}
	}
	if err != nil {
		return "", err
	}
	jsonErr := json.Unmarshal(buf[:numbytes], &response)
	if jsonErr != nil {
		return "", jsonErr
	}
	if response.Result != "ok" {
		return "", Error{"received result: "+response.Result}
	}
	return response.Value, nil
}
