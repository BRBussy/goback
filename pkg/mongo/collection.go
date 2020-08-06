package mongo

import (
	"context"
	mongoBSON "go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"time"
)

// Collection is a mongodb client wrapper exposing mongo collection services
type Collection struct {
	driverCollection *mongoDriver.Collection
	timeout          time.Duration
}

func NewCollection(
	driverCollection *mongoDriver.Collection,
	timeout time.Duration,
) *Collection {
	return &Collection{
		driverCollection: driverCollection,
		timeout:          timeout,
	}
}

func (c *Collection) SetupIndex(model mongoDriver.IndexModel) error {
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	if _, err := c.driverCollection.Indexes().CreateOne(ctx, model); err != nil {
		return NewErrUnexpected(err)
	}
	return nil
}

func (c *Collection) SetupIndices(models []mongoDriver.IndexModel) error {
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	if _, err := c.driverCollection.Indexes().CreateMany(ctx, models); err != nil {
		return NewErrUnexpected(err)
	}
	return nil
}

func (c *Collection) CreateOne(document interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	_, err := c.driverCollection.InsertOne(ctx, document)
	return err
}

func (c *Collection) CreateMany(documents []interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	_, err := c.driverCollection.InsertMany(ctx, documents)
	return err
}

func (c *Collection) DeleteOne(filter Filter) error {
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	if _, err := c.driverCollection.DeleteOne(ctx, filter); err != nil {
		switch err {
		case mongoDriver.ErrNoDocuments:
			return NewErrUnexpected(err)
		default:
			return err
		}
	}
	return nil
}

func (c *Collection) FindOne(document interface{}, filter Filter) error {
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	if err := c.driverCollection.FindOne(ctx, filter).Decode(document); err != nil {
		switch err {
		case mongoDriver.ErrNoDocuments:
			return NewErrNotFound()
		default:
			return NewErrUnexpected(err)
		}
	}
	return nil
}

func (c *Collection) FindMany(documents interface{}, filter Filter, query Query) (int64, error) {
	// get options
	findOptions, err := query.ToMongoFindOptions()
	if err != nil {
		return 0, err
	}

	// perform find
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	cur, err := c.driverCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return 0, NewErrUnexpected(err)
	}

	// decode the results
	if err := cur.All(ctx, documents); err != nil {
		return 0, NewErrUnexpected(err)
	}

	// get document count
	count, err := c.driverCollection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, NewErrUnexpected(err)
	}

	return count, nil
}

func (c *Collection) UpdateOne(document interface{}, filter Filter) error {
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)
	if _, err := c.driverCollection.ReplaceOne(ctx, filter, document); err != nil {
		return NewErrUnexpected(err)
	}
	return nil
}

type aggregationCountHolder struct {
	Count int64 `bson:"count"`
}

func (c *Collection) Aggregate(pipeline mongoDriver.Pipeline, query Query, entities interface{}) (int64, error) {
	ctx, _ := context.WithTimeout(context.Background(), c.timeout)

	// perform aggregation and output count
	countCursor, err := c.driverCollection.Aggregate(
		ctx,
		append(
			pipeline,
			mongoBSON.D{{Key: "$count", Value: "count"}},
		),
	)
	if err != nil {
		return -1, NewErrUnexpected(err)
	}
	var countResults []aggregationCountHolder
	if err := countCursor.All(ctx, &countResults); err != nil {
		return -1, NewErrUnexpected(err)
	}
	var count int64
	if len(countResults) == 1 {
		count = countResults[0].Count
	} else if len(countResults) == 0 {
		count = 0
	} else {
		return -1, NewErrUnexpected(err)
	}

	// perform aggregation and output documents with query applied
	cursor, err := c.driverCollection.Aggregate(
		ctx,
		append(pipeline, query.ToPipelineStages()...),
	)
	if err != nil {
		return -1, NewErrUnexpected(err)
	}

	// decode the results
	if err := cursor.All(ctx, entities); err != nil {
		return -1, NewErrUnexpected(err)
	}

	return count, nil
}
