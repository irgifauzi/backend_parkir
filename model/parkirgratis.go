package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Tempat struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nama_Tempat string             `bson:"nama_tempat,omitempty" json:"nama_tempat,omitempty"`
	Lokasi      string             `bson:"lokasi,omitempty" json:"lokasi,omitempty"`
	Fasilitas   string             `bson:"fasilitas,omitempty" json:"fasilitas,omitempty"`
	Lon         float64            `bson:"lon,omitempty" json:"lon,omitempty"`
	Lat         float64            `bson:"lat,omitempty" json:"lat,omitempty"`
	Gambar      string            `bson:"gambar,omitempty" json:"gambar,omitempty"`
}

type Koordinat struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Markers [][]float64 `json:"markers"`
}
