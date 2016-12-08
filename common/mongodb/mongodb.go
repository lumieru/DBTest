package mongodb

import (
	"gopkg.in/mgo.v2"
	"log"
	"os"
	"sync/atomic"
)

var (
	id 	int32 = 0
)

type Email struct {
	Email string `bson:"em"`
}

type Player struct {
	ID 		int32 `bson:"_id"`
	Email	string `bson:"em"`
	Name 	string `bson:"nm"`
	Pass 	string `bson:"pw"`
	Verify	bool `bson:"ev"`
	Sex 	bool `bson:"sx"`
	Gold 	int32 `bson:"gd"`
	Version int32 `bson:"vn"`
	Friends []int32 `bson:"fd"`
}

func InitDatabase() *mgo.Session {

	Sess, err := mgo.Dial("mongodb://192.168.0.118")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
		os.Exit(1)
	} else {
		log.Print("Connected to server.\n")
	}

	Sess.SetSafe(&mgo.Safe{WMode:"majority",J:true})

	db := Sess.DB("sdCloud")
	coll := db.C("account")

	/*
		_id int
		em 	string
		nm	string
		pw	string
		ev	bool
		sx	bool
		gd	int
		vn	int
		fd	[]int
	 */

	//drop old collection
	coll.DropCollection()

	//drop old index
	coll.DropIndex("_id", "em")

	//_id is indexed automatically

	//ensure new index
	index := mgo.Index{
	 Key: []string{"em"},
	 Unique: true,
	 DropDups: true,
	}
	err = coll.EnsureIndex(index)
	if err != nil {
		log.Fatalf("Error ensure index: %v", err)
		os.Exit(1)
	}

	return Sess
}

func ID() int32 {
	return atomic.AddInt32(&id, 1)
}
