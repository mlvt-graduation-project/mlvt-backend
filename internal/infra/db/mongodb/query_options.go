package mongodb

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SortDirection represents ascending or descending order
type SortDirection int

const (
	SortAsc  SortDirection = 1
	SortDesc SortDirection = -1
)

// SortCondition defines a field and its direction
type SortCondition struct {
	Field     string
	Direction SortDirection
}

type QueryOptions struct {
	Filters []FilterCondition
	Sorts   []SortCondition
	Fields  []string
	Limit   *int
	Offset  *int 
}

func BuildQuery(qo QueryOptions) (bson.M, *options.FindOptions, error) {
	bsonFilter, err := BuildBsonFilter(qo.Filters)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build BSON filter: %v", err)
	}

	findOptions := options.Find()

	if qo.Limit != nil {
		findOptions.SetLimit(int64(*qo.Limit))
	}
	if qo.Offset != nil {
		findOptions.SetSkip(int64(*qo.Offset))
	}

	if qo.Sorts != nil {
		sortDoc := bson.D{}
		for _, s := range qo.Sorts {
			sortDoc = append(sortDoc, bson.E{Key: s.Field, Value: s.Direction})
		}
		findOptions.SetSort(sortDoc)
	}

	if qo.Fields != nil && len(qo.Fields) > 0 {
		projection := bson.M{}
		for _, field := range qo.Fields {
			projection[field] = 1
		}
		findOptions.SetProjection(projection)
	}

	return bsonFilter, findOptions, nil
}


func buidFindOptions(sorts []SortCondition, fields []string) *options.FindOptions {
	findOpts := options.Find()

	if len(sorts) > 0 {
		sortDoc := bson.D{}
		for _, s := range sorts {
			sortDoc = append(sortDoc, bson.E{
				Key:   s.Field,
				Value: s.Direction,
			})
		}
		findOpts.SetSort(sortDoc)
	}

	if len(fields) > 0 {
		projDoc := bson.M{}
		for _, field := range fields {
			projDoc[field] = 1
		}
		findOpts.SetProjection(projDoc)
	}

	return findOpts
}
