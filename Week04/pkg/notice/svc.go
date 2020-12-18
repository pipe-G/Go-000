package notice

import (
	"Week04/lib/db"
	mdn "Week04/model/notice"
	"sync"

	micro "github.com/micro/go-micro"
	gomail "gopkg.in/gomail.v2"
)

type NoticeService struct {
	config        Notice

	mailDialer *gomail.Dialer

	sendingEmails map[int64]bool
	mu            sync.Mutex
}

func NewNoticeService() *NoticeService {
	return &NoticeService{
		sendingEmails: make(map[int64]bool),
	}
}


type SMTP struct {
	Host       string
	Port       int
	From       string
	SenderName string
	Username   string
	Password   string
}

type Notice struct {
	SMTP SMTP
}

type DB struct {
	URL   string
	Debug bool
}

func (svc *NoticeService) init(microSvc micro.Service) error {
	var (
		dbConf     DB
		noticeConf Notice
	)


	if err := db.Open(dbConf.Debug, "mysql", dbConf.URL); err != nil {
		return err
	}

	if err := mdn.AutoMigrate(); err != nil {
		return err
	}

	svc.mailDialer = gomail.NewDialer(
		noticeConf.SMTP.Host,
		noticeConf.SMTP.Port,
		noticeConf.SMTP.Username,
		noticeConf.SMTP.Password,
	)


	svc.config = noticeConf

	go svc.retryEmailLoop()

	return nil
}
