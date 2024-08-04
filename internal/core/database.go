package core

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	_ "modernc.org/sqlite"
)

func NewDatabase(env Env, logger *zap.Logger) *sqlx.DB {

	db, err := sqlx.Connect(env.DatabaseDriver, env.DatabaseUrl)
	logger.Info(fmt.Sprintf("%s %s", env.DatabaseDriver, env.DatabaseUrl),
		zap.Any("database", db))
	if err != nil {
		logger.Error("database fail", zap.Error(err))
	}
	//db.Exec("DROP TABLE IF EXISTS user;")
	//db.Exec(`CREATE TABLE user (
	//	user_id    INTEGER PRIMARY KEY,
	//	first_name VARCHAR(80)  DEFAULT '',
	//	last_name  VARCHAR(80)  DEFAULT '',
	//	email      VARCHAR(250) DEFAULT '',
	//	password   VARCHAR(250) DEFAULT NULL
	//);`)
	return db
}
