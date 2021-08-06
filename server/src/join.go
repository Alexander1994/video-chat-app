package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Join
type joinRequestMsg struct {
	RoomName string `json:"roomName"`
}

func JoinHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		Log(err.Error())
		return
	}

	jwtClaims, authenticated := Authenticate(c.Request.Header)
	if !authenticated {
		return
	}

	req := joinRequestMsg{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		Log(err.Error())
		return
	}
	joinResponse(c, req, jwtClaims)
}

type joinResponseMsg struct {
	roomName string `json:"roomName"`
	roomId   string `json:"roomId"`
	success  bool   `json:"response"`
}

/*
 /join {roomId: "", name: "", token: %TOKEN%}
 return {success bool, roomName: ""}
*/
func joinResponse(c *gin.Context, msg joinRequestMsg, jwtClaims *JwtLoginClaims) {
	room, roomFound := FindRoomByName(msg.RoomName)
	token := jwtClaims.UserToken
	user, userFound := USERS[token]
	LogUser(token)

	success := roomFound && userFound
	if success {
		AddUserToRoom(room, user)
		LogRoom(room.roomId)
	} else {
		room = NullRoom
	}

	resp := gin.H{
		"roomName":         room.roomName,
		"roomId":           room.roomId,
		"success":          success,
		"previousMessages": room.messages,
	}
	c.JSON(http.StatusOK, resp)

}
