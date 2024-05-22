package main

import (
	"testing"
)

func TestServer(t *testing.T) {
	/*
		peers, _ := lru.New(1024)
		serverAddr, err := punchServer(peers, 9999)
		if err != nil {
			t.Fatal(err)
		}
		log.Printf("server addr=%s", serverAddr)

		client, err := net.ListenPacket("udp", "127.0.0.1:")
		if err != nil {
			t.Fatalf("error starting client err=%s", err)
		}
		defer func() { _ = client.Close() }()

		msg := []byte("ping")
		_, err = client.WriteTo(msg, serverAddr)
		if err != nil {
			t.Fatalf("error sending data to server err=%s", err)
		}
	*/
}
