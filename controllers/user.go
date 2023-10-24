package controllers

import(
	"fmt"
	"net/http"
	"go.mongodb.org/mongo-driver/mongo"
	"encoding/json"
	"context"

	"github.com/ayushman101/warden_go_mongo/models"
)

type UserController struct{
	Client *mongo.Client
}

func NewUserController(c *mongo.Client) *UserController{
	return &UserController{c}
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request) error{

	var user models.User

	err:= json.NewDecoder(r.Body).Decode(&user)
	if err!=nil {
		return fmt.Errorf("Error while  decoding json: %w",err)
	}

	collection:= uc.Client.Database("go_test_db").Collection("users")

	result,err:= collection.InsertOne(context.Background(),user)

	if err!=nil{
		return fmt.Errorf("Error while creating user: %w",err)
	}

	json.NewEncoder(w).Encode(result.InsertedID)

	return nil
}
