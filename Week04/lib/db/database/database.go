package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var defaultDB *gorm.DB

func openDB(debug bool, dialect, url string) (*gorm.DB, error) {
	db, err := gorm.Open(dialect, url)
	if err != nil {
		return nil, err
	}

	if debug {
		db.LogMode(true)
	}

	return db, nil
}

func Open(debug bool, dialect, dburl string) error {
	db, err := openDB(debug, dialect, dburl)
	if err != nil {
		return err
	}
	defaultDB = db

	return nil
}

func DB() *gorm.DB {
	return defaultDB
}

func Begin() *gorm.DB {
	return defaultDB.Begin()
}

func Commit() *gorm.DB {
	return defaultDB.Commit()
}

func Rollback() *gorm.DB {
	return defaultDB.Rollback()
}
