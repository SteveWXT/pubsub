package server_test

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/SteveWXT/pubsub/server"
)

// TestAuthStart tests an auth start process
func TestStart(t *testing.T) {
	fmt.Println("Starting SERVER test...")

	go func() {
		if err := server.Start([]string{"tcp://127.0.0.1:1446"}); err == nil {
			t.Fatalf("Expecting error - %s", err.Error())
		}
	}()
	<-time.After(time.Second)
}

func TestStartWithLS(t *testing.T) {
	fmt.Println("Starting SERVER test...")
	ls, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal("Error: TCP listener cannot start")
	}

	go func() {
		if err := server.StartWithLS(ls); err == nil {
			t.Fatalf("Expecting error - %s", err.Error())
		}
	}()
	<-time.After(time.Second)
}
