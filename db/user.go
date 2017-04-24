package db

import (
	"github.com/qsz13/ooxxbot/logger"
	"github.com/qsz13/ooxxbot/model"
)

func (db *DB) RegisterUser(user *model.User) error {
	stmt, err := db.sqldb.Prepare("INSERT INTO \"user\"(id, first_name, last_name, user_name) VALUES ($1, $2, $3, $4) ON CONFLICT (id) DO UPDATE SET first_name=excluded.first_name,last_name=excluded.last_name,user_name=excluded.user_name;")
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(user.ID, user.FirstName, user.LastName, user.Username)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	return nil
}
