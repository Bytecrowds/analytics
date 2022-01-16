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
		Room string
		Data struct {
			BitecrowdText struct{
				Type string
				Content string
			}

		}
	}

	type StoredBitecrowd struct {
		Name string
		Text string
	}

    r := chi.NewRouter()

	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))

	bitecrowds := client.Database("testingDB").Collection("bitecrowds")

    r.Use(middleware.Logger)

	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		var data Bitecrowd 
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			 http.Error(w, err.Error(), http.StatusBadRequest)
			 return
		}

		bitecrowd := bson.D{{"name", data.Room}, {"text", data.Data.BitecrowdText.Content}}
		modifiedBitecrowd := bson.D{{"$set", bson.D{{"text", data.Data.BitecrowdText.Content}}}}
		filter := bson.D{{"name", data.Room}}

		var result StoredBitecrowd
		bitecrowds.FindOne(context.TODO(), filter).Decode(&result)

		if result.Name != "" {
			bitecrowds.UpdateOne(context.TODO(), filter, modifiedBitecrowd)
		} else {
			result, _ := bitecrowds.InsertOne(context.TODO(), bitecrowd)
			if result != nil {
				w.Write([]byte("Bitecrowd updated!"))
			}
		}
		
	})

    http.ListenAndServe("127.0.0.1:5000", r)
}
