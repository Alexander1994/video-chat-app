package main

import (
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/gorilla/websocket"
)

//Define upgrade policy.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

//Ugrade policty from http request to websocket, to be defined
func checkOrigin(r *http.Request) bool {
	//For example: Check in a blacklist if the address is present
	//if blacklist_check(r.RemoteAddr) { return false }
	// check if origin header is correct url
	// if isBlackListed(r) {
	// 	return false
	// }
	//origin := r.Header.Get("Origin")
	return true
	//if !ENV.Public && origin == "null" {
	//	return true
	//}
	//
	//if origin != "" {
	//	if origin == "http://alexmccallum.me" || origin == "file://" {
	//		return true
	//	}
	//}
	//return false
}

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

func newOutMessage(userName string, msg string) Message {
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
//func sendMessage(userToken Token, msg OutMessage) (err error) {
//	if userOut, ok := USERS[userToken]; ok {
//		var out []byte
//		out, err = json.Marshal(msg)
//		if err != nil { // programming error, can't marshall type
//			log.Println("Error - connHandler - MarshalError:", err)
//			return err
//		}
//		if userOut.conn == nil {
//			log.Println("Error - connHandler - Programmer error, user conn not properly set")
//			return UserConnInvalid
//		}
//		if err = userOut.conn.WriteMessage(1, out); err != nil {
//			log.Println("Error - connHandler - WriteMessage Response:", err)
//			return err
//		}
//	} else {
//		return UserNotFoundError
//	}
//	return err
//}

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
const baseNamespace = "/"

// https://github.com/googollee/go-socket.io
func SetupSocketIoServer() *socketio.Server {
	server := socketio.NewServer(nil)
	// join
	server.OnConnect(baseNamespace, func(c socketio.Conn) error {
		return onConnect(server, c)
	})
	server.OnDisconnect(baseNamespace, func(c socketio.Conn, s string) {
		onDisconnect(server, c, s)
	})

	server.OnEvent(baseNamespace, "message", func(s socketio.Conn, msg Message) {
		onMessage(server, s, msg)
	})
	server.OnEvent(baseNamespace, "leave", func(s socketio.Conn, msg Message) {
		onLeave(server, s, msg)
	})
	server.OnEvent(baseNamespace, "connect_error", func(s socketio.Conn, msg string) {
		Log(msg)
	})
	//server.OnEvent(baseNamespace, "join", func(s socketio.Conn, msg Message) {
	//	onJoin(server, s, msg)
	//})
	return server
}

func onMessage(server *socketio.Server, s socketio.Conn, msg Message) Message {
	Log(msg.Message)
	s.SetContext(msg)
	return msg
}

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

func onLeave(server *socketio.Server, s socketio.Conn, msg Message) Message {
	return Message{}
}

func onDisconnect(server *socketio.Server, c socketio.Conn, reason string) {

}

// https://github.com/golang-jwt/jwt
func onConnect(server *socketio.Server, s socketio.Conn) error {
	Log("WHAT")
	_, authenticated := Authenticate(s.RemoteHeader())
	if !authenticated {
		s.LeaveAll()
		s.Close()
	}
	return nil
}

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
//func SocketHandler(c *gin.Context) {
//	socketHandlerWrapped(c.Writer, c.Request)
//}
//
////Catches HTTP Requests, upgrade them if needed and let connHandler managing the connection
//func socketHandlerWrapped(w http.ResponseWriter, r *http.Request) {
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
