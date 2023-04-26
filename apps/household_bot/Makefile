.SILENT:
.PHONY:
.DEFAULT_GOAL := run

run:
	./scripts/run-dev.sh

unit-test:
	go test -count=2 -short ./...

ci: unit-test
	go vet ./...
	./scripts/e2e.sh

BOT_MOCK_SRC=internal/telegram/handler/handler.go
# keep them not under internal because of e2e tests
BOT_MOCK_DST=mocks/mocks.go

mocks:
	mockgen -source ${BOT_MOCK_SRC} -destination ${BOT_MOCK_DST}
