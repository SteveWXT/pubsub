package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/SteveWXT/pubsub/server"
)

// TestAuthStart tests an auth start process
func TestStart(t *testing.T) {
	fmt.Println("Starting SERVER test...")

	// test for error if an auth is provided w/o a token
	go func() {
		if err := server.Start([]string{"tcp://127.0.0.1:1446"}); err == nil {
			t.Fatalf("Expecting error - %s", err.Error())
		}
	}()
	<-time.After(time.Second)
}
