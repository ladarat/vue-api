package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/mgo.v2/bson"

	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://admin:clinic@cluster0-vqe7a.gcp.mongodb.net/test?retryWrites=true&w=majority"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.TODO()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("clinic")
	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(databases)

	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/login", func(c echo.Context) error {
		user := new(User)
		err := c.Bind(user)
		if err != nil {
			return err
		}

		coll := db.Collection("user")
		result := coll.FindOne(ctx, bson.M{"username": user.Username, "password": user.Password})

		resultText, err := result.DecodeBytes()
		fmt.Println(result.DecodeBytes())
		return c.String(http.StatusOK, resultText.String())
	})

	e.Logger.Fatal(e.Start(":1323"))

}
