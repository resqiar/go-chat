package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", helloWorld)

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal(err)
	}
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}
