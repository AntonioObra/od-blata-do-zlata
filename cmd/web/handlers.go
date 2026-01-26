/**
*
*  /\_/\
* ( o.o )
*  > ^ <
*
* Handlers
*
 */

package main

import (
	"fmt"
	"net/http"
	"strconv"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")
	w.Write([]byte("Hello from Od blata do zlata"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "Display a snippet with ID %d...", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating new snippet"))
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)

	w.Write([]byte("Create a new snippet"))
}
