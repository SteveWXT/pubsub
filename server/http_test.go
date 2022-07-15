package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/SteveWXT/pubsub/server"
)

// TestHTTPStart tests to ensure a server will start
func TestHTTPStart(t *testing.T) {
	fmt.Println("Starting HTTP test...")

	go func() {
		if err := server.Start([]string{"http://127.0.0.1:8080"}); err != nil {
			t.Fatalf("Unexpected error - %s", err.Error())
		}
	}()
	<-time.After(time.Second)
}
