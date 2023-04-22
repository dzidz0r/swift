package cmd

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection:", err)
		return
	}
	fmt.Println("Web socket connectin established")
	conn.WriteMessage(1, []byte("This is the socket"))
	defer conn.Close()

	for {
		// Read message from WebSocket client
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Failed to read message:", err)
			break
		}

		// Print received message to console
		fmt.Println("Received message:", string(message))

		// Write message back to WebSocket client
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("Failed to write message:", err)
			break
		}
	}
}

func Main() {
	// Serve the frontend directory
	fs := http.FileServer(http.Dir("ui"))
	http.Handle("/", fs)

	// Handle WebSocket connections
	http.HandleFunc("/ws", handleWebSocket)

	fmt.Println("Starting server on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
}
