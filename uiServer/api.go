package uiServer

import (
	"fmt"
	"net/http"
)

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
