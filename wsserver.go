package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/monodop/devlog/env"
	"github.com/monodop/devlog/log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type wsMessage struct {
	Message string
}

func startWsListener(exitChannel chan bool, messageChannel chan string) {
	address := env.HttpListenAddress()

	man := workerManager{
		nextWorkerId: 0,
		nextDataId:   0,
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
			log.Exception(err)
			return
		}
		defer connection.Close(websocket.StatusInternalError, "Unexpected error. Connection closing")

		handleWsConnection(connection, request, &man)
	})
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	log.Info("HTTP/WebSocket server now listening on %s", address)
	err := http.ListenAndServe(address, nil)
	log.Exception(err)
}

func handleWsConnection(connection *websocket.Conn, request *http.Request, man *workerManager) {

	ctx, cancel := context.WithCancel(request.Context())
	defer cancel()

	channel := make(chan string)
	id := man.AddWorker(channel)
	defer man.RemoveWorker(id)

	log.Info("Opened WS connection %d to %s", id, request.RemoteAddr)
	defer log.Info("Closed WS connection %d to %s", id, request.RemoteAddr)

	ctx = connection.CloseRead(ctx)

	for _, m := range man.data {
		err := wsjson.Write(ctx, connection, wsMessage{
			Message: m,
		})
		if err != nil {
			log.Exception(err)
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
				log.Exception(err)
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
		log.Exception(err)
		return
	}
	defer connection.Close(websocket.StatusInternalError, "Unexpected error, closing connection")

	for {
		message := wsMessage{}
		err = wsjson.Read(ctx, connection, &message)
		if err != nil {
			log.Exception(err)
			return
		}

		log.Info("ws: %s", message.Message)
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
	locked := true
	unlock := func() {
		if locked {
			man.Unlock()
		}
		locked = false
	}
	defer unlock()

	// Generate new unique id for message
	id := man.nextDataId
	man.nextDataId++

	// Parse message
	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(message), &parsed)
	if err != nil {
		log.Error("Error serializing message: %s", err)
		return
	}

	// Tag message with unique id
	parsed["_id"] = id

	// Re-serialize message
	bytes, err := json.Marshal(parsed)
	if err != nil {
		log.Error("Error serializing message: %s", err)
		return
	}

	// Add message to short-term memory
	finalMessage := string(bytes)
	man.data = append(man.data, finalMessage)

	// Clear short-term memory that's older than InternalMessageBufferSize messages
	maxLength := env.InternalMessageBufferSize()
	numToRemove := max(0, len(man.data)-maxLength)
	man.data = man.data[numToRemove:]

	// Send message to all listeners
	unlock()
	man.Iter(func(w chan string) { w <- finalMessage })
}
