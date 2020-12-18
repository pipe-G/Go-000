package notice

import (
	mdn "Week04/model/notice"
	pbnotice "Week04/proto/notice"
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
)

func (svc *NoticeService) SendEmail(
	ctx context.Context, req *pbnotice.SendEmailReq) error {

	now := time.Now()
	mailModel := &mdn.NoticeEmail{
		Notice: mdn.Notice{
			UserID:    req.UserID,
			Status:    mdn.ES_INIT,
			SentTimes: 0,
			UpdatedAt: &now,
			CreatedAt: &now,
		},
		Email:   req.Email,
		Subject: req.Subject,
		Body:    req.Body,
	}

	if err := mdn.SaveNoticeEmail(mailModel); err != nil {
		log.WithField("mail", mailModel).
			WithError(err).
			Error("SaveNoticeEmail failed")
		return err
	}

	go svc.sendMail(mailModel)

	return nil
}

type HandlerAttr struct {
	NoLogReq         bool
}


var registeredHandlers = map[string]*HandlerAttr{}
func ErrExit(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func RegisterHandler(handlerName string, attr *HandlerAttr) {
	handler, ok := registeredHandlers[handlerName]
	if ok {
		ErrExit(fmt.Errorf(
			"handler already existed: %s, right: %+v",
			handlerName, handler,
		))
	}
	registeredHandlers[handlerName] = attr
}

func (svc *NoticeService) setupHandlersRights() {
	RegisterHandler("Notice.SendEmail", &HandlerAttr{
		NoLogReq: true,
	})
}