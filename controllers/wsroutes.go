package controllers

import (
	"coeus/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Connection struct {
	conn   *websocket.Conn
	send   chan []byte
	closed chan bool
}

// Connections for all changes to the database
var generalConnections = make(map[*Connection]bool)

// A map of class session ids to a map of connections
var connections = make(map[int]map[*Connection]bool)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// BROADCAST FUNCTIONS START

// Broadcasts a message to all active connections
func generalBroadcast(message []byte) {
	for connection := range generalConnections {
		select {
		case connection.send <- message:
		default:
			delete(generalConnections, connection)
			close(connection.send)
		}
	}
}

// Broadcasts a message to all active connections for a specific class session
func classSessionBroadcast(classSessionID int, message []byte) {
	if conns, ok := connections[classSessionID]; ok {
		for connection := range conns {
			select {
			case connection.send <- message:
			default:
				delete(conns, connection)
				close(connection.send)
			}
		}
	}
}

func constructVoteUp(classSessionID, userID, questionID, updatedVoteCount int) {
	// Construct a JSON object
	voteUp := map[string]interface{}{
		"action":     "vote-up",
		"userID":     userID,
		"questionID": questionID,
		"votes":      updatedVoteCount,
	}
	// Marshal the voteUp message
	voteUpBytes, err := json.Marshal(voteUp)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return
	}
	// Broadcast the "vote-up" action to all active connections
	classSessionBroadcast(classSessionID, voteUpBytes)
}

func broadcastNewQuestion(classSessionID int, question models.Question) {
	// Construct a JSON object
	newQuestion := map[string]interface{}{
		"action":     "new-question",
		"userID":     question.UserID,
		"questionID": question.ID,
		"sessionID":  question.SessionID,
		"text":       question.Text,
		"votes":      question.Votes,
		"answered":   question.Answered,
		"createdAt":  question.CreatedAt,
	}
	// Marshal the newQuestion message
	questionBytes, err := json.Marshal(newQuestion)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return
	}
	// Broadcast the newQuestion message to all active connections
	classSessionBroadcast(classSessionID, questionBytes)
}

func constructMarkQuestion(classSessionID, questionID int) {
	// Construct a JSON object
	markQuestion := map[string]interface{}{
		"action":     "mark-question",
		"questionID": questionID,
		"answered":   true,
	}
	// Marshal the markQuestion message
	markQuestionBytes, err := json.Marshal(markQuestion)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return
	}
	// Broadcast the "mark-question" action to all active connections
	classSessionBroadcast(classSessionID, markQuestionBytes)
}

func constructStartSession(sectionID, attendanceID int) {
	// Construct a JSON object
	startSession := map[string]interface{}{
		"action":       "start-session",
		"sectionID":    sectionID,
		"attendanceID": attendanceID,
	}
	// Marshal the startSession message
	startSessionBytes, err := json.Marshal(startSession)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return
	}
	// Broadcast the "start-session" action to all active connections
	generalBroadcast(startSessionBytes)
	classSessionBroadcast(sectionID, startSessionBytes)
}

func constructEndSession(classSessionID, sectionID int) {
	// Construct a JSON object
	endSession := map[string]interface{}{
		"action":    "end-session",
		"sectionID": sectionID,
	}
	// Marshal the endSession message
	endSessionBytes, err := json.Marshal(endSession)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return
	}
	// Broadcast the "end-session" action to all active connections
	generalBroadcast(endSessionBytes)
	classSessionBroadcast(classSessionID, endSessionBytes)
}

func constructParticipantJoined(classSessionID, participantCount int) {
	// Construct a JSON object
	participantJoined := map[string]interface{}{
		"action": "participant-joined",
		"count":  participantCount,
	}
	// Marshal the participantJoined message
	participantJoinedBytes, err := json.Marshal(participantJoined)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		return
	}
	// Broadcast the "participant-joined" action to all active connections
	classSessionBroadcast(classSessionID, participantJoinedBytes)
}

// BROADCAST FUNCTIONS END

// DATA PUMPS START
func (c *Connection) writePump() {
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.closed <- true
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c *Connection) readPump() {
	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		// Unmarshal the received JSON message
		var input map[string]interface{}
		if err := json.Unmarshal(messageBytes, &input); err != nil {
			log.Println("Error unmarshalling JSON:", err)
			break
		}
	}
	c.closed <- true
}

// DATA PUMPS END

func WSRoutes(g *gin.RouterGroup) {
	g.GET("/ws", HandleGeneralWebsocketConnection)
	g.GET("/ws/:classSessionID", HandleClassWebsocketConnection)
}
