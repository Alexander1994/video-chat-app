package main

// https://github.com/Durgaprasad-Budhwani/docker-azure-web-app-golang

import (
	"errors"
	"log"
	"mime"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/gorilla/websocket"
)

/*
 %TOKEN% => {conn, name, roomId}
 roomId => {user %TOKEN%s, roomName}

 /data {typ: JOIN,    message: %TOKEN% }
 /data {typ: MESSAGE, message: "" }
 /data {typ: LEAVE,   message: %TOKEN% }
*/

type Token = string

type User struct {
	token  Token
	name   string
	roomId string
	conn   *websocket.Conn
}

var UserConnInvalid = errors.New("User conn invalid")

func (user *User) String() string {
	str := user.token + " " + user.name
	room, ok := ROOMS[user.roomId]
	roomStr := "room: (invalid room)"
	if ok {
		roomStr = "room: (" + room.roomId + " " + room.roomName + ")"
	}
	str += " " + roomStr
	return str
}

var userID uint64 = 0
var NullUserToken string = strconv.FormatUint(userID, 10)
var NullUser = User{NullUserToken, "", NullRoomId, nil}

func newUserId() Token {
	userID++
	return strconv.FormatUint(userID, 10)
}

func NewUser(name string) User {
	return User{newUserId(), name, NullRoomId, nil}
}

type Room struct {
	roomId     string
	roomName   string
	userTokens []Token
	messages   []Message
}

var UserNotFoundInRoomError = errors.New("User not found in room")

func (room *Room) String() string {
	str := room.roomId + " " + room.roomName + " users: ("
	if len(room.userTokens) == 0 {
		str += ")"
	} else {
		lastElIndex := len(room.userTokens) - 1
		for i := 0; i < lastElIndex; i++ {
			userToken := room.userTokens[i]
			if user, ok := USERS[userToken]; ok {
				str += user.String() + ", "
			} else {
				str += "invalid user, "
			}
		}
		u := USERS[room.userTokens[lastElIndex]]
		str += u.String() + ")"
	}
	return str
}

var roomId uint64 = 0
var NullRoomId string = strconv.FormatUint(roomId, 10)
var NullRoom = Room{NullRoomId, "", nil, nil}

func newRoomId() Token {
	roomId++
	return strconv.FormatUint(roomId, 10)
}

func NewRoom(roomName string) Room {
	return Room{newRoomId(), roomName, make([]Token, 0), make([]Message, 0)}
}

var CONNECTIONS map[*websocket.Conn]Token
var USERS map[Token]User
var ROOMS map[string]Room

var ConnectionNotFoundError = errors.New("Connection Not Found")
var UserNotFoundError = errors.New("User Not Found")
var RoomNotFoundError = errors.New("Room Not Found")

const SvcId string = "server"

func loggerPrefix() string {
	return SvcId + " " // + " " + time.Now().Format(time.RFC3339) +
}

func Log(s string) {
	log.Println(s) // loggerPrefix() +
}

func LogConnection(conn *websocket.Conn) {
	token, ok := CONNECTIONS[conn]
	if !ok {
		token = NullUserToken
	}
	Log(token)
}

func LogConnections() {
	for conn := range CONNECTIONS {
		LogConnection(conn)
	}
}

func LogRoom(roomId string) {
	room, ok := ROOMS[roomId]
	roomStr := "invalid room"
	if ok {
		roomStr = room.String()
	}
	Log(roomStr)
}

func LogRooms() {
	for room := range ROOMS {
		LogRoom(room)
	}
}

func LogUser(token Token) {
	user, ok := USERS[token]
	userStr := "invalid user"
	if ok {
		userStr = user.String()
	}
	Log(userStr)
}

func LogUsers() {
	for user := range USERS {
		LogUser(user)
	}
}

func AddConnection(conn *websocket.Conn, token Token) (ok bool) {
	var user User
	connFound := false
	if user, ok = USERS[token]; ok {
		if _, connFound = CONNECTIONS[conn]; !connFound {
			CONNECTIONS[conn] = token
			user.conn = conn
			USERS[token] = user
		}
	}

	return ok && !connFound
}

func AddUser(u User) (ok bool) {
	if _, ok = USERS[u.token]; !ok {
		USERS[u.token] = u
	}
	return !ok
}

func GetUser(conn *websocket.Conn) (user User, ok bool) {
	if userID, ok := CONNECTIONS[conn]; ok {
		user, ok := USERS[userID]
		return user, ok
	}
	return NullUser, false
}

func FindUser(name string) (u User, ok bool) {
	for _, user := range USERS {
		if user.name == name {
			return user, true
		}
	}
	return NullUser, false
}

func RemoveUser(token string) {
	_, ok := USERS[token]
	if ok {
		delete(USERS, token)
	}
}

func RemoveConn(conn *websocket.Conn) {
	_, ok := CONNECTIONS[conn]
	if ok {
		delete(CONNECTIONS, conn)
	}
}

func AddRoom(r Room) (ok bool) {
	if _, ok = ROOMS[r.roomId]; !ok {
		ROOMS[r.roomId] = r
	}
	return !ok
}

func FindRoomByName(roomName string) (room Room, ok bool) {
	for _, room = range ROOMS {
		if room.roomName == roomName {
			return room, true
		}
	}
	return NullRoom, false
}

func AddUserToRoom(room Room, user User) {
	user.roomId = room.roomId
	USERS[user.token] = user

	room.userTokens = append(room.userTokens, user.token)
	ROOMS[room.roomId] = room
}

func RemoveRoom(roomId string) {
	_, ok := ROOMS[roomId]
	if ok {
		delete(ROOMS, roomId)
	}
}

// Helper functions
func jsTimeStamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func AddMsgToRoom(roomId string, msg Message) (ok bool) {
	var room Room
	if room, ok = ROOMS[roomId]; ok {
		room.messages = append(room.messages, msg)
		ROOMS[roomId] = room
	}
	return ok
}

func broadCastMessageToRoom(roomId string, msg Message, onConn func(*websocket.Conn)) {
	room := ROOMS[roomId]
	room.messages = append(room.messages, msg)
	ROOMS[roomId] = room
	for _, token := range room.userTokens {
		if user, ok := USERS[token]; ok {
			onConn(user.conn)
		}
	}
}

func forEachUsertoken(userTokens []string, cb func(u User)) {
	for _, token := range userTokens {
		if user, ok := USERS[token]; ok {
			cb(user)
		}
	}
}

func setup() {
	USERS = make(map[Token]User)
	ROOMS = make(map[string]Room)
	CONNECTIONS = make(map[*websocket.Conn]Token)

	mime.AddExtensionType(".js", "text/javascript")
	mime.AddExtensionType(".css", "text/css")

	u := NewUser("tod")
	r := NewRoom("cool")
	AddUser(u)
	AddRoom(r)
}

func CORSMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}

var sslkey string = "server.key"
var sslcert string = "server.crt"

/*
 user logs in: 'connection
 user attemps join room: 'join-room' (roomId, jwt) => jwt?
 user disconnects: 'disconnect'
*/
// https://github.com/googollee/go-socket.io/blob/master/_examples/gin-gonic/main.go
// https://github.com/WebDevSimplified/Zoom-Clone-With-WebRTC
func main() {
	setup()
	isHttps := false
	isDevelopment := true
	loadDir := "../../client/dist/video-app/"

	PORT := ":8000"
	if isHttps {
		PORT = ":443"
	}

	router := gin.Default()

	if isDevelopment {
		router.Use(CORSMiddleware("http://localhost:3000"))
	}

	api := router.Group("/api")
	{
		api.POST("/login", LoginHandler)
		api.POST("/room/join", JoinHandler)
		api.POST("/room/create", CreateHandler)
		api.DELETE("/room", RoomDeleteHandler)
	}

	websocketRouter := CreateSocketHandlers()

	router.GET("/messages", func(c *gin.Context) {
		websocketRouter.RequestReceiver(c.Writer, c.Request)
	})

	router.Use(static.Serve("/", static.LocalFile(loadDir, true)))

	var err error
	if isHttps {
		err = router.RunTLS(PORT, sslcert, sslkey)
	} else {
		err = router.Run(PORT)
	}
	if err != nil {
		log.Fatal("failed run app: ", err)
	}

	Log("Signaling Server started")
}
