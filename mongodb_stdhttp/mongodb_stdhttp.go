package main

import (
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"github.com/lumieru/DBTest/common"
	"github.com/lumieru/DBTest/common/mongodb"
	"os"
	"net/http"
)

var (
	globalSession *mgo.Session
)

func main() {
	globalSession = mongodb.InitDatabase()

	http.HandleFunc("/insert", insertHandler)
	http.HandleFunc("/queries", queriesHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/array", arrayHandler)
	http.HandleFunc("/size", sizeHandler)
	http.HandleFunc("/exit", exitHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	sess := globalSession.Clone()
	defer sess.Close()

	err := sess.DB("sdCloud").C("account").Insert(bson.M{
		"_id":mongodb.ID(),
		"em":common.RandString(62),
		"nm":common.RandString(32),
		"pw":common.RandString(16),
		"ev":common.RandBool(),
		"sx":common.RandBool(),
		"gd":common.RandInt32(),
		"vn":common.RandInt32(),
		"fd":common.RandInt32Array(),
	})

	if err == nil {
		common.RowsInserted ++

		fmt.Fprint(w, "rowsInserted=", common.RowsInserted)
	} else {
		http.Error(w, fmt.Sprintf("insert failed:%s\n", err.Error()) , http.StatusInternalServerError)
	}
}

func sizeHandler(w http.ResponseWriter, r *http.Request) {
	sess := globalSession.Clone()
	defer sess.Close()

	count, err := sess.DB("sdCloud").C("account").Count()
	switch err {
	case nil:
		fmt.Fprint(w, "size=", count)
	default:
		http.Error(w, fmt.Sprintf("query size failed:%s\n", err.Error()), http.StatusInternalServerError)
	}
}

func queriesHandler(w http.ResponseWriter, r *http.Request) {
	sess := globalSession.Clone()
	defer sess.Close()

	var email mongodb.Email
	targetID := common.RandID()

	err := sess.DB("sdCloud").C("account").Find(bson.M{"_id":targetID}).Select(bson.M{"em": 1, "_id":0}).One(&email)

	switch err {
	case nil:
	default:
		http.Error(w, fmt.Sprintf("query _id with %d failed:%s\n", targetID, err.Error()), http.StatusInternalServerError)
	}

	var player mongodb.Player
	err = sess.DB("sdCloud").C("account").Find(bson.M{"em":email.Email}).One(&player)
	switch err {
	case nil:
		fmt.Fprintf(w, "_id=%d,email=%s,name=%s,password=%s,verified=%v,sex=%v,gold=%d,version=%d,friends=%v\n",
			player.ID,player.Email,player.Name,player.Pass,player.Verify,player.Sex,player.Gold,player.Version,player.Friends)
	default:
		http.Error(w, fmt.Sprintf("query email with %s failed:%s\n", email.Email, err.Error()) , http.StatusInternalServerError)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	sess := globalSession.Clone()
	defer sess.Close()

	err := sess.DB("sdCloud").C("account").UpdateId(common.RandID(), bson.M{"$set":bson.M{"pw":common.RandString(16)}})
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func arrayHandler(w http.ResponseWriter, r *http.Request) {
	sess := globalSession.Clone()
	defer sess.Close()

	err := sess.DB("sdCloud").C("account").UpdateId(common.RandID(), bson.M{"$push":bson.M{"fd":common.RandInt32()}})
	if err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func exitHandler(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}