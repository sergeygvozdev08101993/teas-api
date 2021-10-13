package pkg

import (
	"context"
	"github/gvozdev08101993/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getTeaByName(teaName string) (Tea, error) {

	var (
		tmpData bson.M
		data Tea
		err error
	)

	teasDB := db.ClientDB.Database("teas-api")
	teasCollection := teasDB.Collection("teas")

	err = teasCollection.FindOne(context.TODO(), bson.D{primitive.E{Key: "name", Value: teaName}}).Decode(&tmpData)
	if err != nil {
		return data, err
	}
	if err == mongo.ErrNoDocuments {
		return data, err
	}

	data = Tea{
		Name: tmpData["name"].(string),
		Category: tmpData["category"].(string),
	}

	return data, nil
}

func createTea(tea Tea) error {

	teasDB := db.ClientDB.Database("teas-api")
	teasCollection := teasDB.Collection("teas")

	_, err := teasCollection.InsertOne(context.TODO(), bson.D{
		primitive.E{Key: "name", Value: tea.Name},
		primitive.E{Key: "category", Value: tea.Category}})
	if err != nil {
		return err
	}

	return nil
}

func deleteTeaByName(teaName string) error {

	teasDB := db.ClientDB.Database("teas-api")
	teasCollection := teasDB.Collection("teas")

	_, err := teasCollection.DeleteOne(context.TODO(), bson.D{primitive.E{
		Key: "name", Value: teaName,
	}})
	if err != nil {
		return err
	}

	return nil
}

func updateTea(tea Tea, name string) error {

	teasDB := db.ClientDB.Database("teas-api")
	teasCollection := teasDB.Collection("teas")

	_, err := teasCollection.UpdateOne(
		context.TODO(),
		bson.M{"name": name},
		bson.D{
			{Key: "$set", Value: bson.D{primitive.E{Key: "name", Value: tea.Name}}},
			{Key: "$set", Value: bson.D{primitive.E{Key: "category", Value: tea.Category}}},
		})
	if err != nil {
		return err
	}

	return nil
}

func getAllTeas() ([]Tea, error) {

	var (
		tmpTeas []bson.D
		teas []Tea
		err error
	)

	teasDB := db.ClientDB.Database("teas-api")
	teasCollection := teasDB.Collection("teas")

	ctx := context.TODO()
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "name", Value: 1}})

	cursor, err := teasCollection.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &tmpTeas); err != nil {
		return nil, err
	}

	for _, tea := range tmpTeas {
		name := tea[1].Value.(string)
		category := tea[2].Value.(string)

		teas = append(teas, Tea{Name: name, Category: category})
	}

	return teas, nil
}
