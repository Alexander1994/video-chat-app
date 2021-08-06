package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func videoConnHandler(conn *websocket.Conn) (err error) {
	var raw []byte
	_, raw, err = conn.ReadMessage()
	if err != nil {
		return err
	}
	if userToken, ok := CONNECTIONS[conn]; ok {
		if user, ok := USERS[userToken]; ok {
			err = broadcastVideo(user.roomId, raw)
		} else {
			return RoomNotFoundError
		}
	} else {
		return UserNotFoundError
	}
	return err
}

func broadcastVideo(roomId string, out []byte) (err error) {
	if room, found := ROOMS[roomId]; found {
		for _, userId := range room.userTokens {
			err = sendVideoData(userId, out)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func sendVideoData(userToken Token, out []byte) (err error) {
	if userOut, ok := USERS[userToken]; ok {
		if userOut.conn == nil {
			log.Println("Error - connHandler - Programmer error, user conn not properly set")
			return err
		}
		if err = userOut.conn.WriteMessage(1, out); err != nil {
			log.Println("Error - connHandler - WriteMessage Response:", err)
			return err
		}
	}
	return err
}

func SocketVideoHandler(c *gin.Context) {
	writer := c.Writer
	req := c.Request

	conn, err := upgrader.Upgrade(writer, req, nil)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println(conn.RemoteAddr(), "reached the server")

	for {
		if err := videoConnHandler(conn); err != nil && err.Error() == "websocket: close 1001 (going away)" {
			if user, isRegistered := CONNECTIONS[conn]; isRegistered {
				log.Println("Connection closed for", user)
				forceLeave(conn)
			} else {
				log.Println(conn.RemoteAddr(), "reached the server")
			}
			return
		}
	}
}
