package database

import (
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"os"

	"context"
)

type Bytecrowd struct {
	Room string
	Data struct {
		BytecrowdText struct {
			Type    string
			Content string
		}
	}
}

type StoredBytecrowd struct {
	Name string
	Text string
}

type Language struct {
	Bytecrowd string
	Language  string
}

var database, gasit = os.LookupEnv("DATABASE")
var _, __ = fmt.Print(database, gasit)
var client, _ = mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("CONNECTION_STRING")))
var Bytecrowds = client.Database(database).Collection("bytecrowds")
var Languages = client.Database(database).Collection("languages")
