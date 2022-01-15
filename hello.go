package main

import (
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"

	"encoding/json"
)

const connectionString = "mongodb+srv://Tudor:u22hfwcAxwkt@bitecrowdsmaindb.qadvh.mongodb.net/testingDB?authSource=admin&replicaSet=atlas-1x2m3n-shard-0&w=majority&readPreference=primary&appname=MongoDB%20Compass&retryWrites=true&ssl=true"


func main() {
	type Bitecrowd struct {
		Name string
		BitecrowdText string
	}

    r := chi.NewRouter()

	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))

	bitecrowds := client.Database("testingDB").Collection("bitecrowds")

    r.Use(middleware.Logger)

    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello World!"))
    })

	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		var data Bitecrowd 

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			 http.Error(w, err.Error(), http.StatusBadRequest)
			 return
		}
		
		bitecrowd := bson.D{{"name", data.Name}, {"text", data.BitecrowdText}}
		filter := bson.D{{"name", data.Name}}

		result, _ := bitecrowds.FindOne(context.TODO(), bitecrowd)
		if result != nil {
			bitecrowds.UpdateOne(context.TODO(), filter, bitecrowd)
		} else {
			result, _ := bitecrowds.InsertOne(context.TODO(), bitecrowd)
				if result != nil {

				}
		}
	})

    http.ListenAndServe("127.0.0.1:5000", r)
}
