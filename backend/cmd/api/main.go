package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// Upgrader to upgrade HTTP connection to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for dev
	},
}

// Store connected clients
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan string)

func main() {
	r := gin.Default()

	// 1. WebSocket Endpoint for Live Match Updates
	r.GET("/live-updates", handleWebSocket)

	// 2. RAG Chat Endpoint (The AI Logic)
	r.POST("/ask-assistant", handleRAGChat)

	// Start a goroutine to simulate live match events
	go simulateMatchEvents()
	
	// Start handling broadcasts
	go handleMessages()

	fmt.Println("Server starting on :8080...")
	r.Run(":8080")
}

func handleWebSocket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	// Register client
	clients[ws] = true

	for {
		// Keep connection open and read (even if we don't use client msgs yet)
		_, _, err := ws.ReadMessage()
		if err != nil {
			delete(clients, ws)
			break
		}
	}
}

func handleMessages() {
	for {
		// Grab message from broadcast channel
		msg := <-broadcast
		// Send to every connected client
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// Simulates fetching data from Football-API and pushing to clients
func simulateMatchEvents() {
	events := []string{
		"MATCH START: Arsenal vs Man City",
		"12' - Attempt missed. Bukayo Saka (Arsenal) left footed shot...",
		"24' - GOAL! Haaland scores for Man City!",
		"45' - Half Time. City leads 1-0.",
	}

	for _, event := range events {
		time.Sleep(5 * time.Second) // Simulate time passing
		broadcast <- event
	}
}

type ChatRequest struct {
	UserQuery string `json:"query"`
}

// This is where your RAG logic lives
func handleRAGChat(c *gin.Context) {
	var req ChatRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Request"})
		return
	}

	// TODO: 1. Fetch relevant vector embeddings (Rules/Player Stats)
	// TODO: 2. Send context + req.UserQuery to OpenAI/Gemini API
	
	// Mock Response for now
	fakeAIResponse := fmt.Sprintf("AI Insight: You asked about '%s'. Based on the current match state (City 1-0), Haaland is exploiting the high line.", req.UserQuery)

	c.JSON(http.StatusOK, gin.H{"response": fakeAIResponse})
}