package main

import (
	"fmt"
	"net"
	"os"

	"encoding/json"

	"github.com/inkah-trace/common"
	"github.com/inkah-trace/daemon"
)

/* A Simple function to verify error */
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

func main() {
	ServerAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9800")
	CheckError(err)

	// Now listen at selected port
	ServerConn, err := net.ListenUDP("udp", ServerAddr)
	CheckError(err)
	defer ServerConn.Close()

	buf := make([]byte, 1024)

	fmt.Println("Daemon running on 127.0.0.1:9800")

	for {
		n, addr, err := ServerConn.ReadFromUDP(buf)

		message := buf[0:n]
		fmt.Println("Received ", message, " from ", addr)

		e := inkah.Event{}
		err = daemon.Unmarshall(message, &e)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		go SendMessageUpstream(&e)
	}
}

func SendMessageUpstream(event *inkah.Event) {
	conn, _ := net.Dial("tcp", "127.0.0.1:9810")
	defer conn.Close()

	b, err := json.Marshal(event)
	if err != nil {
		fmt.Println("JSON Marshalling error: ", err)
	}

	fmt.Fprintf(conn, string(b))

}
