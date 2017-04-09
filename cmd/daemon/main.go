package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"google.golang.org/grpc"
	pb "github.com/inkah-trace/daemon/protobuf"
	context "golang.org/x/net/context"
	"time"
	"os"
	"net/http"
	"bytes"
	"io/ioutil"
)

type inkahDaemonServer struct {
	inkahServer *string
	eventChan chan *pb.ForwardedEvent
}

func (ids *inkahDaemonServer) RegisterEvent(ctx context.Context, event *pb.Event) (*pb.EventResponse, error) {
	hn, _ := os.Hostname()
	fe := pb.ForwardedEvent{
		TraceId: event.TraceId,
		SpanId: event.SpanId,
		ParentSpanId: event.ParentSpanId,
		EventType: event.EventType,
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Hostname: hn,
	}

	select {
	case ids.eventChan <- &fe: // Put 2 in the channel unless it is full
	default:
		fmt.Println("Event channel full. Discarding event.")
	}

	return &pb.EventResponse{}, nil
}

func processEvents(eventChan chan *pb.ForwardedEvent) {
	for e := range eventChan {
		sendEventToServer(e)
	}
}

func newInkahDaemonServer(inkahServer *string) pb.InkahServer {
	eventChan := make(chan *pb.ForwardedEvent, 1000)
	go processEvents(eventChan)
	return &inkahDaemonServer{
		inkahServer: inkahServer,
		eventChan: eventChan,
	}
}

func sendEventToServer(event *pb.ForwardedEvent) {
	url := "http://localhost:50052/event"

	b, err := json.Marshal(event)
	if err != nil {
		fmt.Println(err)
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	ioutil.ReadAll(resp.Body)
}

func main() {
	port := flag.Int("port", 50051, "Port for the server to run on")
	inkahServer := flag.String("server", "localhost:50052", "Host and port of Inkah Server.")
	flag.Parse()

	fmt.Printf("Starting Inkah daemon on port %d...\n", *port)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	grpcServer := grpc.NewServer()
	pb.RegisterInkahServer(grpcServer, newInkahDaemonServer(inkahServer))
	grpcServer.Serve(conn)
}
