package idname

import (
	"context"
	"fmt"

	"github.com/gocroot/helper/atdb"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func IDName(Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	longitude := fmt.Sprintf("%f", Pesan.Longitude)
	latitude := fmt.Sprintf("%f", Pesan.Latitude)

	return "Hai.. hai.. kakak atas nama:\n" + Pesan.Alias_name + "\nLongitude: " + longitude + "\nLatitude: " + latitude + "\nLokasi:" + GetLokasi(db, Pesan.Longitude, Pesan.Latitude) + "\nberhasil absen\nmakasih"
}

func GetLokasi(mongoconn *mongo.Database, long float64, lat float64) (namalokasi string) {
	lokasicollection := mongoconn.Collection("lokasi")
	filter := bson.M{
		"batas": bson.M{
			"$geoIntersects": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": []float64{long, lat},
				},
			},
		},
	}
	var lokasi Lokasi
	err := lokasicollection.FindOne(context.TODO(), filter).Decode(&lokasi)
	if err != nil {
		fmt.Printf("GetLokasi: %v\n", err)
	}
	return lokasi.Nama
}

func Postdatapresensi(mongoconn *mongo.Database, long float64, lat float64) {
	lokasi := GetLokasi(mongoconn, long, lat)
	data := DataPresensi{
		Lokasi: lokasi,
	}
	atdb.InsertOneDoc(mongoconn, "presensi", data)
}
