# gRPC Proof of Concept

This is a simple proof of concept using gRPC. 

The transmitter (client) reads a dictionary of words and transmits the ones with the letter 'e'. This is like the crawler, in that it traverses a set of resources and transmits only certain ones. 

The receiver (server) prints out the words that contain the letter 'u' and the time since the message was transmitted. This is like doing work with the machine-learning model.

To run this, first start the receiver in one terminal.

```
go run receiver/receiver.go
```

Then start the transmitter in another terminal.

```
go run transmitter/transmitter.go
```
