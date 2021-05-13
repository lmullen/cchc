// The transmitter reads a dictionary of words, and sends the ones that contain
// the letter 'e' to the receiver.
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/lmullen/cchc/sketches/grpc_proof/grpcproof"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	serverAddr = "localhost:10000"
	dict       = "/usr/share/dict/words" // Unix dictionary on MacOS
)

func sendDocument(client pb.TransmitterClient, doc *pb.Document) {
	fmt.Printf("Sending %s at %v. It has an 'e'.\n", doc.Word, doc.Sent.AsTime())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := client.ProcessDocument(ctx, doc)
	if err != nil {
		log.Println(err)
	}
}

func createDocument(word string) *pb.Document {
	return &pb.Document{
		Word: word,
		Sent: timestamppb.New(time.Now()),
	}
}

func hasE(word string) bool {
	return strings.Contains(word, "e")
}

func main() {
	log.Println("Dialing the server")
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewTransmitterClient(conn)

	f, err := os.Open(dict)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)

	log.Println("Start sending documents")

	for scanner.Scan() {
		word := scanner.Text()
		// Only transmit words that have an "e"
		if hasE(word) {
			sendDocument(client, createDocument(word))
		}
		time.Sleep(100 * time.Millisecond)
	}

	log.Println("Done sending documents")
}
