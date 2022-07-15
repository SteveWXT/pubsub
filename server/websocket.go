package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/pat"
	"github.com/gorilla/websocket"
	"github.com/jcelliott/lumber"

	"github.com/SteveWXT/pubsub/core"
)

// init adds ws/wss as available core server types
func init() {
	Register("ws", StartWS)
}

// StartWS starts a core server listening over a websocket
func StartWS(uri string, errChan chan<- error) {
	router := pat.New()
	router.Get("/subscribe/websocket", func(rw http.ResponseWriter, req *http.Request) {

		// prepare to upgrade http to ws
		upgrader := websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		}

		// upgrade to websocket conn
		conn, err := upgrader.Upgrade(rw, req, nil)
		if err != nil {
			errChan <- fmt.Errorf("Failed to upgrade connection - %s", err.Error())
			return
		}
		defer conn.Close()

		proxy := core.NewProxy()
		defer proxy.Close()

		// add basic WS handlers for this socket
		handlers := GenerateHandlers()

		// read and publish core messages to connected clients (non-blocking)
		go func() {
			for msg := range proxy.Pipe {

				// failing to write is probably because the connection is dead; we dont
				// want core just looping forever tyring to write to something it will
				// never be able to.
				if err := conn.WriteJSON(msg); err != nil {
					if err.Error() != "websocket: close sent" {
						errChan <- fmt.Errorf("Failed to WriteJSON message to WS connection - %s", err.Error())
					}

					break
				}
			}
		}()

		// connection loop (blocking); continually read off the connection. Once something
		// is read, check to see if it's a message the client understands to be one of
		// its commands. If so attempt to execute the command.
		for {

			msg := core.Message{}

			// failing to read is probably because the connection is dead; we dont
			// want core just looping forever tyring to write to something it will
			// never be able to.
			if err := conn.ReadJSON(&msg); err != nil {
				// todo: better logging here too
				if !strings.Contains(err.Error(), "websocket: close 1001") &&
					!strings.Contains(err.Error(), "websocket: close 1005") &&
					!strings.Contains(err.Error(), "websocket: close 1006") { // don't log if client disconnects
					errChan <- fmt.Errorf("Failed to ReadJson message from WS connection - %s", err.Error())
				}

				break // todo: continue?
			}

			// look for the command
			handler, found := handlers[msg.Command]

			// if the command isn't found, return an error
			if !found {
				lumber.Trace("Command '%s' not found", msg.Command)
				if err := conn.WriteJSON(&core.Message{Command: msg.Command, Error: "Unknown Command"}); err != nil {
					errChan <- fmt.Errorf("WS Failed to respond to client with 'command not found' - %s", err.Error())
				}
				continue
			}

			// attempt to run the command
			lumber.Trace("WS Running '%s'...", msg.Command)
			if err := handler(proxy, msg); err != nil {
				lumber.Debug("WS Failed to run '%s' - %s", msg.Command, err.Error())
				if err := conn.WriteJSON(&core.Message{Command: msg.Command, Error: err.Error()}); err != nil {
					errChan <- fmt.Errorf("WS Failed to respond to client with error - %s", err.Error())
				}
				continue
			}
		}
	})

	lumber.Info("WS server listening at '%s'...\n", uri)
	// go http.ListenAndServe(uri, router)
	http.ListenAndServe(uri, router)
}
