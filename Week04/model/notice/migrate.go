package notice

import (
	"Week04/lib/db"
)

func AutoMigrate() error {
	return db.DB.AutoMigrate(
		&NoticeSMS{}, &NoticeEmail{},
	).Error
}
