package repo

import (
	"context"
	"encounters/model"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type EncounterExecutionRepository struct {
	cli    *mongo.Client
	logger *log.Logger
}

func NewEncounterExecutionRepository(ctx context.Context, logger *log.Logger) (*EncounterExecutionRepository,error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &EncounterExecutionRepository{
		cli:    client,
		logger: logger,
	}, nil
}

func (pr *EncounterExecutionRepository) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := pr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		pr.logger.Println(err)
	}

	// Print available databases
	databases, err := pr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
	}
	fmt.Println(databases)
}

func (pr *EncounterExecutionRepository) Disconnect(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (enc *EncounterExecutionRepository) Insert(encounterExecution *model.EncounterExecution) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	encountersExecutionCollection := enc.getCollection()

	result,err := encountersExecutionCollection.InsertOne(ctx,&encounterExecution)
	if err != nil {
		enc.logger.Println(err)
		return err
	}
	enc.logger.Printf("Documents ID: %v\n", result.InsertedID)
	return nil
}

func (enc *EncounterExecutionRepository) getCollection() *mongo.Collection {
	encounterExecutionDatabase := enc.cli.Database("mongoDemo")
	encounterExecutionCollection := encounterExecutionDatabase.Collection("encountersExecution")
	return encounterExecutionCollection
}

func (er *EncounterExecutionRepository) GetAllEncounterExecutions() (model.EncounterExecutions, error) {
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel();

	executionsCollection := er.getCollection()

	var encounterExecutions model.EncounterExecutions

	executionsCursor,err := executionsCollection.Find(ctx,bson.M{})
	if err != nil {
		er.logger.Println(err)
		return nil, err
	}

	for executionsCursor.Next(ctx) {
		var exec model.EncounterExecution
		if err := executionsCursor.Decode(&exec); err != nil {
			er.logger.Println(err)
			return nil, err
		}
		
		// // Konvertuj CompletionTime u željeni format stringa
		// t, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", exec.CompletionTime)
        // if err != nil {
        //     er.logger.Println(err)
        //     return nil, err
        // }
        // exec.CompletionTime = t.Format("2006-01-02T15:04:05")

        encounterExecutions = append(encounterExecutions, &exec)
    }
	if err := executionsCursor.Err(); err != nil {
		er.logger.Println(err)
		return nil, err
	}

	return encounterExecutions, nil

}

func (er *EncounterExecutionRepository) GetByUserIDAndNotCompleted(userID int) (*model.EncounterExecution, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    executionsCollection := er.getCollection()

    var encounterExecution *model.EncounterExecution

    err := executionsCollection.FindOne(ctx, bson.M{"userId": userID, "iscompleted": false}).Decode(&encounterExecution)
    if err != nil {
        er.logger.Println(err)
        return nil, err
    }

    // Ako ne pronađemo nijedan susret koji nije završen, ne vraćamo grešku, već samo nil za susret
    if encounterExecution == nil {
        return nil, nil
    }

    // Konvertuj CompletionTime u željeni format stringa
    // t, err := time.Parse("2006-01-02 15:04:05 +0000 UTC", encounterExecution.CompletionTime)
    // if err != nil {
    //     er.logger.Println(err)
    //     return nil, err
    // }
    // encounterExecution.CompletionTime = t.Format("2006-01-02T15:04:05")

    return encounterExecution, nil
}

func (eer *EncounterExecutionRepository) Update(id string,execution *model.EncounterExecution) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	executionCollection := eer.getCollection()

	objID, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        eer.logger.Println(err)
        return err
    }
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"completionTime":    execution.CompletionTime,
		"iscompleted": execution.IsCompleted,
	}}
	result, err := executionCollection.UpdateOne(ctx, filter, update)
	eer.logger.Printf("Documents matched: %v\n", result.MatchedCount)
	eer.logger.Printf("Documents updated: %v\n", result.ModifiedCount)

	if err != nil {
		eer.logger.Println(err)
		return err
	}
	return nil
}

