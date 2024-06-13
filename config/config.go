package config

import (
	"log"

	"github.com/gocroot/helper"
	"github.com/gocroot/helper/atdb"
	"github.com/gocroot/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var IPPort, Net = helper.GetAddress()

func SetEnv() {
	if ErrorMongoconn != nil {
		log.Println(ErrorMongoconn.Error())
	}
	profile, err := atdb.GetOneDoc[model.Profile](Mongoconn, "profile", primitive.M{})
	if err != nil {
		log.Println(err)
	}
	PublicKeyWhatsAuth = profile.PublicKey
	WAAPIToken = profile.Token
}
