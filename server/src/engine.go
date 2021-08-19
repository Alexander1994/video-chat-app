package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

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

type DataType = json.RawMessage

type MsgHandle = func(*websocket.Conn, DataType)
type ErrorHandle = func(*websocket.Conn, *http.Request, error)
type ConnectionHandle = func(*websocket.Conn, *http.Request) error
type CheckOrigin = func(*http.Request) bool

type Router struct {
	upgrader websocket.Upgrader

	// handles
	msgHandlers map[string]MsgHandle
	connHandle  ConnectionHandle
	errHandle   ErrorHandle
}

type NetWorkLayerMessage struct {
	Typ  string   `json:"type,omitempty"`
	Data DataType `json:"data,omitempty"`
}

func NewRouter() (r *Router) {
	r = &Router{}
	r.msgHandlers = make(map[string]MsgHandle)
	r.errHandle = func(c *websocket.Conn, r *http.Request, e error) {
		fmt.Println(e)
	}
	r.connHandle = func(c *websocket.Conn, r *http.Request) error {
		log.Println(c.RemoteAddr(), "reached the server")
		return nil
	}
	r.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     checkOrigin,
	}

	return r
}

// setters
func (r *Router) Err(errHandle ErrorHandle) {
	r.errHandle = errHandle
}

func (r *Router) Upgrader(upgrader websocket.Upgrader) {
	r.upgrader = upgrader
}

func (r *Router) On(typ string, fn MsgHandle) {
	if _, ok := r.msgHandlers[typ]; ok {
		panic("Pre-existing route")
	} else {
		r.msgHandlers[typ] = fn
	}
}

func (r *Router) OnConnection(connHandle ConnectionHandle) {
	r.connHandle = connHandle
}

func (r *Router) RequestReceiver(w http.ResponseWriter, req *http.Request) {
	conn, err := r.upgrader.Upgrade(w, req, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = r.connHandle(conn, req)
	if err != nil {
		r.errHandle(conn, req, err)
		return
	}

	for {
		if err := r.connHandler(conn); err != nil && err.Error() == "websocket: close 1001 (going away)" {

			//if user, isRegistered := CONNECTIONS[conn]; isRegistered {
			//	log.Println("Connection closed for", user)
			//	forceLeave(conn)
			//} else {
			//}
			return
		} else {
			r.errHandle(conn, req, err)
		}
	}
}

func (r *Router) connHandler(conn *websocket.Conn) (err error) {
	var msg NetWorkLayerMessage
	msg, err = r.getMessageDataSentToServer(conn)

	if err != nil {
		if err.Error() != "websocket: close 1001 (going away)" {
			log.Println("Error - connHandler - ReadMessage:", err)
		}
		return err
	}
	log.Println(msg.Typ)
	if handle, ok := r.msgHandlers[msg.Typ]; ok {
		handle(conn, msg.Data)
	} else {
		// oh no
	}

	return err
}

func (r *Router) getMessageDataSentToServer(conn *websocket.Conn) (msg NetWorkLayerMessage, err error) {
	msg = NetWorkLayerMessage{}
	var raw []byte
	_, raw, err = conn.ReadMessage()
	if err != nil {
		return msg, err
	}
	err = json.Unmarshal(raw, &msg)
	if err != nil {
		log.Print(err)
		//log.Println("Error - connHandler - Unmarshal - Incorrect data format:", string(raw), ":", err)
		//out, err := json.Marshal(Message{Typ: "error", Message: "Incorrect data format"})
		//if err != nil {
		//	log.Println("Error - connHandler - MarshalError:", err)
		//	return msg, err
		//}
		//if err = conn.WriteMessage(1, out); err != nil {
		//	log.Println("Error - connHandler- WriteMessage Response:", err)
		//	return msg, err
		//}
		//return msg, err
	}

	return msg, err
}
