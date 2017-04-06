package main

import (
	"flag"
	"fmt"
	"net"
	"google.golang.org/grpc"
	pb "github.com/inkah-trace/daemon/protobuf"
	context "golang.org/x/net/context"
)

type inkahDaemonServer struct{
	eventChan chan *pb.Event
}

func (ids *inkahDaemonServer) RegisterEvent(ctx context.Context, event *pb.Event) (*pb.EventResponse, error) {
	ids.eventChan <- event
	return &pb.EventResponse{}, nil
}

func processEvents(eventChan chan *pb.Event) {
	for e := range eventChan {
		fmt.Println(e)
	}
}

func newInkahDaemonServer() pb.InkahServer {
	eventChan := make(chan *pb.Event, 1000)
	go processEvents(eventChan)
	return &inkahDaemonServer{
		eventChan: eventChan,
	}
}

func main() {
	port := flag.Int("port", 50051, "Port for the server to run on")
	flag.Parse()

	fmt.Printf("Starting Inkah daemon on port %d...\n", *port)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	grpcServer := grpc.NewServer()
	pb.RegisterInkahServer(grpcServer, newInkahDaemonServer())
	grpcServer.Serve(conn)
}
