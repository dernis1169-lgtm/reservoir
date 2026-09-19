package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	socketio "github.com/googollee/go-socket.io"
)

func main() {
	server := socketio.NewServer(nil)

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "joinRoom", func(s socketio.Conn, roomId string) {
		s.Join(roomId)
		fmt.Printf("User %s joined room %s\n", s.ID(), roomId)
	})

	server.OnEvent("/", "draw", func(s socketio.Conn, data interface{}) {
		// Правильный способ рассылки в Go Socket.IO
		if m, ok := data.(map[string]interface{}); ok {
			if roomId, ok := m["roomId"].(string); ok {
				server.BroadcastToRoom("/", roomId, "draw", data)
			}
		}
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()
	defer server.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.Handle("/socket.io/", server)
	log.Printf("Server started on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
