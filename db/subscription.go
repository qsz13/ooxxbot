package db

import (
	"github.com/qsz13/ooxxbot/logger"
)

func (db *DB) SubscribeJandanPic(uid int) error {
	stmt, err := db.sqldb.Prepare("INSERT INTO subscription(\"user\", pic) VALUES ( $1, $2) ON CONFLICT (\"user\") DO UPDATE SET \"user\"=excluded.\"user\",pic=excluded.pic;")
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(uid, true)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (db *DB) UnsubscribeJandanPic(uid int) error {
	stmt, err := db.sqldb.Prepare("INSERT INTO subscription(\"user\", pic) VALUES ( $1, $2) ON CONFLICT (\"user\") DO UPDATE SET \"user\"=excluded.\"user\", pic=excluded.pic;")
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(uid, false)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (db *DB) SubscribeJandanOOXX(uid int) error {
	stmt, err := db.sqldb.Prepare("INSERT INTO subscription(\"user\", ooxx) VALUES ( $1, $2) ON CONFLICT (\"user\") DO UPDATE SET \"user\"=excluded.\"user\",ooxx=excluded.ooxx;")
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(uid, true)

	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil

}

func (db *DB) UnsubscribeJandanOOXX(uid int) error {
	stmt, err := db.sqldb.Prepare("INSERT INTO subscription(\"user\", ooxx) VALUES ( $1, $2) ON CONFLICT (\"user\") DO UPDATE SET \"user\"=excluded.\"user\",ooxx=excluded.ooxx;")
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(uid, false)

	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil

}

func (db *DB) GetPicSubscriber() ([]int, error) {
	subscribers := []int{}
	rows, err := db.sqldb.Query("SELECT \"user\" FROM subscription where pic=true;")
	if err != nil {
		logger.Error(err.Error())
		return subscribers, err
	}

	defer rows.Close()
	sid := 0
	for rows.Next() {
		err = rows.Scan(&sid)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		subscribers = append(subscribers, sid)
	}
	if err != nil {
		logger.Error(err.Error())
	}
	return subscribers, err
}

func (db *DB) GetOOXXSubscriber() ([]int, error) {
	subscribers := []int{}
	rows, err := db.sqldb.Query("SELECT \"user\" FROM subscription where ooxx=true;")
	if err != nil {
		logger.Error(err.Error())
		return subscribers, err
	}

	defer rows.Close()
	sid := 0
	for rows.Next() {
		err = rows.Scan(&sid)
		if err != nil {
			logger.Error(err.Error())
			continue
		}
		subscribers = append(subscribers, sid)
	}
	if err != nil {
		logger.Error(err.Error())
	}
	return subscribers, err
}
