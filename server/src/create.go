package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create
type createRequestMsg struct {
	RoomName string `json:"roomName"`
	Name     string `json:"name"`
	Token    Token  `json:"token"`
}

func CreateHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		Log(err.Error())
		return
	}

	req := createRequestMsg{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		Log(err.Error())
		return
	}
	createResponse(c, req)
}

type createResponseMsg struct {
	RoomId  string `json:"roomId"`
	Success bool   `json:"success"`
}

/*
 /create {roomName: "", name: "", token: %TOKEN%}
 return {roomId: ""}
*/
func createResponse(c *gin.Context, msg createRequestMsg) {
	room := NewRoom(msg.RoomName)
	user, ok := USERS[msg.Token]
	roomAdded := AddRoom(room)
	success := ok && roomAdded
	if success {
		AddUserToRoom(room, user)
		LogRoom(room.roomId)
	}

	resp := gin.H{
		"roomId":           room.roomId,
		"success":          success,
		"roomName":         room.roomName,
		"previousMessages": room.messages,
	}
	c.JSON(http.StatusOK, resp)
}
