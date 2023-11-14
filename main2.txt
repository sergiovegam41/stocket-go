package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Adjust the origin check based on your needs
	},
}

// DBNames holds the collection names
var DBNames = struct {
	Services       string
	Professions    string
	NotifyMeOrders string
}{
	Services:       "services",
	Professions:    "professions",
	NotifyMeOrders: "notifyMeOrders",
}

// watchChanges watches the MongoDB collection for changes
func watchChanges(client *mongo.Client, ctx context.Context, sendUpdate func(slugName string)) {
	collection := client.Database("medify").Collection("mensajes")
	changeStream, err := collection.Watch(ctx, mongo.Pipeline{})
	if err != nil {
		log.Fatal(err)
	}
	defer changeStream.Close(ctx)

	for changeStream.Next(ctx) {
		var change bson.M
		if err := changeStream.Decode(&change); err != nil {
			log.Fatal(err)
		}

		// Handle the change document
		operationType, ok := change["operationType"].(string)
		if !ok {
			continue
		}

		if operationType == "insert" || operationType == "update" {
			documentKey := change["documentKey"].(bson.M)
			id := documentKey["_id"]

			var service bson.M
			if err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&service); err != nil {
				log.Println(err)
				continue
			}

			professionID := service["profession_id"]
			professionsCollection := client.Database("test-dservice-backend").Collection("Professions")
			var profession bson.M
			if err := professionsCollection.FindOne(ctx, bson.M{"_id": professionID}).Decode(&profession); err != nil {
				log.Println(err)
				continue
			}

			slugName, ok := profession["slug_name"].(string)
			if !ok {
				continue
			}

			sendUpdate(slugName)
		}
	}
	if err := changeStream.Err(); err != nil {
		log.Fatal(err)
	}
}

// handleConnections handles incoming websocket connections
func handleConnections(w http.ResponseWriter, r *http.Request, client *mongo.Client) {
	// ctx := context.Background()

	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	// Send current data to the client
	// getCurrentData() should be implemented
	// currentData, err := getCurrentData()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// if err := ws.WriteJSON(currentData); err != nil {
	// 	log.Println(err)
	// 	return
	// }

	for {
		var msg Message // Message is a struct that you'll need to define according to your application's protocol
		// Wait for a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		// Handle new messages
		switch msg.Type {
		case "client:setNotifyMeOrders":
			// Handle the message
		case "client:getData":
			// Handle the message
		}

		// Send a message to the client
		response := Message{} // Create a response message according to your needs
		if err := ws.WriteJSON(response); err != nil {
			log.Println(err)
			return
		}
	}
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

// Message defines the structure for our messages
type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func main() {


    // Define la carpeta pública que quieres exponer, digamos que se llama "public"
    fs := http.FileServer(http.Dir("./public"))

    // Servir los archivos estáticos de la carpeta 'public' en la ruta '/'
    http.Handle("/", fs)

    // El resto 
	
	
	
	de tu código para manejar las conexiones WebSocket y arrancar el servidor...
    log.Println("http server started on :8000")

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb+srv://stock-manager:r4mEHcjNNw3z9u3K@cluster0.0eyr1.mongodb.net/stock-manager?retryWrites=true&w=majority")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	// Start listening for changes in the background
	go watchChanges(client, context.Background(), func(slugName string) {
		// This function will be called when there's a change in the services collection
		// You need to implement the logic to send updates to the connected clients
		for client := range clients {
			if err := client.WriteJSON(Message{Type: "server:refresh", Data: slugName}); err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	})

	// Start handling websocket connections
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleConnections(w, r, client)
	})

	// Start the server on localhost port 8000 and log any errors
	log.Println("http server started on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
