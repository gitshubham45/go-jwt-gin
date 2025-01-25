package controllers

import (
	"goJwt/database"
	helper "goJwt/helpers"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client , "user")
var validate = validator.New()

func HashPassword()

func VerifyPassword()

func Signup()

func Login()

func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context){
		userId := c.Params("user_id")

		if err := helper.MatchUserTypeToUid(c , userId); err != nil{
			c.JSON(http.StatusBadRequest , gin.H{"error":err.Error()})
			return
		}
	}

}

func GetUser()
