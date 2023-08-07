package main

import (
	"go-chat/handlers"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", &handlers.TemplateHandler{File: "chat.html"})

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal(err)
	}
}
