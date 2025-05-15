package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sumit-verma-trackier/todo/helper"
	TodoModel "github.com/sumit-verma-trackier/todo/model"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb://localhost:27017/"
const dbName = "todoDB"
const collectionName = "todo"

var collection *mongo.Collection

// connect with mongo

func init() {

	// client option
	clientOptions := options.Client().ApplyURI(connectionString)

	// connect to mongo

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mongo connected!")

	collection = client.Database(dbName).Collection(collectionName) // mongo reference is ready now

}

// ------ Controller ------

func AddTodoRecordController(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Content-Allow-Method", "POST")

	var todo TodoModel.Todo

	err := json.NewDecoder(r.Body).Decode(&todo)

	isExist := helper.GetTodoByName(todo.Task, collection)

	if len(isExist.Task) != 0 {
		json.NewEncoder(w).Encode(map[string]any{"status": false, "response": "Task with name already exist"})
		return
	}

	if err != nil {
		log.Fatal(err)
	}

	res := helper.InsertOneRecord(todo, collection)

	json.NewEncoder(w).Encode(map[string]any{"status": true, "response": res})

}

func GetAllTodoController(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	response := helper.GetAllRecord(collection)

	json.NewEncoder(w).Encode(response)

}

func GetTodoByIdController(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// path := r.URL.RawPath

	// fmt.Println(path, "---path")

	mongoId := mux.Vars(r)

	res := helper.GetOneRecord(mongoId["id"], collection)

	json.NewEncoder(w).Encode(res)
}

func UpdateTodoRecordController(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	// path := r.URL.RawPath

	mongoId := mux.Vars(r)

	var payload TodoModel.Todo

	json.NewDecoder(r.Body).Decode(&payload)

	helper.UpdateOneRecord(mongoId["id"], payload, collection)

	json.NewEncoder(w).Encode(map[string]any{"status": true, "response": "Todo Updated!"})

}

func DeleteTodoRecordController(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Content-Allow-Method", "DELETE")

	// path := r.URL.RawPath

	// fmt.Println(path, "path")

	mongoId := mux.Vars(r)

	helper.DeleteOneRecord(mongoId["id"], collection)

	json.NewEncoder(w).Encode(map[string]any{"status": true, "response": "Todo deleted!"})
}
