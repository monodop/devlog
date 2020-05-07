package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"nhooyr.io/websocket/wsjson"

	"nhooyr.io/websocket"
)

type wsMessage struct {
	Message string
}

func startWsListener(exitChannel chan bool, messageChannel chan string) {
	man := workerManager{
		nextWorkerId: 0,
		nextDataId:   100,
		workers:      make(map[int]chan string),
		data:         []string{
			// `{"id": 1, "app": "test", "logger": "MyApp.MyLogger", "message": "Hello World!"}`,
			// `{"id": 2, "app": "test", "logger": "MyApp.MyLogger", "message": "New message"}`,
			// `{"id": 3, "app": "other", "logger": "MyApp.Boi", "message": "New message"}`,
			// `{"id": 4, "app": "other", "logger": "MyApp.Boi", "message": "World Hello!"}`,
			// `{"id": 5, "app": "other", "logger": "MyApp.Boi", "message": "World Hello!"}`,
			// `{"id": 6, "app": "other", "logger": "MyApp.Boi", "message": "World Hello!"}`,
			// `{"id": 7, "app": "other", "logger": "MyApp.Boi", "message": "World Hello!"}`,
			// `{"id": 8, "app": "other", "logger": "MyApp.Boi", "message": "World Hello!"}`,
			// `{"id": 9, "app": "other", "logger": "MyApp.Boi", "message": "World Hello!"}`,
			// `{"id": 10, "app": "other", "logger": "MyApp.Boi", "message": "World Hello!"}`,
			// `{"id": 11, "app": "other", "logger": "MyApp.Boi", "message": "World Hello!"}`,
			// `{"id": 12, "app": "other", "logger": "MyApp.Boi", "message": "World Hello!"}`,
			// `{"id": 13, "app": "other", "logger": "MyApp.Boi", "message": "World Hello!"}`,
			// `{"id": 14, "app": "other", "logger": "MyApp.Boi", "message": "World Hello!"}`,
			// `{"id": 15, "app": "other", "logger": "MyApp.Boi", "message": "World Hello!"}`,
			// `{"id": 16, "superduperlong": "other", "logger": "MyApp.Boi", "message": "World Hello!"}`,
		},
	}
	go func() {
		for {
			msg := <-messageChannel
			man.AddMessage(msg)
		}
	}()
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		connection, err := websocket.Accept(writer, request, &websocket.AcceptOptions{
			OriginPatterns: []string{"*"},
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		defer connection.Close(websocket.StatusInternalError, "Unexpected error. Connection closing")

		handleWsConnection(connection, request, &man)
	})
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	err := http.ListenAndServe(":9091", nil)
	fmt.Println(err)
}

func handleWsConnection(connection *websocket.Conn, request *http.Request, man *workerManager) {

	ctx, cancel := context.WithCancel(request.Context())
	defer cancel()

	channel := make(chan string)
	id := man.AddWorker(channel)
	defer man.RemoveWorker(id)

	fmt.Printf("Opened WS connection %d to %s\n", id, request.RemoteAddr)
	defer fmt.Printf("Closed WS connection %d to %s\n", id, request.RemoteAddr)

	ctx = connection.CloseRead(ctx)

	for _, m := range man.data {
		err := wsjson.Write(ctx, connection, wsMessage{
			Message: m,
		})
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	for {
		select {
		case <-ctx.Done():
			connection.Close(websocket.StatusNormalClosure, "")
			return
		case msg := <-channel:
			err := wsjson.Write(ctx, connection, wsMessage{
				Message: msg,
			})
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func startWsTestConnection() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	connection, _, err := websocket.Dial(ctx, "ws://localhost:9091/ws", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer connection.Close(websocket.StatusInternalError, "Unexpected error, closing connection")

	for {
		message := wsMessage{}
		err = wsjson.Read(ctx, connection, &message)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("ws: %s\n", message.Message)
	}
}

type workerManager struct {
	sync.Mutex
	workers      map[int]chan string
	data         []string
	nextDataId   int
	nextWorkerId int
}

func (man *workerManager) AddWorker(worker chan string) int {
	man.Lock()
	defer man.Unlock()

	id := man.nextWorkerId
	man.nextWorkerId++

	man.workers[id] = worker
	return id
}

func (man *workerManager) RemoveWorker(id int) {
	man.Lock()
	defer man.Unlock()

	delete(man.workers, id)
}

func (man *workerManager) Iter(routine func(chan string)) {
	man.Lock()
	defer man.Unlock()

	for _, worker := range man.workers {
		routine(worker)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (man *workerManager) AddMessage(message string) {
	man.Lock()

	id := man.nextDataId
	man.nextDataId++

	var parsed map[string]interface{}
	json.Unmarshal([]byte(message), &parsed)
	parsed["id"] = id
	bytes, _ := json.Marshal(parsed)

	finalMessage := string(bytes)
	man.data = append(man.data, finalMessage)

	maxLength := 100
	numToRemove := max(0, len(man.data)-maxLength)
	man.data = man.data[numToRemove:]

	man.Unlock()
	man.Iter(func(w chan string) { w <- finalMessage })
}
