package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"testing"

	lru "github.com/hashicorp/golang-lru"
	"github.com/yinheli/udppunch"
)

func TestServer(t *testing.T) {
	port := 9999
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		l.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	peers, _ := lru.New(1024)
	err = punchServer(ctx, peers, serverAddr)
	if err != nil {
		t.Fatal(err)
	}
	defer cancel()

	client, err := net.ListenPacket("udp", "127.0.0.1:")
	if err != nil {
		t.Fatalf("error starting client err=%s", err)
	}
	defer func() { _ = client.Close() }()

	_, err = client.WriteTo([]byte("foo"), serverAddr)
	if err != nil {
		t.Fatalf("error sending data to server err=%s", err)
	}
	got := len(peers.Keys())
	expected := 0
	if got != expected {
		t.Fatal("I sent a bad payload, we should not have peers")
	}

	pubKeyB64 := "szyWfOODAAu6Ma1pufrAng+atPHnSn/dfGm61JcvDQE="
	pubKey, _ := base64.StdEncoding.DecodeString(pubKeyB64)

	data := make([]byte, 0, 32+1)
	data = append(data, udppunch.HandshakeType)
	data = append(data, pubKey[:]...)
	_, err = client.WriteTo(data, serverAddr)
	if err != nil {
		t.Fatalf("error sending data to server err=%s", err)
	}

	data = make([]byte, 0, 32+1)
	data = append(data, udppunch.ResolveType)
	data = append(data, pubKey[:]...)
	_, err = client.WriteTo(data, serverAddr)
	if err != nil {
		t.Fatalf("error sending data to server err=%s", err)
	}

	buf := make([]byte, 1024)
	n, _, err := client.ReadFrom(buf)
	if err != nil {
		t.Fatal(err)
	}
	readPeers := udppunch.ParsePeers(buf[:n])
	if len(readPeers) != 1 {
		t.Fatal("should have been able to read one peer")
	}

	readKey, _ := readPeers[0].Parse()
	expectedKey := udppunch.NewKeyFromStr(pubKeyB64)
	if readKey != expectedKey {
		t.Fatalf("server returned invalid key got=%s expected=%s", readKey.String(), pubKey)
	}

}
