package tokens

import (
	"context"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pradytpk/go-ecommerce/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SignedDetials datatype
type SignedDetials struct {
	Email     string
	FirstName string
	LastName  string
	UID       string
	jwt.StandardClaims
}

// UserData mongo collection for users
var UserData *mongo.Collection = database.UserData(database.Client, "Users")

// SECRET_KEY get the key from env file
var SECRET_KEY = os.Getenv("SECRET_LOVE")

// TokenGenerator for the application
//
//	@param email
//	@param firstname
//	@param lastname
//	@param uiid
//	@return signedtoken
//	@return signedrefershtoken
//	@return err
func TokenGenerator(email string, firstname string, lastname string, uiid string) (signedtoken string, signedrefershtoken string, err error) {
	claims := &SignedDetials{
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
		UID:       uiid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshclaims := &SignedDetials{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(186)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic(err)
		return "", "", err
	}
	refreshtoken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshclaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panicln(err)
		return
	}
	return token, refreshtoken, err
}

// ValidateToken check the token
//
//	@param signedtoken
//	@return claims
//	@return msg
func ValidateToken(signedtoken string) (claims *SignedDetials, msg string) {
	token, err := jwt.ParseWithClaims(signedtoken, &SignedDetials{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetials)
	if !ok {
		msg = "the token is invalid"
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
		return
	}
	return claims, msg
}

// UpdateAllTokens Update operation
//
//	@param signedtoken
//	@param signedrefreshtoken
//	@param userid
func UpdateAllTokens(signedtoken string, signedrefreshtoken string, userid string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var updateobj primitive.D
	updateobj = append(updateobj, bson.E{Key: "token", Value: signedtoken})
	updateobj = append(updateobj, bson.E{Key: "refresh_token", Value: signedrefreshtoken})
	updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateobj = append(updateobj, bson.E{Key: "updatedat", Value: updated_at})
	upsert := true
	filter := bson.M{"user_id": userid}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := UserData.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: updateobj},
	},
		&opt)
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}
}
