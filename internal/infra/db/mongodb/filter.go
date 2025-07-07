package mongodb

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OperationType defines the type for filter operations.
type OperationType string

const (
	OpEqual       OperationType = "=="
	OpNotEqual    OperationType = "!="
	OpGreaterThan OperationType = ">"
	OpLessThan    OperationType = "<"
	OpGTE         OperationType = ">="
	OpLTE         OperationType = "<="
	OpLike        OperationType = "LIKE"
	OpIn          OperationType = "IN"
	OpNotIn       OperationType = "NOT IN"
	OpRegex		  OperationType = "regex"
)

// FilterCondition represents a single filter condition for querying MongoDB.
type FilterCondition struct {
	Key       string
	Operation OperationType
	Value     interface{}
}

// Validate checks if the FilterCondition has valid fields.
func (fc *FilterCondition) Validate() error {
	if fc.Key == "" {
		return fmt.Errorf("filter condition key cannot be empty")
	}

	switch fc.Operation {
	case OpEqual, OpNotEqual, OpGreaterThan, OpLessThan, OpGTE, OpLTE, OpLike, OpIn, OpNotIn, OpRegex:
		// Supported operations
		return nil
	default:
		return fmt.Errorf("unsupported operation: %s", fc.Operation)
	}
}

// BuildBsonFilter converts a slice of FilterCondition into a bson.M filter.
func BuildBsonFilter(filters []FilterCondition) (bson.M, error) {
	if len(filters) == 0 {
		return bson.M{}, nil
	}

	var andConditions []bson.M

	for _, cond := range filters {
		// Validate each condition
		if err := cond.Validate(); err != nil {
			return nil, err
		}

		var condition bson.M

		switch cond.Operation {
		case OpEqual:
			condition = bson.M{cond.Key: cond.Value}
		case OpNotEqual:
			condition = bson.M{cond.Key: bson.M{"$ne": cond.Value}}
		case OpGreaterThan:
			condition = bson.M{cond.Key: bson.M{"$gt": cond.Value}}
		case OpLessThan:
			condition = bson.M{cond.Key: bson.M{"$lt": cond.Value}}
		case OpGTE:
			condition = bson.M{cond.Key: bson.M{"$gte": cond.Value}}
		case OpLTE:
			condition = bson.M{cond.Key: bson.M{"$lte": cond.Value}}
		case OpLike:
			valueStr, ok := cond.Value.(string)
			if !ok {
				return nil, fmt.Errorf("LIKE operation requires a string value")
			}
			// Convert SQL-like wildcard '*' to MongoDB regex '.*'
			regexPattern := "^" + EscapeRegex(valueStr) + "$"
			regexPattern = ReplaceWildcards(regexPattern)
			condition = bson.M{cond.Key: bson.M{"$regex": regexPattern, "$options": "i"}}
		case OpIn:
			values, ok := cond.Value.([]interface{})
			if !ok {
				return nil, fmt.Errorf("IN operation requires a slice of values")
			}
			condition = bson.M{cond.Key: bson.M{"$in": values}}
		case OpNotIn:
			values, ok := cond.Value.([]interface{})
			if !ok {
				return nil, fmt.Errorf("NOT IN operation requires a slice of values")
			}
			condition = bson.M{cond.Key: bson.M{"$nin": values}}
		case OpRegex:
			if regex, ok := cond.Value.(primitive.Regex); ok {
				condition = bson.M{cond.Key: bson.M{"$regex": regex.Pattern, "$options": regex.Options}}
			} else {
				return nil, fmt.Errorf("OpRegex requires primitive.Regex value")
			}
		default:
			return nil, fmt.Errorf("unsupported operation: %s, operationType: (%T)", cond.Operation, cond.Operation)
		}
		andConditions = append(andConditions, condition)
	}

	// Combine all conditions with $and
	finalFilter := bson.M{"$and": andConditions}

	return finalFilter, nil
}

// EscapeRegex escapes special regex characters in the input string except for '*'.
func EscapeRegex(input string) string {
	// Escaping regex special characters except for '*' which is used as wildcard
	specialChars := []string{".", "^", "$", "+", "?", "(", ")", "[", "]", "{", "}", "|", "\\"}
	escaped := input
	for _, char := range specialChars {
		escaped = strings.ReplaceAll(escaped, char, "\\"+char)
	}
	return escaped
}

// ReplaceWildcards replaces '*' with '.*' for regex matching.
func ReplaceWildcards(pattern string) string {
	return strings.ReplaceAll(pattern, "*", ".*")
}
