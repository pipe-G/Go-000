package notice

import (
	"Week04/lib/db"
	"time"

	goqu "gopkg.in/doug-martin/goqu.v5"
)

const (
	ES_INIT    int32 = 0
	ES_SUCCESS       = 1
	ES_FAILED        = 2
)

type Notice struct {
	ID        int64 `gorm:"primary_key"`
	UserID    int64
	Status    int32
	SentTimes int32
	ExpireAt  *time.Time
	UpdatedAt *time.Time
	CreatedAt *time.Time
}

type NoticeEmail struct {
	Notice

	Email   string
	Subject string
	Body    string
}


func SaveNoticeEmail(email *NoticeEmail) error {
	return db.DB.Save(email).Error
}

func GetNoticeEmails(params map[string]interface{}, limit uint) ([]*NoticeEmail, error) {
	var emails []*NoticeEmail

	builder := db.DB.Goqu.From(goqu.I("notice_emails").As("e")).
		Where(
			db.QueryFilter(params).
				Where("id", db.EQ, "ID").
				Where("user_id", db.EQ, "UserID").
				Where("email", db.EQ, "Email").
				Where("status", db.EQ, "Status").
				Where("updated_at", db.LT, "UpdatedAtLess").
				Where("sent_times", db.LT, "MaxSentTimes").
				End()...,
		)

	if limit > 0 {
		builder = builder.Limit(limit)
	}

	err := db.DB.QueryAll(builder, &emails)
	return emails, err
}

