package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Join
type loginRequestMsg struct {
	Name string `json:"name"`
}

func LoginHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		Log(err.Error())
		return
	}

	req := loginRequestMsg{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		Log(err.Error())
		return
	}
	loginResponse(c, req)
}

type loginResponseMsg struct {
	name    string `json:"name"`
	token   string `json:"token"`
	success bool   `json:"success"`
}

/*
 /login {name: "" }
 return {success bool, token: %TOKEN%}
*/
func loginResponse(c *gin.Context, msg loginRequestMsg) {
	user, ok := FindUser(msg.Name)

	tokenStr, err := CreateToken(msg.Name, user.token)

	if err == nil {
		resp := gin.H{
			"name":    user.name,
			"token":   tokenStr,
			"success": ok,
		}
		c.JSON(http.StatusOK, resp)
	}
}
