package controllers

import(
//	"fmt"
//	"net/http"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct{
	Client *mongo.Client
}

func NewUserController(c *mongo.Client) *UserController{
	return &UserController{c}
}


