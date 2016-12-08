package postgresql

import (
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"os"

	//	log15 "gopkg.in/inconshreveable/log15.v2"
)

const (
	maxConnectionCount = 40
)

var (
	Pool *pgx.ConnPool
	tableCreated bool = false
)

// afterConnect creates the prepared statements that this application uses
func afterConnect(conn *pgx.Conn) (err error) {
	if !tableCreated {
		err = prepareTestTable(conn)
		if err != nil {
			return
		}
		tableCreated = true
	}

	_, err = conn.Prepare("insert", `
    INSERT INTO account (email,name,password,verified,sex,gold,version,friends) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("size", `
    SELECT count(*) FROM account
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("query_id", `
    SELECT email FROM account WHERE _id = $1
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("queries", `
    SELECT * FROM account WHERE email = $1
  `)
	if err != nil {
		return
	}

	// There technically is a small race condition in doing an upsert with a CTE
	// where one of two simultaneous requests to the shortened URL would fail
	// with a unique index violation. As the point of this demo is pgx usage and
	// not how to perfectly upsert in PostgreSQL it is deemed acceptable.
	_, err = conn.Prepare("update", `
    UPDATE account SET password = $1, version = version + 1 WHERE _id = $2
  `)
	if err != nil {
		return
	}

	_, err = conn.Prepare("array", `
    UPDATE account SET friends = array_append(friends, $1), version = version + 1 WHERE _id = $2
  `)
	return
}

func InitDatabase() {
	var err error
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "192.168.0.118",
			User:     "postgres",
			Password: "",
			Database: "sdCloud",
			//	Logger:   log15.New("module", "pgx"),
		},
		MaxConnections: maxConnectionCount,
		AfterConnect:   afterConnect,
	}
	Pool, err = pgx.NewConnPool(connPoolConfig)
	if err != nil {
		log.Fatal("Unable to create connection pool:", err)
		os.Exit(1)
	}
}

func prepareTestTable(conn *pgx.Conn) error {
	createTableSql :=
			`-- Table: public.account
	BEGIN;

 	DROP TABLE IF EXISTS public.account;

	CREATE TABLE public.account
	(
		_id serial,
		email character varying(64) COLLATE pg_catalog."default" NOT NULL,
		name character varying(32) COLLATE pg_catalog."default" NOT NULL,
		password character varying(16) COLLATE pg_catalog."default" NOT NULL,
		verified boolean NOT NULL,
		sex boolean NOT NULL,
		gold integer NOT NULL,
		version integer NOT NULL,
		friends integer[] NOT NULL,
		CONSTRAINT account_pkey PRIMARY KEY (_id)
	)
	WITH (
		OIDS = FALSE
	)
	TABLESPACE pg_default;

	ALTER TABLE public.account
		OWNER to postgres;

	-- Index: email_index

	DROP INDEX IF EXISTS public.email_index;

	CREATE INDEX email_index
		ON public.account USING btree
		(email COLLATE pg_catalog."default")
		TABLESPACE pg_default;

	-- Index: id_index

	DROP INDEX IF EXISTS public.id_index;

	CREATE INDEX id_index
		ON public.account USING btree
		(_id)
		TABLESPACE pg_default;

	COMMIT;`;

	if _, err := conn.Exec(createTableSql); err == nil {
		return nil;
	} else {
		return fmt.Errorf("Create table failed with error: %s\n", err.Error());
	}
}
