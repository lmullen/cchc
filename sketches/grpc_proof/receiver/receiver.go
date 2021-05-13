// The receiver listens for messages from the transmitter, and prints only words
// that contain the letter 'u'.
package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"

	pb "github.com/lmullen/cchc/sketches/grpc_proof/grpcproof"
)

var (
	port = 10000
)

// Receiver gets documents and does stuff with them.
type Receiver struct {
	pb.UnimplementedTransmitterServer
}

// ProcessDocument prints out the received document.
func (s *Receiver) ProcessDocument(ctx context.Context, doc *pb.Document) (*pb.DocumentAcknowledgement, error) {
	if hasU(doc.Word) {
		t := time.Since(doc.Sent.AsTime())
		fmt.Printf("Received %s after %s. It has a 'u'.\n", doc.Word, t)
	}
	return &pb.DocumentAcknowledgement{}, nil
}

func newReceiver() *Receiver {
	r := &Receiver{}
	return r
}

func hasU(word string) bool {
	return strings.Contains(word, "u")
}

func main() {
	log.Println("Starting the receiver")

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterTransmitterServer(grpcServer, newReceiver())
	grpcServer.Serve(lis)

}
