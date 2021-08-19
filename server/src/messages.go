package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type MessageType string

const (
	JOIN    MessageType = "join"
	LEAVE   MessageType = "leave"
	ERROR   MessageType = "error"
	MESSAGE MessageType = "message"
)

/*
 // from john ["name"]
 InMessage{"join", "%TOKEN%"}
 // send to all
 OutMessage{"JOIN", "john"}
*/

const ServerName = "server"

type Message struct {
	UserName string `json:"name,omitempty"`
	Message  string `json:"message,omitempty"`
	Date     int64  `json:"date,omitempty"`
}

type JoinOutMessage struct {
	Username    string   `json:"name,omitempty"`
	JoinMessage string   `json:"joinMessage,omitempty"`
	Messages    []string `json:"messages,omitempty"`
	Date        int64    `json:"date,omitempty"`
	Success     bool     `json:"success,omitempty"`
}

type JoinInMessage struct {
	UserId   string `json:"name,omitempty"`
	RoomName string `json:"roomName,omitempty"`
}

func newJoinOutMessage(username string, messages []string, success bool) JoinOutMessage {
	return JoinOutMessage{username, username + " has joined", messages, jsTimeStamp(), success}
}

func newMessage(userName string, msg string) Message {
	return Message{userName, msg, jsTimeStamp()}
}

//func getMessageDataSentToServer(conn *websocket.Conn) (msg Message, err error) {
//	var raw []byte
//	_, raw, err = conn.ReadMessage()
//	if err != nil {
//		return msg, err
//	}
//	err = json.Unmarshal(raw, &msg)
//
//	if err != nil {
//		log.Println("Error - connHandler - Unmarshal - Incorrect data format:", string(raw), ":", err)
//		out, err := json.Marshal(Message{Typ: "error", Message: "Incorrect data format"})
//		if err != nil {
//			log.Println("Error - connHandler - MarshalError:", err)
//			return msg, err
//		}
//		if err = conn.WriteMessage(1, out); err != nil {
//			log.Println("Error - connHandler- WriteMessage Response:", err)
//			return msg, err
//		}
//		return msg, err
//	}
//
//	return msg, err
//}
//
//func connHandler(conn *websocket.Conn) (err error) {
//	var msg InMessage
//	msg, err = getMessageDataSentToServer(conn)
//
//	if err != nil {
//		if err.Error() != "websocket: close 1001 (going away)" {
//			log.Println("Error - connHandler - ReadMessage:", err)
//		}
//		return err
//	}
//	switch msg.Typ {
//	case JOIN:
//		err = processJoin(conn, msg)
//	case MESSAGE:
//		err = processMessage(conn, msg)
//	case LEAVE:
//		err = processLeave(conn)
//	}
//	return err
//}

//
//func processLeave(conn *websocket.Conn) (err error) {
//	if token, ok := CONNECTIONS[conn]; ok {
//		if user, ok := USERS[token]; ok {
//			msg := newOutMessage(LEAVE, user.name, user.name+" left the room")
//			err = broadcastMessage(user.roomId, msg)
//			if err != nil {
//				AddMsgToRoom(user.roomId, msg)
//			}
//			user.roomId = NullRoomId
//			USERS[token] = user
//		}
//	}
//	return err
//}
//
//func processJoin(conn *websocket.Conn, msg InMessage) (err error) {
//	token := msg.Message
//	if ok := AddConnection(conn, token); ok {
//		user := USERS[token]
//		msg := newOutMessage(JOIN, user.name, user.name+" joined the room")
//		err = broadcastMessage(user.roomId, msg)
//		if err != nil {
//			AddMsgToRoom(user.roomId, msg)
//		}
//	}
//
//	return err
//}
//
//func processMessage(conn *websocket.Conn, msg InMessage) (err error) {
//	if token, ok := CONNECTIONS[conn]; ok {
//		if currUser, ok := USERS[token]; ok {
//			msg := newOutMessage(MESSAGE, currUser.name, msg.Message)
//			err = broadcastMessage(currUser.roomId, msg)
//			if err != nil {
//				AddMsgToRoom(currUser.roomId, msg)
//			}
//
//		}
//	}
//	return err
//}
//
//func broadcastMessage(roomId string, msg OutMessage) (err error) {
//	if room, found := ROOMS[roomId]; found {
//		for _, userId := range room.userTokens {
//			err = sendMessage(userId, msg)
//			if err != nil {
//				return err
//			}
//		}
//	} else {
//		return RoomNotFoundError
//	}
//	return err
//}
//

func sendMessage(conn *websocket.Conn, marshalToJson interface{}) (err error) {
	if conn == nil {
		log.Println("Error - sendMessage - Programmer error, user conn not properly set")
		return UserConnInvalid
	}
	out, err := json.Marshal(marshalToJson)
	if err != nil {
		return err
	}
	if err = conn.WriteMessage(1, out); err != nil {
		log.Println("Error - sendMessage - WriteMessage Response:", err)
		return err
	}
	return nil
}

func closeConnAndRemoveData(conn *websocket.Conn) (err error) {
	if token, ok := CONNECTIONS[conn]; ok {
		delete(CONNECTIONS, conn)
		if _, ok := USERS[token]; ok {
			delete(USERS, token)
		}
	}
	return err
}

func forceLeave(conn *websocket.Conn) (err error) {
	defer closeConnAndRemoveData(conn)
	return err
	/*
		var out []byte
		var msg Message
		msg.Type = LEAVE
		msg.Src = CONNECTIONS[conn]

		defer closeConnAndRemoveData(conn)
		log.Println("Leave message received from", CONNECTIONS[conn])

		if isGame(conn) { // if sent by game, let controllers know it left
			controllersConnectedToGame := getAllControllerConnsFromGameConn(conn)
			msg.Src = CONNECTIONS[conn]
			out, err = json.Marshal(msg)
			for _, controllerConn := range controllersConnectedToGame {
				if err := controllerConn.WriteMessage(1, out); err != nil {
					log.Println("Error - forceLeave - WriteMessage:", err)
					return err
				}
			}
		} else { // if controller send, let game know it left.
			gameConn, found := getGameConnFromControllerConn(conn)
			msg.Src = CONNECTIONS[conn]
			if found {
				out, err = json.Marshal(msg)
				if err := gameConn.WriteMessage(1, out); err != nil {
					log.Println("Error - forceLeave - WriteMessage:", err)
					return err
				}
			}
		}
	*/
}

//type loadUserRequest struct {
//	UserId string `json:"userId"`
//}

// https://github.com/googollee/go-socket.io
//func SetupSocketIoServer() *socketio.Server {
//
//	//server.OnEvent(baseNamespace, "join", func(s socketio.Conn, msg Message) {
//	//	onJoin(server, s, msg)
//	//})
//	return server
//}

//func onJoin(server *socketio.Server, s socketio.Conn, msg Message) Message {
//	token := msg.Message
//	if ok := AddConnection(conn, token); ok {
//		user := USERS[token]
//		msg := newOutMessage(user.name, user.name+" joined the room")
//		err = broadcastMessage(user.roomId, msg)
//		if err != nil {
//			AddMsgToRoom(user.roomId, msg)
//		}
//	}
//
//	return Message{}
//}

// https://github.com/golang-jwt/jwt

//const userIdKey = "userId"
//
//func LoadUserData(r *http.Request) string {
//	return r.URL.Query().Get(userIdKey)
//}
//
//func LoadUser(conn *websocket.Conn, userId string) {
//	res := AddConnection(conn, userId)
//	if !res {
//		Log("failed to connect User")
//	}
//	LogConnections()
//}
//
//func SocketHandler(w http.ResponseWriter, r *http.Request) {
//	//Upgrade a HTTP Request to get a pointer to a Conn
//	//userId := LoadUserData(r)
//	//if userId != "" {
//	//	http.Error(w, "User did not provide userId", http.StatusInternalServerError)
//	//	return
//	//}
//
//	conn, err := upgrader.Upgrade(w, r, nil)
//
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	//LoadUser(conn, userId)
//
//	log.Println(conn.RemoteAddr(), "reached the server")
//
//	for {
//		if err := connHandler(conn); err != nil && err.Error() == "websocket: close 1001 (going away)" {
//			if user, isRegistered := CONNECTIONS[conn]; isRegistered {
//				log.Println("Connection closed for", user)
//				forceLeave(conn)
//			} else {
//				log.Println(conn.RemoteAddr(), "reached the server")
//			}
//			return
//		}
//	}
//}

func CreateSocketHandlers() *Router {
	websocketRouter := NewRouter()

	websocketRouter.OnConnection(func(c *websocket.Conn, r *http.Request) error {
		return onConnection(c, r)
	})

	websocketRouter.On("message", func(c *websocket.Conn, dt DataType) {
		onMessage(c, dt)
	})

	websocketRouter.On("leave", func(c *websocket.Conn, dt DataType) {
		onLeave(c, dt)
	})

	//websocketRouter.On("join", func(c *websocket.Conn, dt DataType) {
	//	onJoin(c, dt)
	//})
	return websocketRouter
}

func onMessage(c *websocket.Conn, dt DataType) {
	var msg Message
	err := json.Unmarshal(dt, &msg)
	if err != nil {
		// unmarshall data in failed
	}
	if user, ok := GetUser(c); ok {
		if room, ok := ROOMS[user.roomId]; ok {
			LogRoom(room.roomId)
			data := newMessage(user.name, msg.Message)
			outRaw, err := json.Marshal(data)
			if err != nil {
				log.Print(err)
			}
			broadCastMessageToRoom(room.roomId, msg, func(conn *websocket.Conn) {
				outMsg := NetWorkLayerMessage{Typ: string(MESSAGE), Data: outRaw}
				sendMessage(conn, outMsg)
			})
		}
	}
}

func onLeave(c *websocket.Conn, dt DataType) {
	//var leaveMessage nil
}

//func onJoin(c *websocket.Conn, dt DataType) {
//	var inMsg JoinInMessage
//	err := json.Unmarshal(dt, &inMsg)
//	if err != nil {
//		// unmarshall data in failed
//	}
//	if user, ok := GetUser(c); ok {
//
//	}
//
//}

func onConnection(c *websocket.Conn, r *http.Request) error {
	token, ok := r.URL.Query()["Authorization"]
	if !ok || len(token[0]) < 1 {
		log.Printf("invalid authorization token")
	}
	log.Printf(token[0])

	claims, authenticated := Authenticate(token[0])
	if authenticated {
		AddConnection(c, claims.UserToken)
		return nil
	}

	return errors.New("authentication failed")
}
