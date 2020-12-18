package main

import (
	"Week04/pkg/notice"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.WithError(
		notice.Run(),
	).Error("Service exited")
}
