package helper

import (
	"context"
	"fmt"
	"log"

	TodoModel "github.com/sumit-verma-trackier/todo/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Insert One helper
func InsertOneRecord(todo TodoModel.Todo, collection *mongo.Collection) TodoModel.Todo {

	response, err := collection.InsertOne(context.Background(), bson.M{
		"task":        todo.Task,
		"isCompleted": todo.IsCompleted,
	})

	if err != nil {
		log.Fatal(err)
	}

	todo.ID = response.InsertedID.(primitive.ObjectID)

	fmt.Println("Todo data inserted with id: ", response.InsertedID)
	return todo

}

// update one helper
func UpdateOneRecord(mongoId string, todo TodoModel.Todo, collection *mongo.Collection) (bool, error) {

	// change string id into mongo id
	id, _ := primitive.ObjectIDFromHex(mongoId)

	filter := bson.M{"_id": id}

	todo.ID = id

	updateData, _ := bson.Marshal(todo) // converting struct into raw BSON

	var updateDoc bson.M
	err := bson.Unmarshal(updateData, &updateDoc) // converting raw BSON Byte into BSON.M(like map)
	if err != nil {
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
func DeleteOneRecord(mongoId string, collection *mongo.Collection) {

	id, _ := primitive.ObjectIDFromHex(mongoId)

	filter := bson.M{"_id": id}

	res, _ := collection.DeleteOne(context.Background(), filter)

	fmt.Println("Record deleted successfully", res.DeletedCount)

}

// getAll helper
func GetAllRecord(collection *mongo.Collection) []primitive.M {

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
func GetOneRecord(mongoId string, collection *mongo.Collection) TodoModel.Todo {

	id, _ := primitive.ObjectIDFromHex(mongoId)

	filter := bson.M{"_id": id}

	var todoItem TodoModel.Todo

	res := collection.FindOne(context.Background(), filter)

	res.Decode(&todoItem)

	return todoItem

}

func GetTodoByName(name string, collection *mongo.Collection)TodoModel.Todo{

	filter := bson.M{"task":name}

	res:=collection.FindOne(context.Background(), filter)

	var todoItem TodoModel.Todo
	
	res.Decode(&todoItem)

	return todoItem



}
