package main

import (
	"fmt"
	"log"

	"github.com/jackc/pgx"
	"github.com/valyala/fasthttp"
	"os"
	"net/http"

	"github.com/lumieru/DBTest/common"
	"github.com/lumieru/DBTest/common/postgresql"
)


func main() {
	var err error

	postgresql.InitDatabase(false);

	s := &fasthttp.Server{
		Handler: mainHandler,
		Name:    "go",
	}
	ln := common.GetListener()
	if err = s.Serve(ln); err != nil {
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
	sql := fmt.Sprintf("INSERT INTO account (email,name,password,verified,sex,gold,version,friends) VALUES ('%s','%s','%s',%t,%t,%d,%d,'{%s}');",
		common.RandStringBytesMaskImprSrc(62), 	//email
		common.RandStringBytesMaskImprSrc(32), 	//name
		common.RandStringBytesMaskImprSrc(16),	//password
		common.RandBool(), 					  	//verified
		common.RandBool(),						//sex
		common.RandInt32(),					  	//glod
		common.RandInt32(),					  	//version
		common.RandInt32ArrayString(),				//friends
	 )
	if _, err := postgresql.Pool.Exec(sql); err == nil {
		ctx.Response.SetStatusCode(http.StatusOK)
		common.RowsInserted ++
		ctx.Response.AppendBodyString(fmt.Sprint("rowsInserted=", common.RowsInserted))
	} else {
		ctx.Error(fmt.Sprintf("insert failed:%s\n", err.Error()) , http.StatusInternalServerError)
	}
}

func sizeHandler(ctx *fasthttp.RequestCtx) {
	var count int32
	err := postgresql.Pool.QueryRow("SELECT count(*) FROM account").Scan(&count)
	switch err {
	case nil:
		ctx.Response.AppendBodyString(fmt.Sprint("size=", count))
	case pgx.ErrNoRows:
		ctx.Error(fmt.Sprintf("query size failed:%s\n", err.Error()), http.StatusNotFound)
	default:
		ctx.Error(fmt.Sprintf("query size failed:%s\n", err.Error()), http.StatusInternalServerError)
	}
}

func queriesHandler(ctx *fasthttp.RequestCtx) {
	var email string
	targetID := common.RandID()
	conn, err := postgresql.Pool.Acquire()
	if err != nil {
		ctx.Error(fmt.Sprintf("Acquire connection failed:%s\n", err.Error()), http.StatusInternalServerError)
		return
	}

	defer postgresql.Pool.Release(conn)

	err = conn.QueryRow(fmt.Sprintf("SELECT email FROM account WHERE _id=%d", targetID)).Scan(&email)
	switch err {
	case nil:
	case pgx.ErrNoRows:
		ctx.Error(fmt.Sprintf("query _id with %d failed:%s\n", targetID, err.Error()), http.StatusNotFound)
		return
	default:
		ctx.Error(fmt.Sprintf("query _id with %d failed:%s\n", targetID, err.Error()), http.StatusInternalServerError)
		return
	}

	var _id int32
	var name, password string
	var verified, sex bool
	var gold, version int32
	var friends []int32
	err = conn.QueryRow(fmt.Sprintf("SELECT * FROM account WHERE email='%s'", email)).Scan(&_id, &email, &name, &password, &verified, &sex, &gold, &version, &friends)
	switch err {
	case nil:
		ctx.Response.SetStatusCode(http.StatusOK)
		ctx.Response.AppendBodyString(fmt.Sprintf("_id=%d,email=%s,name=%s,password=%s,verified=%v,sex=%v,gold=%d,version=%d,friends=%v\n",
			_id,email,name,password,verified,sex,gold,version,friends))
	case pgx.ErrNoRows:
		ctx.Error(fmt.Sprintf("query email with %s failed:%s\n", email, err.Error()), http.StatusNotFound)
	default:
		ctx.Error(fmt.Sprintf("query email with %s failed:%s\n", email, err.Error()) , http.StatusInternalServerError)
	}
}

func updateHandler(ctx *fasthttp.RequestCtx) {
	if _, err := postgresql.Pool.Exec(fmt.Sprintf("UPDATE account SET password='%s', version=version+1 WHERE _id=%d",
		common.RandStringBytesMaskImprSrc(16),
		common.RandID(),)); err == nil {
		ctx.Response.SetStatusCode(http.StatusOK)
	} else {
		ctx.Error("Internal server error", http.StatusInternalServerError)
	}
}

func arrayHandler(ctx *fasthttp.RequestCtx) {
	if _, err := postgresql.Pool.Exec(fmt.Sprintf("UPDATE account SET friends=array_append(friends,%d), version=version+1 WHERE _id=%d",
		common.RandInt32(),
		common.RandID())); err == nil {
		ctx.Response.SetStatusCode(http.StatusOK)
	} else {
		ctx.Error("Internal server error", http.StatusInternalServerError)
	}
}
