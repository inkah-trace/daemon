package main

import (
	"flag"
	"fmt"
	"net"
	"google.golang.org/grpc"
	pb "github.com/inkah-trace/daemon/protobuf"
	context "golang.org/x/net/context"
	"time"
	"os"
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
		RequestId: event.RequestId,
		EventType: event.EventType,
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Hostname: hn,
	}

	ids.eventChan <- &fe
	return &pb.EventResponse{}, nil
}

func processEvents(eventChan chan *pb.ForwardedEvent) {
	for e := range eventChan {
		fmt.Println(e)
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
