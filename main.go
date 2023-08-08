package main

import (
	"go-chat/chat"
	"go-chat/handlers"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", &handlers.TemplateHandler{File: "chat.html"})

	// public room
	publicRoom := chat.NewRoom()
	// execute the room as a goroutine
	go publicRoom.Run()
	http.Handle("/ws/public", publicRoom)

	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Fatal(err)
	}
}
