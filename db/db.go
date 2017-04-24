package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	jd "github.com/qsz13/ooxxbot/jandan"
	"github.com/qsz13/ooxxbot/logger"
	"os"
	"strconv"
	"time"
)

type DB struct {
	sqldb *sql.DB
}

func NewDB() *DB {
	sqldb := initDBConn(os.Getenv("DATABASE_URL"))
	return &DB{sqldb: sqldb}
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

func (db *DB) SaveJandanTops(tops []jd.Comment) {
	sqlStr := "INSERT INTO jandan (id, content, category, oo, xx, top) VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (id) DO UPDATE SET top=true;"

	stmt, err := db.sqldb.Prepare(sqlStr)
	defer stmt.Close()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	for _, comment := range tops {
		_, err = stmt.Exec(comment.ID, comment.Content, strconv.Itoa(int(comment.Type)), strconv.Itoa(int(comment.OO)), strconv.Itoa(int(comment.XX)), false)
		if err != nil {
			logger.Error(err.Error())
		}

	}

	//format all vals at once
	logger.Debug("Top Saved!")

}

func (db *DB) SaveJandanToDB(comments []jd.Comment) {
	sqlStr := "INSERT INTO jandan (id, content, category, oo, xx, top) VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT (id) DO UPDATE SET content=excluded.content,category=excluded.category,oo=excluded.oo, xx=excluded.xx;"
	stmt, err := db.sqldb.Prepare(sqlStr)
	defer stmt.Close()

	if err != nil {
		logger.Error(err.Error())
		return
	}
	for _, comment := range comments {
		_, err = stmt.Exec(comment.ID, comment.Content, strconv.Itoa(int(comment.Type)), strconv.Itoa(int(comment.OO)), strconv.Itoa(int(comment.XX)), false)
		if err != nil {
			logger.Error(err.Error())
		}
	}
	//format all vals at once
	logger.Debug("Comment Saved!")

}

func (db *DB) GetRandomComment(jdType jd.JandanType) (string, error) {
	content := ""
	stmt, err := db.sqldb.Prepare("SELECT content FROM jandan WHERE category=$1 AND oo > 2*xx ORDER BY RANDOM() LIMIT 1;")
	if err != nil {
		logger.Error(err.Error())
		return content, err
	}
	err = stmt.QueryRow(jdType).Scan(&content)
	if err != nil {
		logger.Error(err.Error())
	}
	return content, err
}

func (db *DB) TopExists(top *jd.Comment) bool {
	stmt, err := db.sqldb.Prepare("SELECT top from jandan where id = $1")
	if err != nil {
		logger.Error(err.Error())
		return true
	}

	var isTop bool
	err = stmt.QueryRow(top.ID).Scan(&isTop)
	if err != nil {
		return false
	}
	if isTop {
		return true
	}
	return false

}
