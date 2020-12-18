package database

import (
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func Transaction(f func(*gorm.DB) error) error {
	gdb := Begin()

	err := f(gdb)
	if err != nil {
		log.WithError(err).Error("db transaction failed")
		gdb.Rollback()
		return err
	}

	err = gdb.Commit().Error
	if err != nil {
		log.WithError(err).Error("db transaction commit failed")
		return err
	}

	return nil
}

func FoundRecord(err error) (bool, error) {
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err == nil {
		return true, nil
	}

	return false, err
}

func RecordCount(tx *gorm.DB) (int64, *gorm.DB) {
	c := &struct {
		Size int64
	}{
		Size: 0,
	}

	tx = tx.Select("count(0) as `size`").Scan(c)
	return c.Size, tx
}

func SQLExpr(sql string, args []interface{}) (string, interface{}) {
	if len(args) > 0 {
		sql += " ?"
	}

	s := " "
	for range args {
		s += " ? "
	}

	return sql, gorm.Expr(s, args...)
}
