package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	TodoModel "github.com/sumit-verma-trackier/todo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Helper method (can create this into another folder called helper)

// Insert One helper
func insertOneRecord(todo TodoModel.Todo) TodoModel.Todo {

	response, err := collection.InsertOne(context.Background(), todo)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Todo data inserted with id: ", response.InsertedID)
	return todo

}

// update one helper
func updateOneRecord(mongoId string, todo TodoModel.Todo) (bool, error) {

	// change string id into mongo id
	id, _ := primitive.ObjectIDFromHex(mongoId)

	filter := bson.M{"_id": id}

	todo.ID = id

	updateData, _ := bson.Marshal(todo)

	var updateDoc bson.M
	if err := bson.Unmarshal(updateData, &updateDoc); err != nil {
		log.Println("Failed to unmarshal to bson.M:", err)
		return false, err
	}

	update := bson.M{"$set": updateDoc}

	res, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Println("Failed to update record:", err)
		return false, nil
	}

	return res.ModifiedCount > 0, nil

}

// delete one helper
func deleteOneRecord(mongoId string) {

	id, _ := primitive.ObjectIDFromHex(mongoId)

	filter := bson.M{"_id": id}

	res, _ := collection.DeleteOne(context.Background(), filter)

	fmt.Println("Record deleted successfully", res.DeletedCount)

}

// getAll helper
func getAllRecord() []primitive.M {

	res, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		log.Fatal(err)
	}

	var todoList []primitive.M

	for res.Next(context.Background()) {
		var todo bson.M

		err := res.Decode(&todo)

		if err != nil {
			log.Fatal(err)
		}

		todoList = append(todoList, todo)

	}

	defer res.Close(context.Background())

	return todoList

}

// func Get one helper
func getOneRecord(mongoId string) TodoModel.Todo {

	id, _ := primitive.ObjectIDFromHex(mongoId)

	filter := bson.M{"_id": id}

	var todoItem TodoModel.Todo

	res := collection.FindOne(context.Background(), filter)

	res.Decode(&todoItem)

	return todoItem

}






// Controller (only this should be here in file)

func GetAllTodoController(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	response := getAllRecord()

	json.NewEncoder(w).Encode(response)

}

func GetTodoByIdController(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// path := r.URL.RawPath

	// fmt.Println(path, "---path")

	mongoId := mux.Vars(r)

	res := getOneRecord(mongoId["id"])

	json.NewEncoder(w).Encode(res)
}

func AddTodoRecordController(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Content-Allow-Method", "POST")

	var todo TodoModel.Todo

	err := json.NewDecoder(r.Body).Decode(&todo)

	if err != nil {
		log.Fatal(err)
	}

	insertOneRecord(todo)

	json.NewEncoder(w).Encode(map[string]any{"status": true, "response": todo})

}

func UpdateTodoRecordController(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	// path := r.URL.RawPath

	mongoId := mux.Vars(r)

	var payload TodoModel.Todo

	json.NewDecoder(r.Body).Decode(&payload)

	updateOneRecord(mongoId["id"], payload)

	json.NewEncoder(w).Encode(map[string]any{"status": true, "response": "Todo Updated!"})

}

func DeleteTodoRecordController(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Content-Allow-Method", "DELETE")

	// path := r.URL.RawPath

	// fmt.Println(path, "path")

	mongoId := mux.Vars(r)

	deleteOneRecord(mongoId["id"])

	json.NewEncoder(w).Encode(map[string]any{"status": true, "response": "Todo deleted!"})
}
