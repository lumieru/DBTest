package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx"
	"github.com/lumieru/DBTest/common/postgresql"
	"github.com/lumieru/DBTest/common"
)

func main() {
	postgresql.InitDatabase();

	http.HandleFunc("/insert", insertHandler)
	http.HandleFunc("/queries", queriesHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/array", arrayHandler)
	http.HandleFunc("/size", sizeHandler)
	http.HandleFunc("/exit", exitHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func insertHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := postgresql.Pool.Exec("insert",
		common.RandStringBytesMaskImprSrc(62), 	//email
		common.RandStringBytesMaskImprSrc(32), 	//name
		common.RandStringBytesMaskImprSrc(16),	//password
		common.RandBool(), 					  	//verified
		common.RandBool(),						//sex
		common.RandInt32(),					  	//glod
		common.RandInt32(),					  	//version
		common.RandInt32Array(),				//friends
	); err == nil {
		common.RowsInserted ++
		fmt.Fprint(w, "rowsInserted=", common.RowsInserted)
	} else {
		http.Error(w, fmt.Sprintf("insert failed:%s\n", err.Error()), http.StatusInternalServerError)
	}
}

func sizeHandler(w http.ResponseWriter, r *http.Request) {
	var count int32
	err := postgresql.Pool.QueryRow("size").Scan(&count)
	switch err {
	case nil:
		fmt.Fprint(w, "size=", count)
	case pgx.ErrNoRows:
		http.Error(w, fmt.Sprintf("query size failed:%s\n", err.Error()), http.StatusNotFound)
	default:
		http.Error(w, fmt.Sprintf("query size failed:%s\n", err.Error()), http.StatusInternalServerError)
	}
}

func queriesHandler(w http.ResponseWriter, r *http.Request) {
	var email string
	targetID := common.RandID()
	err := postgresql.Pool.QueryRow("query_id", targetID).Scan(&email)
	switch err {
	case nil:
	case pgx.ErrNoRows:
		http.Error(w, fmt.Sprintf("query _id with %d failed:%s\n", targetID, err.Error()), http.StatusNotFound)
	default:
		http.Error(w, fmt.Sprintf("query _id with %d failed:%s\n", targetID, err.Error()), http.StatusInternalServerError)
	}

	var _id int32
	var name, password string
	var verified, sex bool
	var gold, version int32
	var friends []int32
	err = postgresql.Pool.QueryRow("queries", email).Scan(&_id, &email, &name, &password, &verified, &sex, &gold, &version, &friends)
	switch err {
	case nil:
		fmt.Fprintf(w, "_id=%d,email=%s,name=%s,password=%s,verified=%v,sex=%v,gold=%d,version=%d,friends=%v\n",
			_id,email,name,password,verified,sex,gold,version,friends)
	case pgx.ErrNoRows:
		http.Error(w, fmt.Sprintf("query email with %s failed:%s\n", email, err.Error()), http.StatusNotFound)
	default:
		http.Error(w, fmt.Sprintf("query email with %s failed:%s\n", email, err.Error()), http.StatusInternalServerError)
	}
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := postgresql.Pool.Exec("update",
		common.RandStringBytesMaskImprSrc(16),
		common.RandID(),
	); err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func arrayHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := postgresql.Pool.Exec("array",
		common.RandInt32(),
		common.RandID(),
	); err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func exitHandler(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}
