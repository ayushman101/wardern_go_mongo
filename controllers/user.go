package controllers

import(
	"errors"
	"fmt"
	"net/http"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"encoding/json"
	"context"
	"log"
	"github.com/ayushman101/warden_go_mongo/models"
	//"github.com/dgrijalva/jwt-go"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)


var key string ="fejofjeaje335931jfjj3o"


//for user queries
type UserController struct{
	Client *mongo.Client
}

func NewUserController(c *mongo.Client) *UserController{
	return &UserController{c}
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request){


	var user models.User

	err:= json.NewDecoder(r.Body).Decode(&user)
	if err!=nil {
		fmt.Println(err)
		return 
	}

	collection:= uc.Client.Database("go_test_db").Collection("users")

	result,err:= collection.InsertOne(context.Background(),user)

	if err!=nil{
		 fmt.Println(err)
		 return
	}

	tok,err:= signJWT(result.InsertedID);
	
	if err!=nil{
		 log.Fatal(err)
	}

	json.NewEncoder(w).Encode(tok)

}

func (uc UserController) Allusers(w http.ResponseWriter, r *http.Request){

	//Validating the User
	tokenString := r.Header.Get("Authorization")

	ss := strings.Split(tokenString," ");

	if ss[0]!="Bearer"{
		log.Fatal(errors.New("No Bearer"))
	}

	tokenString=ss[1]

    	if tokenString == "" {
        	w.WriteHeader(http.StatusUnauthorized)
        	return
	}

    	// Validate the JWT token.
    	claims, err := validateJWT(tokenString, key)
    	if err != nil {
        	w.WriteHeader(http.StatusUnauthorized)
        	return
    	}

	fmt.Println(claims)




	collection := uc.Client.Database("go_test_db").Collection("users")

	
	//second argument is of type Document which is defined in mongo-bson.
	cursor,err:= collection.Find(context.Background(),bson.D{{}})  //second argument is a query filer
								       // empty for returnig all users
	if err!=nil{
		fmt.Println("1",err)
		log.Fatal(err)
	}

	var users []models.User

	//iterate over the cursor to get all the users
	for cursor.Next(context.Background()){

		var user models.User

		err = cursor.Decode(&user)

		if err!=nil{
			fmt.Println("2")
	
			log.Fatal(err)
		}

		users=append(users,user)
	}


	//close cursor
	err= cursor.Close(context.Background())
	
	if err!=nil{
		fmt.Println("3")
	
		log.Fatal(err)
		return
	}

	//return users in json
	json.NewEncoder(w).Encode(users)

}





//sign a JWT token

func signJWT(id interface{}) (string,error){
	token:=jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":id,
	})

	tokenString, err:= token.SignedString([]byte(key))

	if err!=nil{
		log.Fatal(err)
	}

	return tokenString, nil
}

//Validate a token

func validateJWT(tokenString string, signingKey string) (jwt.MapClaims,error){

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// Don't forget to validate the alg is what you expect:
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	// signingKey is a []byte containing your secret, e.g. []byte("my_secret_key")
	return []byte(signingKey), nil
})
	
	if err!=nil{
		fmt.Println("error in parsing token")
		log.Fatal(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid{
    		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}



