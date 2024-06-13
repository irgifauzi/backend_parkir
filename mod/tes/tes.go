package tes

import "github.com/whatsauth/itmodel"

func Tes(Pesan itmodel.IteungMessage) (reply string) {

	return "Hai.. hai.. ini nama grup :\n" + Pesan.Group_name + "\nsimpan dan catat baik baik ya.. \nmakasih"
}
