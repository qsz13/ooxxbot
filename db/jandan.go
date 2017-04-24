package db

import (
	_ "github.com/lib/pq"
	jd "github.com/qsz13/ooxxbot/jandan"
	"github.com/qsz13/ooxxbot/logger"
	"strconv"
)

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
