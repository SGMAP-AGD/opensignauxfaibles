package main

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// ServerJSFunc Function à injecter dans l'instance MongoDB
type ServerJSFunc struct {
	ID    string          `json:"id" bson:"_id"`
	Value bson.JavaScript `json:"value" bson:"value"`
}

// Add Méthode pour insérer
func (f ServerJSFunc) Add(db *mgo.Database) {
	db.C("system.js").Upsert(bson.M{"_id": f.ID}, f)
}

func declareServerFunctions(db *mgo.Database) {
	f := ServerJSFunc{
		ID:    "generatePeriodSerie",
		Value: bson.JavaScript{Code: "function (date_debut, date_fin) {var date_next = new Date(date_debut.getTime());var serie = [];while (date_next.getTime() < date_fin.getTime()) {serie.push(new Date(date_next.getTime()));date_next.setUTCMonth(date_next.getUTCMonth() + 1);}return serie;}"},
	}
	f.Add(db)
	f = ServerJSFunc{
		ID:    "compareDebit",
		Value: bson.JavaScript{Code: `function(a,b) {if (a.numero_historique < b.numero_historique) return -1;if (a.numero_historique > b.numero_historique) return 1;return 0;}`},
	}
	f.Add(db)
	f = ServerJSFunc{
		ID:    "isRJLJ",
		Value: bson.JavaScript{Code: `function(code) {codes = ['PCL010501','PCL010502','PCL030105','PCL05010102','PCL05010203','PCL05010402','PCL05010302','PCL05010502','PCL05010702','PCL05010802','PCL05010901','PCL05011003','PCL05011101','PCL05011203','PCL05011303','PCL05011403','PCL05011503','PCL05011603','PCL05011902','PCL05012003','PCL0108','PCL0109','PCL030107','PCL030108','PCL030307','PCL030308','PCL05010103','PCL05010104','PCL05010204','PCL05010205','PCL05010303','PCL05010304','PCL05010403','PCL05010404','PCL05010503','PCL05010504','PCL05010703','PCL05010803','PCL05011004','PCL05011005','PCL05011102','PCL05011103','PCL05011204','PCL05011205','PCL05011304','PCL05011305','PCL05011404','PCL05011405','PCL05011504','PCL05011505','PCL05011604','PCL05011605','PCL05011903','PCL05011904','PCL05012004','PCL05012005','PCL040802'];return codes.includes(code);}`},
	}
	f.Add(db)
	f = ServerJSFunc{
		ID:    "DateAddMonth",
		Value: bson.JavaScript{Code: `function(date, nbMonth) {var result = new Date(date.getTime());result.setUTCMonth(result.getUTCMonth() + nbMonth);return result;}`},
	}
}
