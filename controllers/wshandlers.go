package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HandleGeneralWebsocketConnection handles the websocket connection for general changes to the database
func HandleGeneralWebsocketConnection(c *gin.Context) {
	w := c.Writer
	r := c.Request
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}

	connection := &Connection{conn: conn, send: make(chan []byte), closed: make(chan bool, 2)}

	generalConnections[connection] = true

	go connection.writePump()
	go connection.readPump()

	<-connection.closed
	<-connection.closed

	delete(generalConnections, connection)
	conn.Close()
}

// HandleClassWebsocketConnection handles the websocket connection for a specific class session
func HandleClassWebsocketConnection(c *gin.Context) {
	// Extract classSessionID from the request
	classSessionID, err := strconv.Atoi(c.Param("classSessionID"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract class session ID"})
		return
	}

	w := c.Writer
	r := c.Request
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}

	connection := &Connection{conn: conn, send: make(chan []byte), closed: make(chan bool, 2)}

	// Add the connection to the appropriate class session
	if _, ok := connections[classSessionID]; !ok {
		connections[classSessionID] = make(map[*Connection]bool)
	}
	connections[classSessionID][connection] = true

	go connection.writePump()
	go connection.readPump()

	<-connection.closed
	<-connection.closed

	// Remove the connection from the class session
	delete(connections[classSessionID], connection)
	conn.Close()
}
