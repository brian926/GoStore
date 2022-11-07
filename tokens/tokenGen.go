package token

import (
	"context"
	"gostore/database"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserData *mongo.Collection = database.UserData(database.Client, "Users")
var SECRET_KEY = os.Getenv("SECRET_KEY")

type SignedDetails struct {
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

func TokenGenerator(email string, firstName string, lastName string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		Email:      email,
		First_Name: firstName,
		Last_Name:  lastName,
		Uid:        uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return
	}

	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "The token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token is already expired"
		return
	}

	return claims, msg
}

func UpdateAllTokens(signedToken string, signedRefreshToken string, userid string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateobj primitive.D

	updateobj = append(updateobj, bson.E{Key: "token", Value: signedToken})
	updateobj = append(updateobj, bson.E{Key: "refresh_token", Value: signedRefreshToken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateobj = append(updateobj, bson.E{Key: "updatedat", Value: updated_at})

	upsert := true
	filter := bson.M{"user_id": userid}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := UserData.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updateobj},
	}, &opt)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}
}
