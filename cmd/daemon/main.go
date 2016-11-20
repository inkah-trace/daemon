package main

import (
	"fmt"
	"net"
	"os"
	"time"
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
		message := string(buf[0:n])
		fmt.Println("Received ", message, " from ", addr)

		timestamp := time.Now().Unix()

		go SendMessageUpstream(timestamp, message)

		if err != nil {
			fmt.Println("Error: ", err)
		}
	}
}

func SendMessageUpstream(timestamp int64, message string) {
	conn, _ := net.Dial("tcp", "127.0.0.1:9810")
	defer conn.Close()
	hn, _ := os.Hostname()
	fmt.Fprintf(conn, "INKAH\t%d:%s:%s\n", timestamp, hn, message)
}
