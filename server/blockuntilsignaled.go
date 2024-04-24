package server

import (
	"os"
	"os/signal"
)

// BlockUntilSignaled blocks until one of the given signals is received.
// Returns the received signal.
func BlockUntilSignaled(first os.Signal, more ...os.Signal) os.Signal {
	// Before we start the server, sign up for system signals that indicate that we should
	// gracefully shutdown. We set our buffer size to 1 since we are expecting just one signal.
	// See https://golang.org/pkg/os/signal/#Notify.
	interrupt := make(chan os.Signal, 1)
	defer close(interrupt)
	signal.Notify(interrupt, append(more, first)...)

	// Block until we get a signal from the interrupt channel.
	sig := <-interrupt

	return sig
}
