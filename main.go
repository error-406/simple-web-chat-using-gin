package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
)

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
	Color string `json:"color"`
}
var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel
var upgrader = websocket.Upgrader{}
var sender *websocket.Conn


func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	  ws, err := upgrader.Upgrade(w, r, nil)
	  if err != nil {
		log.Fatal(err)
	  }
	  // Make sure we close the connection when the function returns
	  defer ws.Close()
	  // Register our new client
	  clients[ws] = true
	  for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		sender = ws
		err := ws.ReadJSON(&msg)
		if err != nil {
		  log.Printf("error: %v", err)
		  delete(clients, ws)
		  break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	  }
  }
  
  func handleMessages() {
	  for {
		  // Grab the next message from the broadcast channel
		  msg := <-broadcast
		  // Send it out to every client that is currently connected
		  for client := range clients {
			if(client != sender){
					err := client.WriteJSON(msg)
					if err != nil {
						log.Printf("error: %v", err)
						client.Close()
						delete(clients, client)
					}
			}
		  }
	  }
  }

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.LoadHTMLGlob("html/*")
	// Ping test
	
	r.GET("/chat", func(c *gin.Context) {
		//c.String(http.StatusOK, "chat")
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Ruang Julid Terbuka",
		})
	})

	r.GET("/ws", func(c *gin.Context) {
		handleConnections(c.Writer, c.Request)
	})
	//http.HandleFunc("/ws", handleConnections)
	return r
}

func main() {
	r := setupRouter()
	go handleMessages()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8082")
}
