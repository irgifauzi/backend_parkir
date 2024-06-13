package mod

import (
	"github.com/gocroot/mod/idgrup"
	"github.com/gocroot/mod/idname"
	"github.com/gocroot/mod/tes"
	"github.com/whatsauth/itmodel"
	"go.mongodb.org/mongo-driver/mongo"
)

func Caller(Modulename string, Pesan itmodel.IteungMessage, db *mongo.Database) (reply string) {
	switch Modulename {
	case "idgrup":
		reply = idgrup.IDGroup(Pesan)
	case "idname":
		reply = idname.IDName(Pesan, db)
	case "tes":
		reply = tes.Tes(Pesan)
	}
	return
}
