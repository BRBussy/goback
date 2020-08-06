package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// DatabaseConnection is a wrapper around the official go.mongodb.org golang driver Client
// implementing convenience functions to get handles to new Collection wrappers, close the
// connection and get DBStats
type DatabaseConnection struct {
	mongoClient *mongoDriver.Client
	database    *mongoDriver.Database
	timeout     time.Duration
}

// NewDatabaseConnection will connect and return a DatabaseConnection wrapper.
// Using a connectionString is prioritised over the combination of mongoDB hosts, username & password.
func NewDatabaseConnection(
	mongoDBHosts []string,
	mongoDBUsername,
	mongoDBPassword,
	connectionString,
	databaseName string,
	timeout time.Duration,
) (*DatabaseConnection, error) {

	if connectionString != "" {
		return NewFromConnectionString(connectionString, databaseName, timeout)
	} else if len(mongoDBHosts) != 0 {
		return NewFromHosts(mongoDBHosts, mongoDBUsername, mongoDBPassword, databaseName, timeout)
	}

	return nil, NewErrInvalidConfig([]string{"no hosts or connection string"})
}

// NewFromHosts is a will connect and return a DatabaseConnection wrapper.
func NewFromHosts(
	mongoDBHosts []string,
	mongoDBUsername,
	mongoDBPassword,
	databaseName string,
	timeout time.Duration,
) (*DatabaseConnection, error) {
	// create mongo client options
	mongoOptions := &options.ClientOptions{
		Hosts: mongoDBHosts,
	}

	// if a username is provided set auth on mongo client options
	if mongoDBUsername != "" {
		mongoOptions.SetAuth(options.Credential{
			Username:      mongoDBUsername,
			Password:      mongoDBPassword,
			AuthSource:    "admin",
			PasswordSet:   true,
			AuthMechanism: "SCRAM-SHA-1",
		})
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	mongoClient, err := mongoDriver.Connect(
		ctx,
		mongoOptions,
	)
	if err != nil {
		return nil, NewErrConnectionError(err)
	}

	// confirm that the client is connected
	ctx, cancelFn = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, NewErrPingError(err)
	}

	// connection successful populate and return database
	return &DatabaseConnection{
		mongoClient: mongoClient,
		database:    mongoClient.Database(databaseName),
		timeout:     timeout,
	}, nil
}

// NewFromConnectionString is a will connect and return a DatabaseConnection wrapper.
func NewFromConnectionString(connectionString string, databaseName string, timeout time.Duration) (*DatabaseConnection, error) {
	// create a new mongo client
	ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
	defer cancelFn()
	mongoClient, err := mongoDriver.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, NewErrConnectionError(err)
	}

	// confirm that the client is connected with a ping
	ctx, cancelFn = context.WithTimeout(context.Background(), timeout)
	defer cancelFn()
	if err := mongoClient.Ping(ctx, nil); err != nil {
		return nil, NewErrPingError(err)
	}

	return &DatabaseConnection{
		mongoClient: mongoClient,
		database:    mongoClient.Database(databaseName),
		timeout:     timeout,
	}, nil
}

// CloseConnection is a convenience function that can be used to close the wrapped go.mongodb.org golang driver Client.
func (d *DatabaseConnection) CloseConnection() error {
	if err := d.mongoClient.Disconnect(context.Background()); err != nil {
		return NewErrCloseConnectionError(err)
	}
	return nil
}

// CloseConnection is a convenience function that uses the wrapped go.mongodb.org golang driver Client
// to return a wrapped go.mongodb.org collection
func (d *DatabaseConnection) Collection(collectionName string) *Collection {
	return NewCollection(
		d.database.Collection(collectionName),
		d.timeout,
	)
}

// DBStats is used to decode the result of submitting a "dbStats" to the database server
type DBStats struct {
	// name of database
	DB string `json:"db" bson:"db"`

	// number of collections in the database
	Collections int32 `json:"collections" bson:"collections"`

	// number of documents across all collections
	Objects int32 `json:"objects" bson:"objects"`

	// average size of each object in bytes
	// NOT AFFECTED BY SCALE FACTOR
	AvgObjSize float64 `json:"avgObjSize" bson:"avgObjSize"`

	// total size of the uncompressed data held in database in bytes/scaleFactor
	DataSize float64 `json:"dataSize" bson:"dataSize"`

	// total amount of space allocated to collections in this database for document storage
	StorageSize float64 `json:"storageSize" bson:"storageSize"`

	// number of extents in the database across all collections
	NumExtents int32 `json:"numExtents" bson:"numExtents"`

	// number of indexes across all collections in the database
	Indexes int32 `json:"indexes" bson:"indexes"`

	// total size of all indexes created on this database
	IndexSize float64 `json:"indexSize" bson:"indexSize"`

	// scale used by the command
	ScaleFactor float64 `json:"scaleFactor" bson:"scaleFactor"`

	// total sized of all disk space in use on the filesystem where MongoDB stores data
	FSUsedSize float64 `json:"fsUsedSize" bson:"fsUsedSize"`

	// total size of all disk capacity on the filesystem where MongoDB stores data
	FSTotalSize float64 `json:"fsTotalSize" bson:"fsTotalSize"`

	Views int32   `json:"views" bson:"views"`
	Ok    float64 `json:"ok" bson:"ok"`
}

// CloseConnection is a convenience function that submits a "dbStats"
// query to the wrapped go.mongodb.org golang driver Client
func (d *DatabaseConnection) GetDBStats() (*DBStats, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), d.timeout)
	defer cancelFunc()

	result := d.database.RunCommand(
		ctx,
		bson.D{
			{
				Key:   "dbStats",
				Value: 1,
			},
			{
				Key:   "scale",
				Value: 1024,
			},
		},
	)
	dbStats := new(DBStats)
	if err := result.Decode(dbStats); err != nil {
		return nil, NewErrUnexpected(err)
	}

	return dbStats, nil
}
