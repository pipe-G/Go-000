package notice

import (
	"time"

	micro "github.com/micro/go-micro"
	log "github.com/sirupsen/logrus"
)

var (
	ServiceName = "Week04.notice"
	Version     = "1.0.0"
)

func Run() error {
	noticeSvc := NewNoticeService()

	microSvc := micro.NewService(
		micro.Name(ServiceName),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Version(Version),
	)

	if err := noticeSvc.init(microSvc); err != nil {
		log.WithError(err).Error("NoticeService init error")
		return err
	}

	noticeSvc.setupHandlersRights()
	return microSvc.Run()
}
