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

// PostTempatParkir adalah fungsi yang menangani permintaan POST untuk menyimpan data tempat parkir baru.
func PostTempatParkir(respw http.ResponseWriter, req *http.Request) {
	// Membaca data dari body permintaan
	var tempatParkir model.Tempat
	if err := json.NewDecoder(req.Body).Decode(&tempatParkir); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, itmodel.Response{Response: err.Error()})
		return
	}

	// Menyisipkan data tempat ke dalam koleksi MongoDB yang ditentukan
	result, err := config.Mongoconn.Collection("tempat").InsertOne(context.Background(), tempatParkir)
	if err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, itmodel.Response{Response: err.Error()})
		return
	}

	// Mengambil ID dari dokumen yang baru disisipkan
	insertedID := result.InsertedID.(primitive.ObjectID)

	// Mengirimkan respons sukses dengan ID dari data yang baru disisipkan
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


func DeleteKoordinat(respw http.ResponseWriter, req *http.Request) {
	var deleteRequest struct {
		ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		Markers [][]float64 `json:"markers"`
	}

	if err := json.NewDecoder(req.Body).Decode(&deleteRequest); err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, err.Error())
		return
	}

	id, err := primitive.ObjectIDFromHex("6661898bb85c143abc747d03")
	if err != nil {
		helper.WriteJSON(respw, http.StatusBadRequest, "Invalid ID format")
		return
	}

	filter := bson.M{
		"_id": id,
		"markers": bson.M{
			"$in": deleteRequest.Markers,
		},
	}

	update := bson.M{
		"$pull": bson.M{
			"markers": bson.M{
				"$in": deleteRequest.Markers,
			},
		},
	}

	if _, err := atdb.UpdateDoc(config.Mongoconn, "marker", filter, update); err != nil {
		helper.WriteJSON(respw, http.StatusInternalServerError, err.Error())
		return
	}

	helper.WriteJSON(respw, http.StatusOK, "Coordinates deleted")
}

