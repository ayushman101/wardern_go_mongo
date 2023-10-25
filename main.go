package main

import (
	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"log"
	"fmt"
	"os"
	"context"
	"time"
	"github.com/ayushman101/warden_go_mongo/controllers"
)

func main(){

	r := chi.NewRouter()
	
	c,err:= connectToDB()

	if err!=nil{
		log.Fatal(err.Error())
		os.Exit(1)
	}
	

	uc := controllers.NewUserController(c)
	
	r.Get("/users/login",uc.LoginUser)
	r.Get("/users/all",uc.Allusers)
	r.Post("/users/register",uc.CreateUser)
	r.Post("/session/create",uc.CreateSession)

	fmt.Println("Server started at port 8080")
	log.Fatal(http.ListenAndServe(":8080",r))
}

func connectToDB() (*mongo.Client, error){
	
	clientOptions:= options.Client().ApplyURI("mongodb+srv://ayushman:av0HnjtyBeYyAFT7@cluster0.wxxniud.mongodb.net/go_test_db?retryWrites=true&w=majority")
	
	client,err := mongo.NewClient(clientOptions)

	if err!=nil{
		log.Fatal(err)
		os.Exit(1)
	}

	ctx, cancel:= context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)

	return client,nil
}


