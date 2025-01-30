package helpers

import (
	"context"
	"errors"
	"fmt"
	"goJwt/database"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct {
	Email      string
	First_name string
	Last_name  string
	User_type  string
	Uid        string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_name: firstName,
		Last_name:  lastName,
		Uid:        uid,
		User_type:  userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userId string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refresh_token", Value: signedRefreshToken})

	Updated_at , _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updated_at",Value: Updated_at})

	upsert := true
	filter := bson.M{"user_id":userId}
	opt := options.UpdateOptions{
		Upsert : &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter ,
		bson.D{
			{Key: "$set",Value: updateObj},
		},
		&opt,
	)
	defer cancel()

	if err != nil{
		
	}
}

func ValidateToken(signedToken string) (claims *SignedDetails, err error) {
    token, err := jwt.ParseWithClaims(
        signedToken,
        &SignedDetails{},
        func(token *jwt.Token) (interface{}, error) {
            return []byte(SECRET_KEY), nil
        },
    )

    if err != nil {
        return nil, fmt.Errorf("token parsing failed: %w", err) 
    }

    claims, ok := token.Claims.(*SignedDetails)
    if !ok {
        return nil, errors.New("invalid token claims") 
    }

    if claims.ExpiresAt < time.Now().Local().Unix() { 
        return nil, errors.New("token is expired") 
    }

    return claims, nil 
}

