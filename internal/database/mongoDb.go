package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// mongoDatabase holds a collection of activities
type mongoDatabase struct {
	mongoClient *mongo.Client
}

func NewMongoDatabase(dbConnectionString string) *mongoDatabase {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(dbConnectionString).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	db := &mongoDatabase{mongoClient: client}

	return db
}

func (db *mongoDatabase) Disconnect() {
	if err := db.mongoClient.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func (db *mongoDatabase) AddToken(token UserToken) error {
	// Get a handle for your collection
	collection := db.mongoClient.Database("frostPointsDev").Collection("users")

	// Insert a single document
	insertResult, err := collection.InsertOne(context.TODO(), token)
	if err != nil {
		log.Println("Error adding userToken to db: ", err)
		return err
	}

	log.Println("Added userToken to db: ", insertResult.InsertedID)
	return err
}

func (db *mongoDatabase) FindTokenById(athleteId int64) *UserToken {
	collection := db.mongoClient.Database("frostPointsDev").Collection("users")

	var result UserToken
	err := collection.FindOne(context.TODO(), bson.D{{"athlete_id", athleteId}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("No document was found with the specified athlete_id")
		} else {
			log.Printf("Error finding document: %v\n", err)
		}
		return nil
	}

	return &result
}

func (db *mongoDatabase) DeleteToken(athleteId int64) error {
	collection := db.mongoClient.Database("frostPointsDev").Collection("users")

	result, err := collection.DeleteOne(context.TODO(), bson.D{{"athlete_id", athleteId}})
	if err != nil {
		log.Printf("Error deleting id %v\n", err)
		return err
	}

	fmt.Printf("Deleted %v documents in the usertokens collection\n", result.DeletedCount)
	return nil
}
