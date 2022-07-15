package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/SteveWXT/pubsub/server"
)

// TestWSStart tests to ensure a server will start
func TestWSStart(t *testing.T) {
	fmt.Println("Starting WS test...")

	go func() {
		if err := server.Start([]string{"ws://127.0.0.1:8888"}); err != nil {
			t.Fatalf("Unexpected error - %s", err.Error())
		}
	}()
	<-time.After(time.Second)
}
