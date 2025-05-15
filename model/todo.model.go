package TodoModel

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID          primitive.ObjectID `json:"taskId,omitempty" bson:"_id,omitempty"`
	Task        string             `json:"task" bson:"task"`
	IsCompleted bool               `json:"isCompleted" bson:"isCompleted"`
}
