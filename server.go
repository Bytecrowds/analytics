package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"os"

	"github.com/joho/godotenv"

	"github.com/rs/cors"

	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"encoding/json"
)

func main() {
	if os.Getenv("PRODUTION") != "1" && os.Getenv("PRODUCTION") != "true" {
		godotenv.Load()
	}

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

	r := chi.NewRouter()

	database := os.Getenv("DATABASE")
	client, _ := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("CONNECTION_STRING")))
	bytecrowds := client.Database(database).Collection("bytecrowds")
	languages := client.Database(database).Collection("languages")

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
			w.Write([]byte(""))
		}
	})

	r.Get("/getLanguage/{bytecrowd}", func(w http.ResponseWriter, r *http.Request) {
		bytecrowdName := chi.URLParam(r, "bytecrowd")
		filter := bson.D{{"bytecrowd", bytecrowdName}}

		var result Language
		languages.FindOne(context.TODO(), filter).Decode(&result)

		if result.Language != "" {
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

	http.ListenAndServe(":"+os.Getenv("PORT"), r)
}
