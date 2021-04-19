package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var uri string = "mongodb+srv://root:asdf@ssk.3hxej.mongodb.net/GOAPIPROJDB?retryWrites=true&w=majority&authSource=admin"

type User struct {
	Name string
	Age  int
}

//连接数据库操作
func connectToMongoDB(uri string, DBName string, timeout time.Duration, numOfConnection uint64, collectionName string) (*mongo.Collection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	o := options.Client().ApplyURI(uri)
	o.SetMaxPoolSize(numOfConnection)
	client, err := mongo.Connect(ctx, o)
	if err != nil {
		return nil, err
	}

	return client.Database(DBName).Collection(collectionName), nil
}

//C操作
func Create(collection *mongo.Collection, obj []interface{}) []interface{} {
	insertResult, err := collection.InsertMany(context.TODO(), obj)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return insertResult.InsertedIDs
}

//R操作
func Read(collection *mongo.Collection /*, obj []interface{}*/) []interface{} {

	cur, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var results []interface{}
	for cur.Next(context.TODO()) {
		// 创建一个值，将单个文档解码为该值
		var elem interface{}
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// 完成后关闭游标
	cur.Close(context.TODO())
	return results
}

//U操作
func Update(collection *mongo.Collection, fliter, obj *interface{}) *mongo.UpdateResult {

	//log.Fatal("not implemented Error")
	updateResult, err := collection.UpdateMany(context.TODO(), fliter, obj)
	if err != nil {
		log.Fatal(err)
	}
	return updateResult
}

//D操作
func Delete(collection *mongo.Collection, fliter interface{}) *mongo.DeleteResult {
	deleteResult, err := collection.DeleteMany(context.TODO(), fliter)
	if err != nil {
		log.Fatal(err)
	}
	return deleteResult
}

func main() {
	lr := LogRecord{
		jobName: "append",
		command: "-0t",
		err:     "no",
		content: "shabi",
		tp: TimePoint{
			startTime: 1,
			endTime:   2,
		},
	}
	var iResult *mongo.InsertOneResult
	var err error
	collection, _ := connectToMongoDB(uri, "GOAPIPROJDB", 10*time.Second, 5, "test")
	if iResult, err = collection.InsertOne(context.TODO(), &lr); err != nil {
		fmt.Print(err)
		return
	}
	//_id:默认生成一个全局唯一ID
	id := iResult.InsertedID.(primitive.ObjectID)
	fmt.Println("自增ID", id.Hex())
}

//