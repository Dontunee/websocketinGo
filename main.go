package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//Define our Upgrader
//In order to create a WebSocket endpoint,
//we effectively need to upgrade an incoming connection from a standard HTTP endpoint to a long-lasting WebSocket connection
// This will hold information such as the Read and Write buffer size for our WebSocket connection
var upgrader = websocket.Upgrader{}

type MyStruct struct {
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type Response struct {
	Status  bool        `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

var listOfData []MyStruct
var index = 0

func main() {
	route := mux.NewRouter()
	route.HandleFunc("/", IndexHandler)
	route.HandleFunc("/v1/ws", WsHandler)
	route.HandleFunc("/v1/person/add", AddPersonHandler).Methods("POST")

	srv := &http.Server{
		Addr:              ":3000",
		Handler:           route,
		TLSConfig:         nil,
		ReadTimeout:       90 * time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      90 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server - %s\n", err.Error())
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	//Get working directory in environment
	dir, _ := os.Getwd()
	//search for /index.html
	dbs := filepath.Join(dir, "/index.html")
	log.Println(dbs)
	http.ServeFile(w, r, filepath.Join(dir, "/index.html"))
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	//This will determine whether an incoming request from a different domain is allowed to connect,
	////and if it isn’t they’ll be hit with a CORS error.
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	//Upgrading our Connection
	var conn, _ = upgrader.Upgrade(w, r, nil)

	go readMessage(conn)

	go writeNotificationMessage(conn, &index)
}

func AddPersonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data, response = MyStruct{}, Response{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Status = false
		response.Data = nil
		response.Message = fmt.Sprintln("An error occurred: ", err.Error())
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response.Status = true
	response.Data = data
	response.Message = fmt.Sprint("successfully added person")
	_ = json.NewEncoder(w).Encode(response)
	listOfData = append(listOfData, data)
	return

}

func readMessage(conn *websocket.Conn) {
	//continually listen for any incoming messages sent through that WebSocket connection.
	//We’ll call this reader() for now and it will take in a pointer to the WebSocket connection that we received from our call to upgrader.Upgrade:
	for {
		var myData MyStruct
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error on closing: %s\n", err.Error())
			conn.Close()
			return
		}
		err = json.Unmarshal(msg, &myData)
		if err != nil {
			fmt.Printf("error unmarshalling data: %s\n", err.Error())
			return
		}
		listOfData = append(listOfData, myData)
		fmt.Printf("list of data %+v\n", listOfData)
	}
}

func writeMessage(conn *websocket.Conn) {
	ch := time.Tick(5 * time.Second)
	for range ch {
		conn.WriteJSON(MyStruct{
			UserName:  "",
			FirstName: "",
			LastName:  "",
		})
	}
}

func writeNotificationMessage(conn *websocket.Conn, index *int) {
	ch := time.Tick(5 * time.Second)
	for range ch {
		if (len(listOfData)) > *index {
			for _, v := range listOfData[*index:] {
				//Writing Back to the Client
				conn.WriteJSON(MyStruct{
					UserName:  v.UserName,
					FirstName: v.FirstName,
					LastName:  v.LastName,
				})
				*index += 1
				fmt.Printf("%+v\n", v)
				fmt.Printf("all data: %+v\n", listOfData)
			}
		}
	}
}
