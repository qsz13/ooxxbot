package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/qsz13/ooxxbot/logger"
	"strconv"
	"time"
)

type DB struct {
	sqldb *sql.DB
}

func NewDB(db_uri string) *DB {
	db := &DB{sqldb: initDBConn(db_uri)}
	db.CreateTable()
	return db
}

func initDBConn(db_dsn string) *sql.DB {
	logger.Debug("Init DB connection: " + db_dsn)
	for retry, wait := 1, 1; ; retry, wait = retry+1, wait*2 {
		db, err := sql.Open("postgres", db_dsn)
		if err != nil {
			logger.Warning("DB open failed, retry times: " + strconv.Itoa(retry) + ", reason:" + err.Error())
			time.Sleep(time.Duration(wait) * time.Second)
			continue
		}
		err = db.Ping()
		if err != nil {
			logger.Warning("DB ping failed, retry times: " + strconv.Itoa(retry) + ", reason:" + err.Error())
			time.Sleep(time.Duration(wait) * time.Second)
			continue
		}
		logger.Debug("DB connection success: " + db_dsn)
		return db
	}
	logger.Error("DB connection failed.")
	return nil
}

func (db *DB) CreateTable() {
	db.createUserTable()
	db.createSubscriptionTable()
	db.createJandanTable()
	db.createSourceTable()
}

func (db *DB) createUserTable() error {
	logger.Debug("Creating Table user")
	sql_table := `
	CREATE TABLE IF NOT EXISTS "user" (
	id serial primary key,
	first_name varchar(32),
	last_name varchar(32),
	user_name varchar(32)
	) WITH (OIDS=FALSE);`
	_, err := db.sqldb.Exec(sql_table)

	if err != nil {
		logger.Debug(err)
	}
	return err
}

func (db *DB) createSubscriptionTable() error {
	logger.Debug("Creating Table subscription")
	sql_table := `CREATE TABLE IF NOT EXISTS subscription (
	"user" int8 primary key,
	"ooxx" bool,
	"pic" bool) WITH (OIDS=FALSE);
	ALTER TABLE subscription ADD CONSTRAINT "subscribe-user" FOREIGN KEY ("user") REFERENCES "user" ("id") ON UPDATE NO ACTION ON DELETE CASCADE NOT DEFERRABLE INITIALLY IMMEDIATE;`
	_, err := db.sqldb.Exec(sql_table)
	if err != nil {
		logger.Debug(err)
	}

	return err
}

func (db *DB) createJandanTable() error {
	logger.Debug("Creating Table jandan")

	sql_table := `CREATE TABLE IF NOT EXISTS jandan (
	"id" serial primary key,
	"content" text NOT NULL COLLATE "default",
	"category" varchar NOT NULL COLLATE "default",
	"oo" int4,
	"xx" int4,
	"top" bool DEFAULT false) WITH (OIDS=FALSE);`
	_, err := db.sqldb.Exec(sql_table)
	if err != nil {

		logger.Debug(err)
	}

	return err
}

func (db *DB) createSourceTable() error {
	logger.Debug("Create Table source")
	sql_table := `CREATE TABLE IF NOT EXISTS source (
	"id" serial primary key,
	"name" varchar(1024),
	"url" varchar(2048)
	)`
	_, err := db.sqldb.Exec(sql_table)
	if err != nil {

		logger.Debug(err)
	}
	return err
}
