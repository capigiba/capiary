package query

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BuildMongoQuery builds a BSON filter and FindOptions for MongoDB.
func BuildMongoQuery(opts QueryOptions) (bson.M, *options.FindOptions) {
	fmt.Println("opts: ", opts)
	filter := bson.M{}
	findOpts := options.Find()

	for _, f := range opts.Filters {
		op := map[OperationType]string{
			OpEqual:       "$eq",
			OpNotEqual:    "$ne",
			OpGreaterThan: "$gt",
			OpLessThan:    "$lt",
			OpGTE:         "$gte",
			OpLTE:         "$lte",
		}[f.Operator]

		if f.Operator == OpEqual {
			filter[f.Field] = f.Value
			continue
		}

		// ensure a sub‑doc exists
		if _, ok := filter[f.Field]; !ok {
			filter[f.Field] = bson.M{}
		}
		filter[f.Field].(bson.M)[op] = f.Value
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

	fmt.Println("filter: ", filter)
	fmt.Println("findOpts: ", findOpts)

	return filter, findOpts
}
