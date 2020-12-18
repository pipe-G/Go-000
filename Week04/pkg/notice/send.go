package notice

import (
	mdn "Week04/model/notice"
	"time"

	log "github.com/sirupsen/logrus"
	gomail "gopkg.in/gomail.v2"
)

func (svc *NoticeService) sendNotice(
	sendingMap map[int64]bool, notice *mdn.Notice, sendFunc func() error,
) bool {
	svc.mu.Lock()
	{
		sending := sendingMap[notice.ID]

		if sending {
			svc.mu.Unlock()
			return false
		}

		sendingMap[notice.ID] = true
	}
	svc.mu.Unlock()

	defer func() {
		svc.mu.Lock()
		delete(sendingMap, notice.ID)
		svc.mu.Unlock()
	}()

	now := time.Now()
	notice.UpdatedAt = &now
	notice.SentTimes++

	err := sendFunc()
	if err != nil {
		notice.Status = mdn.ES_FAILED
	} else {
		notice.Status = mdn.ES_SUCCESS
	}

	return true
}



func (svc *NoticeService) sendMail(mail *mdn.NoticeEmail) {
	if svc.sendNotice(svc.sendingEmails, &mail.Notice, func() error {
		m := gomail.NewMessage()
		m.SetBody("text/html", mail.Body)
		m.SetHeaders(map[string][]string{
			"From": {m.FormatAddress(
				svc.config.SMTP.From,
				svc.config.SMTP.SenderName,
			)},
			"To":      {mail.Email},
			"Subject": {mail.Subject},
		})

		err := svc.mailDialer.DialAndSend(m)
		if err != nil {
			log.WithError(err).
				WithField("Email", mail).
				Error("Send email failed")
		}
		return err
	}) {
		if err := mdn.SaveNoticeEmail(mail); err != nil {
			log.WithField("mail", mail).
				WithError(err).Error("SaveNoticeEmail failed")
		}
	}
}
func (svc *NoticeService) retryEmailLoop() {
	for {
		emails, err := mdn.GetNoticeEmails(map[string]interface{}{
			"Status":        mdn.ES_FAILED,
			"UpdatedAtLess": time.Now().Add(time.Minute * -3),
			"MaxSentTimes":  3,
		}, 0)

		if err != nil {
			log.WithError(err).
				Error("GetNoticeEmails failed")
			time.Sleep(time.Second * 5)
			continue
		}

		for _, email := range emails {
			svc.sendMail(email)
		}

		time.Sleep(time.Second * 5)
	}
}
