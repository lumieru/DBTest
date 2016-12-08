package main

import (
	"log"
	"net/http"
	"gopkg.in/mgo.v2/bson"
	"github.com/valyala/fasthttp"
	"github.com/lumieru/DBTest/common"
	"os"
	"github.com/lumieru/DBTest/common/mongodb"
	"fmt"
	"gopkg.in/mgo.v2"
)

var (
	globalSession *mgo.Session
)

func main() {
	globalSession = mongodb.InitDatabase()

	s := &fasthttp.Server{
		Handler: mainHandler,
		Name:    "go",
	}
	ln := common.GetListener()
	if err := s.Serve(ln); err != nil {
		log.Fatalf("Error when serving incoming connections: %s", err)
	}
}

func mainHandler(ctx *fasthttp.RequestCtx) {
	path := ctx.Path()
	switch string(path) {
	case "/insert":
		insertHandler(ctx)
	case "/queries":
		queriesHandler(ctx)
	case "/update":
		updateHandler(ctx)
	case "/array":
		arrayHandler(ctx)
	case "/size":
		sizeHandler(ctx)
	case "/exit":
		os.Exit(0)
	default:
		ctx.Error("unexpected path", fasthttp.StatusBadRequest)
	}
}

func insertHandler(ctx *fasthttp.RequestCtx) {
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
		ctx.Response.SetStatusCode(http.StatusOK)
		common.RowsInserted ++
		ctx.Response.AppendBodyString(fmt.Sprint("rowsInserted=", common.RowsInserted))
	} else {
		ctx.Error(fmt.Sprintf("insert failed:%s\n", err.Error()) , http.StatusInternalServerError)
	}
}

func sizeHandler(ctx *fasthttp.RequestCtx) {
	sess := globalSession.Clone()
	defer sess.Close()

	count, err := sess.DB("sdCloud").C("account").Count()
	switch err {
	case nil:
		ctx.Response.AppendBodyString(fmt.Sprint("size=", count))
	default:
		ctx.Error(fmt.Sprintf("query size failed:%s\n", err.Error()), http.StatusInternalServerError)
	}
}

func queriesHandler(ctx *fasthttp.RequestCtx) {
	sess := globalSession.Clone()
	defer sess.Close()

	var email mongodb.Email
	targetID := common.RandID()

	err := sess.DB("sdCloud").C("account").Find(bson.M{"_id":targetID}).Select(bson.M{"em": 1, "_id":0}).One(&email)

	switch err {
	case nil:
	default:
		ctx.Error(fmt.Sprintf("query _id with %d failed:%s\n", targetID, err.Error()), http.StatusInternalServerError)
	}

	var player mongodb.Player
	err = sess.DB("sdCloud").C("account").Find(bson.M{"em":email.Email}).One(&player)
	switch err {
	case nil:
		ctx.Response.SetStatusCode(http.StatusOK)
		ctx.Response.AppendBodyString(fmt.Sprintf("_id=%d,email=%s,name=%s,password=%s,verified=%v,sex=%v,gold=%d,version=%d,friends=%v\n",
			player.ID,player.Email,player.Name,player.Pass,player.Verify,player.Sex,player.Gold,player.Version,player.Friends))
	default:
		ctx.Error(fmt.Sprintf("query email with %s failed:%s\n", email.Email, err.Error()) , http.StatusInternalServerError)
	}
}

func updateHandler(ctx *fasthttp.RequestCtx) {
	sess := globalSession.Clone()
	defer sess.Close()

	err := sess.DB("sdCloud").C("account").UpdateId(common.RandID(), bson.M{"$set":bson.M{"pw":common.RandString(16)}})
	if err == nil {
		ctx.Response.SetStatusCode(http.StatusOK)
	} else {
		ctx.Error("Internal server error", http.StatusInternalServerError)
	}
}

func arrayHandler(ctx *fasthttp.RequestCtx) {
	sess := globalSession.Clone()
	defer sess.Close()

	err := sess.DB("sdCloud").C("account").UpdateId(common.RandID(), bson.M{"$push":bson.M{"fd":common.RandInt32()}})
	if err == nil {
		ctx.Response.SetStatusCode(http.StatusOK)
	} else {
		ctx.Error("Internal server error", http.StatusInternalServerError)
	}
}
