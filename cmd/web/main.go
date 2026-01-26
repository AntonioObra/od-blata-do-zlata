/**
*
*  /\_/\
* ( o.o )
*  > ^ <
*
* Main
*
 */

package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// Define a new command-line flag with the name "addr"
	addr := flag.String("addr", ":4000", "HTTP network address")

	// ! Important, to use flag.Parse() method to parse the command-line flag
	// ! flag.Parse() returns values as a pointer, so we need to dereference it with *
	flag.Parse()

	// Setup custom logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()

	// Create a file server which servers files out of the "./ui/static" dir
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home) // Restrict this route to exact matches on / only.
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	// Print a log message to say that the server is starting.
	logger.Info("starting server", "addr", *addr)

	// Use the http.ListenAndServe() function to start a new web server. We pass in
	// two parameters: the TCP network address to listen on (in this case ":4000")
	// and the servemux we just created. If http.ListenAndServe() returns an error
	// we use the log.Fatal() function to log the error message and terminate the
	// program. Note that any error returned by http.ListenAndServe() is always
	// non-nil.
	err := http.ListenAndServe(*addr, mux)

	// * Log the error using custom logger and terminate the application
	logger.Error(err.Error())
	os.Exit(1)
}
