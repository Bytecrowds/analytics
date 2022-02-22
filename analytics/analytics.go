package analytics

import (
	"bytecrowds-database-server/database"
	"bytecrowds-database-server/database/models"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"

	"context"
)

var IPanalytics = database.IPanalytics

func InterceptRequest(ginContext *gin.Context) {
	var data, result models.Request
	ginContext.BindJSON(&data)

	filter := bson.D{{"ip", data.IP}}
	IPstat := bson.D{{"ip", data.IP}, {"hits", 0}}

	IPanalytics.FindOne(context.TODO(), filter).Decode(&result)

	if result.IP == "" {
		IPanalytics.InsertOne(context.TODO(), IPstat)
	}
}
