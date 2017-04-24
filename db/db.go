package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/qsz13/ooxxbot/logger"
	"os"
	"strconv"
	"time"
)

type DB struct {
	sqldb *sql.DB
}

func NewDB() *DB {
	db := &DB{sqldb: initDBConn(os.Getenv("DATABASE_URL"))}
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
}

func (db *DB) createUserTable() error {
	logger.Debug("Create Table user")
	sql_table := `
	CREATE SEQUENCE IF NOT EXISTS user_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1;
	CREATE TABLE IF NOT EXISTS "user" (
	id int8 NOT NULL DEFAULT nextval('user_id_seq'::regclass),
	first_name varchar(32),
	last_name varchar(32),
	user_name varchar(32),
	constraint pk_user primary key (id)
	) WITH (OIDS=FALSE);
	CREATE UNIQUE INDEX IF NOT EXISTS "user_id_key" ON "user" USING btree("id" "pg_catalog"."int8_ops" ASC NULLS LAST);`
	_, err := db.sqldb.Exec(sql_table)
	return err
}

func (db *DB) createSubscriptionTable() error {
	logger.Debug("Create Table subscription")
	sql_table := `CREATE TABLE IF NOT EXISTS subscription (
	"user" int8 NOT NULL,
	"ooxx" bool,
	"pic" bool) WITH (OIDS=FALSE);
	ALTER TABLE subscription ADD PRIMARY KEY ("user") NOT DEFERRABLE INITIALLY IMMEDIATE;
	ALTER TABLE subscription ADD CONSTRAINT "subscribe-user" FOREIGN KEY ("user") REFERENCES "user" ("id") ON UPDATE NO ACTION ON DELETE CASCADE NOT DEFERRABLE INITIALLY IMMEDIATE;`
	_, err := db.sqldb.Exec(sql_table)
	return err
}

func (db *DB) createJandanTable() error {
	logger.Debug("Create Table jandan")

	sql_table := `CREATE TABLE jandan (
	"id" int4 NOT NULL,
	"content" varchar(1024) NOT NULL COLLATE "default",
	"category" varchar NOT NULL COLLATE "default",
	"oo" int4,
	"xx" int4,
	"top" bool DEFAULT false) WITH (OIDS=FALSE);
	ALTER TABLE jandan ADD PRIMARY KEY ("id") NOT DEFERRABLE INITIALLY IMMEDIATE;`
	_, err := db.sqldb.Exec(sql_table)
	return err
}
