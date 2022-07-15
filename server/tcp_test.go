package server_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/SteveWXT/pubsub/server"
)

func TestTCPStart(t *testing.T) {
	fmt.Println("Starting TCP test...")

	go func() {
		if err := server.Start([]string{"tcp://127.0.0.1:1445"}); err != nil {
			t.Fatalf("Unexpected error - %s", err.Error())
		}
	}()
	<-time.After(time.Second)
}
