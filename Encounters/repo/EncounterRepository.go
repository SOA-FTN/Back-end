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

type EncounterRepository struct {
	cli    *mongo.Client
	logger *log.Logger
}

func NewEncounterRepository(ctx context.Context, logger *log.Logger) (*EncounterRepository, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(dburi))
	if err != nil {
		return nil, err
	}

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return &EncounterRepository{
		cli:    client,
		logger: logger,
	}, nil
}


func (pr *EncounterRepository) Ping() {
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

func (pr *EncounterRepository) Disconnect(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}
/*
func (tr *EncounterRepository) CreateEncounter(encounter *model.Encounter) error {
	dbResult := tr.DatabaseConnection.Create(encounter)
	if dbResult.Error != nil {
		return dbResult.Error
	}
	return nil
}
*/

func(enc *EncounterRepository) Insert(encounter *model.Encounter) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log.Println("=====================================")
	log.Println(encounter);
	encountersCollection := enc.getCollection()

	result,err := encountersCollection.InsertOne(ctx,encounter)
	if err != nil {
		enc.logger.Println(err)
		return err
	}
	enc.logger.Printf("Documents ID: %v\n", result.InsertedID)
	return nil

}

func (enc *EncounterRepository) getCollection() *mongo.Collection {
	encounterDatabase := enc.cli.Database("mongoDemo")
	encounterCollection := encounterDatabase.Collection("encounters")
	return encounterCollection
}

func (er *EncounterRepository) GetAllEncounters() (model.Encounters, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	encountersCollection := er.getCollection()

	var encounters model.Encounters
	encountersCursos , err := encountersCollection.Find(ctx,bson.M{})
	if err != nil {
		er.logger.Println(err)
		return nil, err
	}

	if err = encountersCursos.All(ctx, &encounters); err != nil {
		er.logger.Println(err)
		return nil, err
	}
	return encounters, nil

}
/*
func (repo EncounterRepository) UpdateEncounter(encounter *model.Encounter) (*model.Encounter, error) {
	dbResult := repo.DatabaseConnection.Model(&model.Encounter{}).Where("name=?", encounter.Name).Updates(encounter)
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	return encounter, nil
}
*/
func (er *EncounterRepository) GetEncounterByID(id string) (*model.Encounter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	encounterCollection := er.getCollection()

	var encounter model.Encounter
	objID, _ := primitive.ObjectIDFromHex(id)
	err := encounterCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&encounter)
	if err != nil {
		er.logger.Println(err)
		return nil, err
	}
	return &encounter, nil
}
