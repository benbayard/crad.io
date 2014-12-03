package crad

import "gopkg.in/mgo.v2"

type DBConnection struct {
	Session    *mgo.Session    `json:"-" bson:"-"`
	Collection *mgo.Collection `json:"-" bson:"-"`
}

func (dbconn *DBConnection) GetDB() (*mgo.Session, *mgo.Collection) {
	if dbconn.Session == nil {
		DatabaseConnect()
	}
	return dbconn.Session, dbconn.Collection
}

// DB.GetDB()
