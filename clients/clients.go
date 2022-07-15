package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/jcelliott/lumber"

	"github.com/SteveWXT/pubsub/core"
)

type (
	// TCP represents a TCP connection to the core server
	TCP struct {
		conn     io.ReadWriteCloser // the connection to the core server
		encoder  *json.Encoder      //
		host     string             //
		messages chan core.Message  // the channel that core server 'publishes' updates to
	}
)

// New attempts to connect to a running core server at the clients specified
// host and port.
func New(host string) (*TCP, error) {
	client := &TCP{
		host:     host,
		messages: make(chan core.Message),
	}

	return client, client.connect()
}

// connect dials the remote core server and handles any incoming responses back
// from core
func (c *TCP) connect() error {

	// attempt to connect to the server
	conn, err := net.Dial("tcp", c.host)
	if err != nil {
		return fmt.Errorf("Failed to dial '%s' - %s", c.host, err.Error())
	}

	// set the connection for the client
	c.conn = conn

	// create a new json encoder for the clients connection
	c.encoder = json.NewEncoder(c.conn)

	// ensure we are authorized/still connected (unauthorized clients get disconnected)
	c.Ping()
	decoder := json.NewDecoder(conn)
	msg := core.Message{}
	if err := decoder.Decode(&msg); err != nil {
		conn.Close()
		close(c.messages)
		return fmt.Errorf("Ping failed, possibly bad token, or can't read from core - %s", err.Error())
	}

	// connection loop (blocking); continually read off the connection. Once something
	// is read, check to see if it's a message the client understands to be one of
	// its commands. If so attempt to execute the command.
	go func() {

		for {
			msg := core.Message{}

			// decode an array value (Message)
			if err := decoder.Decode(&msg); err != nil {
				switch err {
				case io.EOF:
					lumber.Debug("[pubsub client] pubsub terminated connection")
				case io.ErrUnexpectedEOF:
					lumber.Debug("[pubsub client] pubsub terminated connection unexpectedly")
				default:
					lumber.Error("[pubsub client] Failed to get message from pubsub - %s", err.Error())
				}
				conn.Close()
				close(c.messages)
				return
			}
			c.messages <- msg // read from this using the .Messages() function
			lumber.Trace("[pubsub client] Received message - %#v", msg)
		}
	}()

	return nil
}

// Ping the server
func (c *TCP) Ping() error {
	return c.encoder.Encode(&core.Message{Command: "ping"})
}

// Subscribe takes the specified tags and tells the server to subscribe to updates
// on those tags, returning the tags and an error or nil
func (c *TCP) Subscribe(tags []string) error {

	if len(tags) == 0 {
		return fmt.Errorf("Unable to subscribe - missing tags")
	}

	return c.encoder.Encode(&core.Message{Command: "subscribe", Tags: tags})
}

// Unsubscribe takes the specified tags and tells the server to unsubscribe from
// updates on those tags, returning an error or nil
func (c *TCP) Unsubscribe(tags []string) error {

	if len(tags) == 0 {
		return fmt.Errorf("Unable to unsubscribe - missing tags")
	}

	return c.encoder.Encode(&core.Message{Command: "unsubscribe", Tags: tags})
}

// Publish sends a message to the core server to be published to all subscribed
// clients
func (c *TCP) Publish(tags []string, data string) error {

	if len(tags) == 0 {
		return fmt.Errorf("Unable to publish - missing tags")
	}

	if data == "" {
		return fmt.Errorf("Unable to publish - missing data")
	}

	return c.encoder.Encode(&core.Message{Command: "publish", Tags: tags, Data: data})
}

// PublishAfter sends a message to the core server to be published to all subscribed
// clients after a specified delay
func (c *TCP) PublishAfter(tags []string, data string, delay time.Duration) error {
	go func() {
		<-time.After(delay)
		c.Publish(tags, data)
	}()
	return nil
}

// List requests a list from the server of the tags this client is subscribed to
func (c *TCP) List() error {
	return c.encoder.Encode(&core.Message{Command: "list"})
}

// listall related
// List requests a list from the server of the tags this client is subscribed to
func (c *TCP) ListAll() error {
	return c.encoder.Encode(&core.Message{Command: "listall"})
}

// who related
// Who requests connection/subscriber stats from the server
func (c *TCP) Who() error {
	return c.encoder.Encode(&core.Message{Command: "who"})
}

// Close closes the client data channel and the connection to the server
func (c *TCP) Close() {
	c.conn.Close()
	// close(c.messages) // we don't close this in case there is a message waiting in the channel
}

// Messages ...
func (c *TCP) Messages() <-chan core.Message {
	return c.messages
}
