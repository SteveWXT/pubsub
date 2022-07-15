package server

import (
	"fmt"
	"strings"

	"github.com/SteveWXT/pubsub/core"
)

// GenerateHandlers ...
func GenerateHandlers() map[string]core.HandleFunc {
	return map[string]core.HandleFunc{
		"ping":        handlePing,
		"subscribe":   handleSubscribe,
		"unsubscribe": handleUnsubscribe,
		"publish":     handlePublish,
		// "publishAfter":     handlePublishAfter,
		"list":    handleList,
		"listall": handleListAll, // listall related
		"who":     handleWho,     // who related
	}
}

// handlePing
func handlePing(proxy *core.Proxy, msg core.Message) error {
	// goroutining any of these would allow a client to spam and overwhelm the server. clients don't need the ability to ping indefinitely
	proxy.Pipe <- core.Message{Command: "ping", Tags: []string{}, Data: "pong"}
	return nil
}

// handleSubscribe
func handleSubscribe(proxy *core.Proxy, msg core.Message) error {
	proxy.Subscribe(msg.Tags)
	return nil
}

// handleUnsubscribe
func handleUnsubscribe(proxy *core.Proxy, msg core.Message) error {
	proxy.Unsubscribe(msg.Tags)
	return nil
}

// handlePublish
func handlePublish(proxy *core.Proxy, msg core.Message) error {
	proxy.Publish(msg.Tags, msg.Data)
	return nil
}

// handlePublishAfter - how do we get the [delay] here?
// func handlePublishAfter(proxy *core.Proxy, msg core.Message) error {
// 	proxy.PublishAfter(msg.Tags, msg.Data, ???)
// 	go func() {
// 		proxy.Pipe <- core.Message{Command: "publish after", Tags: msg.Tags, Data: "success"}
// 	}()
// 	return nil
// }

// handleList
func handleList(proxy *core.Proxy, msg core.Message) error {
	var subscriptions string
	for _, v := range proxy.List() {
		subscriptions += strings.Join(v, ",")
	}
	proxy.Pipe <- core.Message{Command: "list", Tags: msg.Tags, Data: subscriptions}
	return nil
}

// handleListAll - listall related
func handleListAll(proxy *core.Proxy, msg core.Message) error {
	subscriptions := core.Subscribers()
	proxy.Pipe <- core.Message{Command: "listall", Tags: msg.Tags, Data: subscriptions}
	return nil
}

// handleWho - who related
func handleWho(proxy *core.Proxy, msg core.Message) error {
	who, max := core.Who()
	subscribers := fmt.Sprintf("Lifetime  connections: %d\nSubscribers connected: %d", max, who)
	proxy.Pipe <- core.Message{Command: "who", Tags: msg.Tags, Data: subscribers}
	return nil
}
