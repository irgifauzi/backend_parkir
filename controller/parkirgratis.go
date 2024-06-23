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
	id, err := primitive.ObjectIDFromHex("6661898bb85c143abc747d03")
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
	var tempat model.Tempat
	err := json.NewDecoder(req.Body).Decode(&tempat)
	if err != nil {
		var respn model.Response
		respn.Response = err.Error()
		helper.WriteJSON(respw, http.StatusForbidden, respn)
		return
		
	}
	temp, err:= atdb.GetOneDoc[model.Tempat](config.Mongoconn,"tempat json",bson.M{"_id":tempat.ID})
	if err != nil {
		var respn model.Response
		respn.Response = err.Error()
		helper.WriteJSON(respw, http.StatusForbidden, respn)
		return
		
	}
	temp.Nama_Tempat = tempat.Nama_Tempat
	temp.Lokasi = tempat.Lokasi
	temp.Fasilitas = tempat.Fasilitas
	temp.Lon = tempat.Lon
	temp.Lat = tempat.Lat
	temp.Gambar = tempat.Gambar
	_, err= atdb.ReplaceOneDoc(config.Mongoconn,"tempat json",bson.M{"_id":tempat.ID},temp)
	if err != nil {
		var respn model.Response
		respn.Response = err.Error()
		helper.WriteJSON(respw, http.StatusForbidden, respn)
		return
		
	}
	helper.WriteJSON(respw, http.StatusOK, temp)
	
}

func DeleteTempatParkir(respw http.ResponseWriter, req *http.Request) {
	var tempat model.Tempat
	err := json.NewDecoder(req.Body).Decode(&tempat)
	if err != nil {
		var respn model.Response
		respn.Response = err.Error()
		helper.WriteJSON(respw, http.StatusForbidden, respn)
		return

	}
	err = atdb.DeleteOneDoc(config.Mongoconn, "tempat json", bson.M{"_id": tempat.ID})
	if err != nil {
		var respn model.Response
		respn.Response = err.Error()
		helper.WriteJSON(respw, http.StatusForbidden, respn)
		return

	}
	parkirgratis, err := atdb.GetAllDoc[[]model.Tempat](config.Mongoconn, "tempat json", bson.M{"_id": tempat.ID})
	if err != nil {
		var respn model.Response
		respn.Response = err.Error()
		helper.WriteJSON(respw, http.StatusForbidden, respn)
		return

	}
	helper.WriteJSON(respw, http.StatusOK, parkirgratis)

}
