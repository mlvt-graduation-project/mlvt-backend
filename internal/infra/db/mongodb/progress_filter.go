package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"
)

func BuildProgressFilter(filters []FilterCondition) bson.M {
	return bson.M{}
}
