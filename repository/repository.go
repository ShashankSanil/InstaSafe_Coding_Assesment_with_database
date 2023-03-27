package repository

import (
	"context"
	"fmt"

	"github.com/go-chassis/openlog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	DBClient *mongo.Client
	DBName   string
}

var CollectionName = "Transaction_Details"

var UCollection = "User_Details"

func (tr *Repository) FindUserByEmail(Email string) ([]interface{}, string, error) {
	collection := tr.DBClient.Database(tr.DBName).Collection(UCollection)
	var data []interface{}
	cur, err := collection.Find(context.TODO(), bson.M{"Email": Email})
	if err != nil {
		openlog.Error(err.Error())
		return nil, "106", err
	}
	for cur.Next(context.TODO()) {
		var elem interface{}
		err := cur.Decode(&elem)
		if err != nil {
			openlog.Error(err.Error())
			return nil, "106", err
		}
		data = append(data, elem)
	}
	if err := cur.Err(); err != nil {
		openlog.Error(err.Error())
		return nil, "106", err
	}
	cur.Close(context.TODO())
	return data, "", nil
}

func (tr *Repository) CreateEndUser(data interface{}) (interface{}, string, error) {
	collection := tr.DBClient.Database(tr.DBName).Collection(UCollection)
	res, insertErr := collection.InsertOne(context.TODO(), data)
	if insertErr != nil {
		openlog.Error(insertErr.Error())
		return nil, "108", insertErr
	}
	var res1 = map[string]interface{}{
		"_id": res.InsertedID,
	}
	return res1, "109", nil
}

func (tr *Repository) CreateTransaction(data interface{}) (interface{}, string, error) {
	collection := tr.DBClient.Database(tr.DBName).Collection(CollectionName)
	res, insertErr := collection.InsertOne(context.TODO(), data)
	if insertErr != nil {
		openlog.Error(insertErr.Error())
		return nil, "102", insertErr
	}
	fmt.Println(res.InsertedID)
	var res1 = map[string]interface{}{
		"_id": res.InsertedID,
	}
	return res1, "103", nil
}

func (tr *Repository) DeleteAllTransactions() (interface{}, string, error) {
	collection := tr.DBClient.Database(tr.DBName).Collection(CollectionName)
	_, err := collection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		openlog.Error(err.Error())
		return nil, "113", err
	}
	return nil, "114", nil
}

func (tr *Repository) UpdateLocation(uid string, payload map[string]interface{}) (interface{}, string, error) {
	collection := tr.DBClient.Database(tr.DBName).Collection(UCollection)
	objId, err1 := primitive.ObjectIDFromHex(uid)
	if err1 != nil {
		openlog.Error(err1.Error())
		return nil, "125", err1
	}
	_, err := collection.UpdateMany(context.TODO(), bson.M{"_id": objId}, bson.M{"$set": bson.M{"city": payload["city"].(string), "resetLocation": payload["resetLocation"].(bool)}})
	if err != nil {
		openlog.Error(err.Error())
		return nil, "115", err
	}
	return nil, "", nil
}

func (tr *Repository) GetUserByID(Id string) (map[string]interface{}, string, error) {
	collection := tr.DBClient.Database(tr.DBName).Collection(UCollection)
	var data map[string]interface{}
	objId, err1 := primitive.ObjectIDFromHex(Id)
	if err1 != nil {
		openlog.Error(err1.Error())
		return nil, "120", err1
	}
	err := collection.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&data)
	if err != nil {
		openlog.Error(err.Error())
		return nil, "119", err
	}
	return data, "", nil
}

func (tr *Repository) GetAllTransactions(filter primitive.M) ([]map[string]interface{}, string, error) {
	collection := tr.DBClient.Database(tr.DBName).Collection(CollectionName)
	var data []map[string]interface{}
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		openlog.Error(err.Error())
		return nil, "121", err
	}
	for cur.Next(context.TODO()) {
		var elem map[string]interface{}
		err := cur.Decode(&elem)
		if err != nil {
			openlog.Error(err.Error())
			return nil, "121", err
		}
		data = append(data, elem)
	}
	if err := cur.Err(); err != nil {
		openlog.Error(err.Error())
		return nil, "121", err
	}
	cur.Close(context.TODO())
	return data, "", nil
}
