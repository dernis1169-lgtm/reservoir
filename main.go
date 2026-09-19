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
		fmt.Println("New connection:", s.ID())
		return nil
	})

	server.OnEvent("/", "joinRoom", func(s socketio.Conn, roomId string) {
		s.Join(roomId)
		fmt.Printf("User %s joined room %s\n", s.ID(), roomId)
		// Уведомляем пользователя, что он успешно вошел
		s.Emit("joined", roomId)
	})

	server.OnEvent("/", "draw", func(s socketio.Conn, data interface{}) {
		// Рассылаем ВСЕМ в этой комнате, включая отправителя (для синхронизации ID)
		if m, ok := data.(map[string]interface{}); ok {
			if roomId, ok := m["roomId"].(string); ok {
				server.BroadcastToRoom("/", roomId, "draw", data)
			}
		}
	})

	server.OnEvent("/", "reaction", func(s socketio.Conn, data interface{}) {
		if m, ok := data.(map[string]interface{}); ok {
			if roomId, ok := m["roomId"].(string); ok {
				server.BroadcastToRoom("/", roomId, "reaction", data)
			}
		}
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("Error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("Closed:", reason)
	})

	go server.Serve()
	defer server.Close()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.Handle("/socket.io/", server)
	log.Printf("Server running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
