package query

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BuildMongoQuery builds a BSON filter and FindOptions for MongoDB.
func BuildMongoQuery(opts QueryOptions) (bson.M, *options.FindOptions) {
	filter := bson.M{}
	findOpts := options.Find()

	for _, fil := range opts.Filters {
		switch fil.Operator {
		case OpEqual:
			filter[fil.Field] = fil.Value
		case OpNotEqual:
			filter[fil.Field] = bson.M{"$ne": fil.Value}
		case OpGreaterThan:
			filter[fil.Field] = bson.M{"$gt": fil.Value}
		case OpLessThan:
			filter[fil.Field] = bson.M{"$lt": fil.Value}
		case OpGTE:
			filter[fil.Field] = bson.M{"$gte": fil.Value}
		case OpLTE:
			filter[fil.Field] = bson.M{"$lte": fil.Value}
		default:
			// Fallback or handle custom operators
			filter[fil.Field] = fil.Value
		}
	}

	// Sort
	if len(opts.Sorts) > 0 {
		sortDoc := bson.D{}
		for _, s := range opts.Sorts {
			dir := 1
			if s.Desc {
				dir = -1
			}
			sortDoc = append(sortDoc, bson.E{Key: s.Field, Value: dir})
		}
		findOpts.SetSort(sortDoc)
	}

	// Projection (fields)
	if len(opts.Fields) > 0 {
		projection := bson.M{}
		for _, f := range opts.Fields {
			f = strings.TrimSpace(f)
			projection[f] = 1
		}
		findOpts.SetProjection(projection)
	}

	if opts.Skip > 0 {
		findOpts.SetSkip(opts.Skip)
	}
	if opts.Limit > 0 {
		findOpts.SetLimit(opts.Limit)
	}

	return filter, findOpts
}
