/**
*
* Main
*
*  /\_/\
* ( o.o )
*  > ^ <
*
 */

package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

type application struct {
	logger *slog.Logger
}

func main() {
	// Define a new command-line flag with the name "addr"
	addr := flag.String("addr", ":4000", "HTTP network address")

	// ! Important, to use flag.Parse() method to parse the command-line flag
	// ! flag.Parse() returns values as a pointer, so we need to dereference it with *
	flag.Parse()

	// Setup custom logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	// Initialize a new instance of our application struct, containing
	// dependencies (for noew, just the structured logger)
	app := &application{
		logger: logger,
	}

	// Print a log message to say that the server is starting.
	app.logger.Info("starting server", "addr", *addr)

	// Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address to listen on (in this case ":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and terminate the
	// program. Note that any error returned by http.ListenAndServe() is always
	// non-nil.
	err := http.ListenAndServe(*addr, app.routes())

	// * Log the error using custom logger and terminate the application
	app.logger.Error(err.Error())
	os.Exit(1)
}
