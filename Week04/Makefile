PROTOC = protoc
PROTO_MICRO_BUILD = $(PROTOC) --go_out=paths=source_relative:. --micro_out=paths=source_relative:.
BUILD = go build -ldflags '-w'

ALL: notice

.PHONY: proto

proto:
	$(PROTO_MICRO_BUILD) proto/notice/notice.proto


notice:
	go build -o ./_build/bin/notice cmd/notice/*
run-notice:
	make notice
	./_build/bin/notice --broker nats

db-reset:
	  mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD} ivehicle < sql/0000_init.down.sql
	  mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD} ivehicle < sql/0000_init.up.sql
	  mysql -u${MYSQL_USER} -p${MYSQL_PASSWORD} ivehicle < sql/test.sql
