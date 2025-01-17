package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RG-7/goformonogo/models"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct {
	client *mongo.Client
}

func NewUserController(client *mongo.Client) *UserController {
	return &UserController{client}
}

// GetUser retrieves a user by ID
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	collection := uc.client.Database("goformongo").Collection("users")

	var user models.User
	filter := bson.M{"_id": id}

	err := collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	uj, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(uj)
}

// CreateUser adds a new user to the database
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Id = fmt.Sprintf("%v", time.Now().UnixNano())

	collection := uc.client.Database("goformongo").Collection("users")

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	uj, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(uj)
}

// DeleteUser removes a user by ID
func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	collection := uc.client.Database("goformongo").Collection("users")

	filter := bson.M{"_id": id}

	_, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted User: %s\n", id)
}
