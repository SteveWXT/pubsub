package server

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/jcelliott/lumber"
)

var (
	ErrNotImplemented = fmt.Errorf("Error: Not Implemented\n")

	// this is a map of the supported servers that can be started by mist
	servers    = map[string]handleFunc{}
	serversTex sync.RWMutex
)

type (
	handleFunc func(uri string, errChan chan<- error)
)

// Register registers a new pubsub server
func Register(name string, auth handleFunc) {
	serversTex.Lock()
	servers[name] = auth
	serversTex.Unlock()
}

// Start attempts to individually start servers from a list of provided
// listeners; the listeners provided is a comma delimited list of uri strings
// (scheme:[//[user:pass@]host[:port]][/]path[?query][#fragment])
func Start(uris []string) error {

	// this chan is given to each individual server start as a way for them to
	// communicate back their startup status
	errChan := make(chan error, len(uris))

	// iterate over each of the provided listener uris attempting to start them
	// individually; if one isn't supported it gets skipped
	for i := range uris {

		// parse the uri string into a url object
		url, err := url.Parse(uris[i])
		if err != nil {
			return err
		}

		// check to see if the scheme is supported; if not, indicate as such and
		// continue
		serversTex.RLock()
		server, ok := servers[url.Scheme]
		serversTex.RUnlock()
		if !ok {
			lumber.Error("Unsupported scheme '%s'", url.Scheme)
			continue
		}

		// attempt to start the server
		lumber.Info("Starting '%s' server...", url.Scheme)
		go server(url.Host, errChan)
	}

	// handle errors that happen during startup by reading off errChan and returning
	// on any error received. If no errors are received after 1 second per server
	// assume successful starts.
	select {
	case err := <-errChan:
		lumber.Error("Failed to start - %s", err.Error())
		return err
	case <-time.After(time.Second * time.Duration(len(uris))):
		// no errors
	}

	// handle errors that happen after initial start; if any errors are received they
	// are logged and the servers just try to keep running
	for err := range errChan {
		// log these errors and continue
		lumber.Error("Server error - %s", err.Error())
	}

	return nil
}

// StartWithLS attempts to individually start servers from a TCP listeners
func StartWithLS(ls net.Listener) error {

	// this chan is given to each individual server start as a way for them to
	// communicate back their startup status
	errChan := make(chan error, 1)

	// attempt to start the server
	go StartTCPWithLS(ls, errChan)

	// handle errors that happen during startup by reading off errChan and returning
	// on any error received. If no errors are received after 1 second per server
	// assume successful starts.
	select {
	case err := <-errChan:
		lumber.Error("Failed to start - %s", err.Error())
		return err
	case <-time.After(time.Second * time.Duration(1)):
		// no errors
	}

	// handle errors that happen after initial start; if any errors are received they
	// are logged and the servers just try to keep running
	for err := range errChan {
		// log these errors and continue
		lumber.Error("Server error - %s", err.Error())
	}

	return nil
}

// GetPort return the port number of a listener
func GetPort(lis net.Listener) (uint32, error) {
	_, portStr, err := net.SplitHostPort(lis.Addr().String())
	if err != nil {
		return 0, err
	}
	port, err := strconv.ParseUint(portStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(port), nil
}
