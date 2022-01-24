package main

import (
    "net/http"
	"os"

    "github.com/go-chi/chi/v5"
    "github.com/go-chi/chi/v5/middleware"

	"github.com/rs/cors"

	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"

	"encoding/json"
)

const connectionString = "mongodb+srv://Tudor:u22hfwcAxwkt@bitecrowdsmaindb.qadvh.mongodb.net/productionDB?authSource=admin&replicaSet=atlas-1x2m3n-shard-0&w=majority&readPreference=primary&appname=MongoDB%20Compass&retryWrites=true&ssl=true"


func main() {
	type Bytecrowd struct {
		Room string
		Data struct {
			BytecrowdText struct{
				Type string
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
		Language string
	}

    r := chi.NewRouter()

	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))

	bytecrowds := client.Database("productionDB").Collection("bytecrowds")
	languages := client.Database("productionDB").Collection("languages")

    r.Use(middleware.Logger)
	r.Use(cors.Default().Handler)

	r.Post("/update", func(w http.ResponseWriter, r *http.Request) {
		var data Bytecrowd 
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			 http.Error(w, err.Error(), http.StatusBadRequest)
			 return
		}

		bytecrowd := bson.D{{"name", data.Room}, {"text", data.Data.BytecrowdText.Content}}
		modifiedBytecrowd := bson.D{{"$set", bson.D{{"text", data.Data.BytecrowdText.Content}}}}
		filter := bson.D{{"name", data.Room}}

		var result StoredBytecrowd
		bytecrowds.FindOne(context.TODO(), filter).Decode(&result)

		if result.Name != "" {
			bytecrowds.UpdateOne(context.TODO(), filter, modifiedBytecrowd)
		} else {
			result, _ := bytecrowds.InsertOne(context.TODO(), bytecrowd)
			if result != nil {
				w.Write([]byte("Bytecrowd inserted!"))
			}
		}
	})

	r.Get("/get/{bytecrowd}", func(w http.ResponseWriter, r *http.Request) {
		bytecrowdName := chi.URLParam(r, "bytecrowd")
		filter := bson.D{{"name", bytecrowdName}}

		var result StoredBytecrowd
		bytecrowds.FindOne(context.TODO(), filter).Decode(&result)

		if result.Text != "" {
			w.Write([]byte(result.Text))
		} else {
			w.WriteHeader("")
		}
	})

	r.Get("/getLanguage/{bytecrowd}", func(w http.ResponseWriter, r *http.Request) {
		bytecrowdName := chi.URLParam(r, "bytecrowd")
		filter := bson.D{{"bytecrowd", bytecrowdName}}

		var result Language
		languages.FindOne(context.TODO(), filter).Decode(&result)

		if result.Language != ""{
			w.Write([]byte(result.Language))
		} else {
			w.Write([]byte(""))
		}
	})

	r.Post("/updateLanguage", func(w http.ResponseWriter, r *http.Request) {
		var data Language 
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			 http.Error(w, err.Error(), http.StatusBadRequest)
			 return
		}

		language := bson.D{{"bytecrowd", data.Bytecrowd}, {"language", data.Language}}
		modifiedLanguage := bson.D{{"$set", bson.D{{"language", data.Language}}}}
		filter := bson.D{{"bytecrowd", data.Bytecrowd}}

		var result Language
		languages.FindOne(context.TODO(), filter).Decode(&result)

		if result.Bytecrowd != "" {
			languages.UpdateOne(context.TODO(), filter, modifiedLanguage)
		} else {
			result, _ := languages.InsertOne(context.TODO(), language)
			if result != nil {
				w.Write([]byte("Language inserted!"))
			}
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
    	port = "5000"
	}

    http.ListenAndServe(":" + port, r)
}
