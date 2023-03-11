package dbUser

import (
	"CoinAI.net/server/global/dbType"
	"go.mongodb.org/mongo-driver/bson"
)

func (dbObj *AccountType) Update() {
	db := dbObj.DB
	var result dbType.UserTable
	FK := bson.D{{
		Key:   "UserID",
		Value: dbObj.UserID,
	}}
	db.Table.FindOne(db.Ctx, FK).Decode(&result)
	dbObj.Data = result
}
