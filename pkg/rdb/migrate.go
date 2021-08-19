package rdb

import "go-chi-ddd/domain/entity"

func migrate() error {
	return db.AutoMigrate(
		&entity.User{},
	)
}
