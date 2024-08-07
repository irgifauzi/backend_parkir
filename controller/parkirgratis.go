package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"net/http"

	"github.com/gocroot/config"
	"github.com/gocroot/helper"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetLokasi(respw http.ResponseWriter, req *http.Request) {
	var resp itmodel.Response
	kor, err := atdb.GetAllDoc[[]model.Tempat](config.Mongoconn, "tempat", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusBadRequest, resp)
		return
	}
	helper.WriteJSON(respw, http.StatusOK, kor)
}

func GetMarker(respw http.ResponseWriter, req *http.Request) {
	var resp itmodel.Response
	mar, err := atdb.GetOneLatestDoc[model.Koordinat](config.Mongoconn, "marker", bson.M{})
	if err != nil {
		resp.Response = err.Error()
		helper.WriteJSON(respw, http.StatusBadRequest, mar)
		return
	}
	helper.WriteJSON(respw, http.StatusOK, mar)
}

func PostTempatParkir(respw http.ResponseWriter, req *http.Request) {

	var tempatParkir model.Tempat
	if err := json.NewDecoder(req.Body).Decode(&tempatParkir); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: err.Error()})
		return
	}

	if tempatParkir.Gambar != "" {
		tempatParkir.Gambar = "https://raw.githubusercontent.com/parkirgratis/filegambar/main/img/" + tempatParkir.Gambar
	}

	result, err := config.Mongoconn.Collection("tempat").InsertOne(context.Background(), tempatParkir)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, itmodel.Response{Response: err.Error()})
		return
	}

	insertedID := result.InsertedID.(primitive.ObjectID)

	helper.WriteJSON(respw, http.StatusOK, itmodel.Response{Response: fmt.Sprintf("Tempat parkir berhasil disimpan dengan ID: %s", insertedID.Hex())})
}

func PostKoordinat(respw http.ResponseWriter, req *http.Request) {
	var newKoor model.Koordinat
	if err := json.NewDecoder(req.Body).Decode(&newKoor); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, err.Error())
		return
	}

	// Set the specific ID you want to update
	id, err := primitive.ObjectIDFromHex("6679b77450a939208a4a7a28")
	if err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// Create filter and update fields
	filter := bson.M{"_id": id}
	update := bson.M{"$push": bson.M{"markers": bson.M{"$each": newKoor.Markers}}}

	if _, err := atdb.UpdateDoc(config.Mongoconn, "marker", filter, update); err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, err.Error())
		return
	}
	helper.WriteJSON(respw, http.StatusOK, "Markers updated")
}

func PutTempatParkir(respw http.ResponseWriter, req *http.Request) {
	var newTempat model.Tempat
	if err := json.NewDecoder(req.Body).Decode(&newTempat); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("Decoded document:", newTempat)

	if newTempat.ID.IsZero() {
		helper.WriteJSON(respw, http.StatusBadRequest, "ID is required")
		return
	}

	filter := bson.M{"_id": newTempat.ID}
	update := bson.M{"$set": newTempat}
	fmt.Println("Filter:", filter)
	fmt.Println("Update:", update)

	result, err := atdb.UpdateDoc(config.Mongoconn, "tempat", filter, update)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, err.Error())
		return
	}

	if result.ModifiedCount == 0 {
		helper.WriteJSON(respw, http.StatusNotFound, "Document not found or not modified")
		return
	}

	helper.WriteJSON(respw, http.StatusOK, newTempat)
}

func DeleteTempatParkir(respw http.ResponseWriter, req *http.Request) {
	var requestBody struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(req.Body).Decode(&requestBody); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Invalid JSON data"})
		return
	}

	objectId, err := primitive.ObjectIDFromHex(requestBody.ID)
	if err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Invalid ID format"})
		return
	}

	filter := bson.M{"_id": objectId}

	deletedCount, err := atdb.DeleteOneDoc(config.Mongoconn, "tempat", filter)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to delete document", "error": err.Error()})
		return
	}

	if deletedCount == 0 {
		helper.WriteJSON(respw, http.StatusNotFound, map[string]string{"message": "Document not found"})
		return
	}

	helper.WriteJSON(respw, http.StatusOK, map[string]string{"message": "Document deleted successfully"})
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func AdminLogin(respw http.ResponseWriter, req *http.Request) {
	var loginReq LoginRequest

	if err := json.NewDecoder(req.Body).Decode(&loginReq); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Invalid JSON data"})
		return
	}

	clientOptions := options.Client().ApplyURI(config.MongoURI) // Assuming MongoURI is defined in your config
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to connect to MongoDB", "error": err.Error()})
		return
	}
	defer client.Disconnect(context.TODO())

	adminCollection := client.Database("parkir_db").Collection("admin")

	var admin model.Admin
	filter := bson.M{"username": loginReq.Username, "password": loginReq.Password}
	err = adminCollection.FindOne(context.TODO(), filter).Decode(&admin)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			helper.WriteJSON(respw, http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
		} else {
			helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to login", "error": err.Error()})
		}
		return
	}

	helper.WriteJSON(respw, http.StatusOK, map[string]string{"message": "Login successful"})
}

func AdminRegister(respw http.ResponseWriter, req *http.Request) {
	var registerReq model.RegisterRequest

	if err := json.NewDecoder(req.Body).Decode(&registerReq); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Invalid JSON data"})
		return
	}

	if registerReq.Password != registerReq.ConfirmPassword {
		helper.WriteJSON(respw, http.StatusBadRequest, map[string]string{"message": "Passwords do not match"})
		return
	}

	clientOptions := options.Client().ApplyURI(config.MongoURI)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to connect to MongoDB", "error": err.Error()})
		return
	}
	defer client.Disconnect(context.TODO())

	adminCollection := client.Database("parkir_db").Collection("admin")

	// Check if username already exists
	var existingAdmin model.Admin
	filter := bson.M{"username": registerReq.Username}
	err = adminCollection.FindOne(context.TODO(), filter).Decode(&existingAdmin)
	if err == nil {
		helper.WriteJSON(respw, http.StatusConflict, map[string]string{"message": "Username already exists"})
		return
	} else if err != mongo.ErrNoDocuments {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to register", "error": err.Error()})
		return
	}

	admin := model.Admin{
		ID:       primitive.NewObjectID(),
		Username: registerReq.Username,
		Password: registerReq.Password,
	}

	_, err = adminCollection.InsertOne(context.TODO(), admin)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, map[string]string{"message": "Failed to register", "error": err.Error()})
		return
	}

	helper.WriteJSON(respw, http.StatusOK, map[string]string{"message": "Registration successful"})
}

func DeleteKoordinat(respw http.ResponseWriter, req *http.Request) {
	var deleteRequest struct {
		Markers [][]float64 `json:"markers"`
	}

	if err := json.NewDecoder(req.Body).Decode(&deleteRequest); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, err.Error())
		return
	}

	id, err := primitive.ObjectIDFromHex("6679b77450a939208a4a7a28")
	if err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, "Invalid ID format")
		return
	}

	filter := bson.M{"_id": id}
	update := bson.M{
		"$pull": bson.M{
			"markers": bson.M{
				"$in": deleteRequest.Markers,
			},
		},
	}

	// Perform the update
	if _, err := atdb.UpdateDoc(config.Mongoconn, "marker", filter, update); err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, err.Error())
		return
	}

	// Respond with success
	helper.WriteJSON(respw, http.StatusOK, "Coordinates deleted")
}

func PutKoordinat(respw http.ResponseWriter, req *http.Request) {
	var updateRequest struct {
		ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Markers [][]float64        `json:"markers"`
	}

	// Decode the JSON request body into the updateRequest struct
	if err := json.NewDecoder(req.Body).Decode(&updateRequest); err != nil {
		http.Error(respw, err.Error(), http.StatusBadRequest)
		return
	}

	// Use the provided ID from the request
	id := updateRequest.ID
	if id.IsZero() {
		http.Error(respw, "ID is required", http.StatusBadRequest)
		return
	}

	// Define the filter for the document
	filter := bson.M{"_id": id}

	// Establish a MongoDB connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb+srv://irgifauzi:%40Sasuke123@webservice.rq9zk4m.mongodb.net/"))
	if err != nil {
		http.Error(respw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(context.TODO())

	collection := client.Database("parkir_db").Collection("marker")

	// Fetch the current document to find the index of the marker to update
	var document struct {
		Markers [][]float64 `bson:"markers"`
	}
	if err := collection.FindOne(context.TODO(), filter).Decode(&document); err != nil {
		http.Error(respw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find the index of the marker to update
	var index int = -1
	for i, marker := range document.Markers {
		if len(marker) == 2 && marker[0] == updateRequest.Markers[0][0] && marker[1] == updateRequest.Markers[0][1] {
			index = i
			break
		}
	}

	if index == -1 {
		http.Error(respw, "Marker not found", http.StatusBadRequest)
		return
	}

	// Update the specific marker in the array
	update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("markers.%d", index): updateRequest.Markers[1],
		},
	}

	// Apply the update operation to the MongoDB collection
	if _, err := collection.UpdateOne(context.TODO(), filter, update); err != nil {
		http.Error(respw, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	respw.WriteHeader(http.StatusOK)
	respw.Write([]byte("Coordinate updated"))
}
