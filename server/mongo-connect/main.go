package main

import (
   "context"
   "fmt"
   "log"
   "net/http"
   "encoding/json"

   "go.mongodb.org/mongo-driver/mongo"
   "go.mongodb.org/mongo-driver/mongo/options"
   "go.mongodb.org/mongo-driver/bson"
   "github.com/gorilla/mux"
)

// Define your MongoDB connection string
const connectionString = "mongodb://localhost:27017"

// Define your MongoDB database and collection names
const dbName = "test"
const collectionName = "item-price"

// Define a MongoDB client
var client *mongo.Client

// Define a struct to represent your data model
type Item struct {
   ID    string `json:"id" bson:"_id,omitempty"`
   Name  string `json:"name"`
   Price float64 `json:"price"`
}

// Initialize the MongoDB client
func init() {
   var err error
   client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(connectionString))
   if err != nil {
      log.Fatal(err)
   }
}

// Create a new item
func createItem(w http.ResponseWriter, r *http.Request) {
   var newItem Item
   json.NewDecoder(r.Body).Decode(&newItem)

   collection := client.Database(dbName).Collection(collectionName)
   _, err := collection.InsertOne(context.Background(), newItem)
   if err != nil {
      log.Fatal(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
   }

   w.WriteHeader(http.StatusCreated)
}

// Get all items
func getAllItems(w http.ResponseWriter, r *http.Request) {
   collection := client.Database(dbName).Collection(collectionName)
   cursor, err := collection.Find(context.Background(), bson.D{})
   if err != nil {
      log.Fatal(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
   }
   defer cursor.Close(context.Background())

   var items []Item
   if err := cursor.All(context.Background(), &items); err != nil {
      log.Fatal(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
   }

   w.Header().Set("Content-Type", "application/json")
   json.NewEncoder(w).Encode(items)
}

// Update an item by ID
func updateItem(w http.ResponseWriter, r *http.Request) {
   params := mux.Vars(r)
   id := params["id"]

   var updatedItem Item
   json.NewDecoder(r.Body).Decode(&updatedItem)

   collection := client.Database(dbName).Collection(collectionName)
   filter := bson.D{{"_id", id}}
   update := bson.D{{"$set", updatedItem}}

   _, err := collection.UpdateOne(context.Background(), filter, update)
   if err != nil {
      log.Fatal(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
   }

   w.WriteHeader(http.StatusOK)
}

// Delete an item by ID
func deleteItem(w http.ResponseWriter, r *http.Request) {
   params := mux.Vars(r)
   id := params["id"]

   collection := client.Database(dbName).Collection(collectionName)
   filter := bson.D{{"_id", id}}

   _, err := collection.DeleteOne(context.Background(), filter)
   if err != nil {
      log.Fatal(err)
      w.WriteHeader(http.StatusInternalServerError)
      return
   }

   w.WriteHeader(http.StatusNoContent)
}

func main() {
   // Initialize router
   router := mux.NewRouter()

   // Define routes
   router.HandleFunc("/items", createItem).Methods("POST")
   router.HandleFunc("/items", getAllItems).Methods("GET")
   router.HandleFunc("/items/{id}", updateItem).Methods("PUT")
   router.HandleFunc("/items/{id}", deleteItem).Methods("DELETE")

   // Start the server
   fmt.Println("Server is running on :8080")
   log.Fatal(http.ListenAndServe(":8080", router))
}
